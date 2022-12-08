// SPDX-License-Identifier: AGPL-3.0-only

package encoding

import (
	"bufio"
	"encoding/binary"
	"hash/crc32"
	"os"
	"sync"

	"github.com/grafana/dskit/multierror"
	"github.com/pkg/errors"
)

// readerBufferSize is the size of the buffer used for reading index-header files. This
// value is arbitrary and will likely change in the future based on profiling results.
const readerBufferSize = 4096

// DecbufFactory creates new file-backed decoding buffer instances for a specific index-header file.
type DecbufFactory struct {
	files   *filePool
	buffers sync.Pool
}

func NewDecbufFactory(path string, maxFileHandles uint) *DecbufFactory {
	return &DecbufFactory{
		files: newFilePool(path, maxFileHandles),
		buffers: sync.Pool{
			New: func() any {
				return bufio.NewReaderSize(nil, readerBufferSize)
			},
		},
	}
}

// NewDecbufAtChecked returns a new file-backed decoding buffer positioned at offset + 4 bytes.
// It expects the first 4 bytes after offset to hold the big endian encoded content length, followed
// by the contents and the expected checksum. This method checks the CRC of the content and will
// return an error Decbuf if it does not match the expected CRC.
func (df *DecbufFactory) NewDecbufAtChecked(offset int, table *crc32.Table) Decbuf {
	f, err := df.files.get()
	if err != nil {
		return Decbuf{E: errors.Wrap(err, "open file for decbuf")}
	}

	// If we return early and don't include a Reader for our Decbuf, we are responsible
	// for putting the file handle back in the pool.
	closeFile := true
	defer func() {
		if closeFile {
			_ = df.files.put(f)
		}
	}()

	// TODO: A particular index-header only has symbols and posting offsets. We should only need to read
	//  the length of each of those a single time per index-header (DecbufFactory). Should the factory
	//  cache the length? Should the table of contents be passed to the factory?
	lengthBytes := make([]byte, 4)
	n, err := f.ReadAt(lengthBytes, int64(offset))
	if err != nil {
		return Decbuf{E: err}
	}
	if n != 4 {
		return Decbuf{E: errors.Wrapf(ErrInvalidSize, "insufficient bytes read for size (got %d, wanted %d)", n, 4)}
	}

	bufReader := df.buffers.Get().(*bufio.Reader)
	contentLength := int(binary.BigEndian.Uint32(lengthBytes))
	bufferLength := len(lengthBytes) + contentLength + crc32.Size
	r, err := NewFileReader(f, offset, bufferLength, bufReader)
	if err != nil {
		return Decbuf{E: errors.Wrap(err, "create file reader")}
	}

	closeFile = false
	d := Decbuf{r: r}

	if d.ResetAt(4); d.Err() != nil {
		return d
	}

	if table != nil {
		if d.CheckCrc32(table); d.Err() != nil {
			return d
		}

		// Reset to the beginning of the content after reading it all for the CRC.
		d.ResetAt(4)
	}

	return d
}

// NewDecbufAtUnchecked returns a new file-backed decoding buffer positioned at offset + 4 bytes.
// It expects the first 4 bytes after offset to hold the big endian encoded content length, followed
// by the contents and the expected checksum. This method does NOT compute the CRC of the content.
// To check the CRC of the content, use NewDecbufAtChecked.
func (df *DecbufFactory) NewDecbufAtUnchecked(offset int) Decbuf {
	return df.NewDecbufAtChecked(offset, nil)
}

// NewRawDecbuf returns a new file-backed decoding buffer positioned at the beginning of the file,
// spanning the entire length of the file. It does not make any assumptions about the contents of the
// file, nor does it perform any form of integrity check. To create a decoding buffer for some subset
// of the file or perform integrity checks use NewDecbufAtUnchecked or NewDecbufAtChecked.
func (df *DecbufFactory) NewRawDecbuf() Decbuf {
	f, err := df.files.get()
	if err != nil {
		return Decbuf{E: errors.Wrap(err, "open file for decbuf")}
	}

	// If we return early and don't include a Reader for our Decbuf, we are responsible
	// for putting the file handle back in the pool.
	closeFile := true
	defer func() {
		if closeFile {
			_ = df.files.put(f)
		}
	}()

	stat, err := f.Stat()
	if err != nil {
		return Decbuf{E: errors.Wrap(err, "stat file for decbuf")}
	}

	bufReader := df.buffers.Get().(*bufio.Reader)
	fileSize := stat.Size()
	reader, err := NewFileReader(f, 0, int(fileSize), bufReader)
	if err != nil {
		return Decbuf{E: errors.Wrap(err, "file reader for decbuf")}
	}

	closeFile = false
	return Decbuf{r: reader}
}

// Stop cleans up resources associated with this DecbufFactory
func (df *DecbufFactory) Stop() {
	df.files.stop()
}

// Close cleans up any resources associated with the Decbuf
func (df *DecbufFactory) Close(d Decbuf) error {
	if d.r != nil {
		df.buffers.Put(d.r.buf)
		return df.files.put(d.r.file)
	}

	return nil
}

// CloseWithErrCapture cleans up any resources associated with d,
// capturing any errors that occur while cleaning up d in err.
func (df *DecbufFactory) CloseWithErrCapture(err *error, d Decbuf, format string, a ...interface{}) {
	merr := multierror.MultiError{}

	merr.Add(*err)
	merr.Add(errors.Wrapf(df.Close(d), format, a...))

	*err = merr.Err()
}

// filePool maintains a pool of file handles up to a maximum number, creating
// new handles and closing them when required. get and put operations on this
// pool never block. If there are no available file handles, one will be created
// on get. If the pool is full, the file handle is closed on put.
type filePool struct {
	path    string
	handles chan *os.File
}

// newFilePool creates a new file pool for path with cap capacity. If cap is 0,
// get always opens new file handles and put always closes them immediately.
func newFilePool(path string, cap uint) *filePool {
	return &filePool{
		path: path,
		// We don't care if cap is 0 which means the channel will be unbuffered. Because
		// we have default cases for reads and writes to the channel, we will always open
		// new files and close file handles immediately if the channel is unbuffered.
		handles: make(chan *os.File, cap),
	}
}

// get returns a pooled file handle if available or opens a new one if there
// are no pooled handles available.
func (p *filePool) get() (*os.File, error) {
	select {
	case f := <-p.handles:
		return f, nil
	default:
		return os.Open(p.path)
	}
}

// put returns a file handle to the pool if there is space available or closes
// the file handle if there is not.
func (p *filePool) put(f *os.File) error {
	select {
	case p.handles <- f:
		return nil
	default:
		return f.Close()
	}
}

// stop closes all pooled file handles.
func (p *filePool) stop() {
	for {
		select {
		case f := <-p.handles:
			_ = f.Close()
		default:
			return
		}
	}
}
