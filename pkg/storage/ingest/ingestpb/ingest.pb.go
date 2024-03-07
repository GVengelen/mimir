// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ingest.proto

package ingestpb

import (
	bytes "bytes"
	context "context"
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	mimirpb "github.com/grafana/mimir/pkg/mimirpb"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

type Segment struct {
	Pieces []*Piece `protobuf:"bytes,1,rep,name=pieces,proto3" json:"pieces,omitempty"`
}

func (m *Segment) Reset()      { *m = Segment{} }
func (*Segment) ProtoMessage() {}
func (*Segment) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{0}
}
func (m *Segment) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Segment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Segment.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Segment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Segment.Merge(m, src)
}
func (m *Segment) XXX_Size() int {
	return m.Size()
}
func (m *Segment) XXX_DiscardUnknown() {
	xxx_messageInfo_Segment.DiscardUnknown(m)
}

var xxx_messageInfo_Segment proto.InternalMessageInfo

func (m *Segment) GetPieces() []*Piece {
	if m != nil {
		return m.Pieces
	}
	return nil
}

type Piece struct {
	ObsoleteWriteRequest *mimirpb.WriteRequest `protobuf:"bytes,1,opt,name=obsolete_write_request,json=obsoleteWriteRequest,proto3" json:"obsolete_write_request,omitempty"`
	TenantId             string                `protobuf:"bytes,2,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	Data                 []byte                `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	CreatedAtMs          int64                 `protobuf:"varint,4,opt,name=created_at_ms,json=createdAtMs,proto3" json:"created_at_ms,omitempty"`
	SnappyEncoded        bool                  `protobuf:"varint,6,opt,name=snappy_encoded,json=snappyEncoded,proto3" json:"snappy_encoded,omitempty"`
}

func (m *Piece) Reset()      { *m = Piece{} }
func (*Piece) ProtoMessage() {}
func (*Piece) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{1}
}
func (m *Piece) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Piece) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Piece.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Piece) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Piece.Merge(m, src)
}
func (m *Piece) XXX_Size() int {
	return m.Size()
}
func (m *Piece) XXX_DiscardUnknown() {
	xxx_messageInfo_Piece.DiscardUnknown(m)
}

var xxx_messageInfo_Piece proto.InternalMessageInfo

func (m *Piece) GetObsoleteWriteRequest() *mimirpb.WriteRequest {
	if m != nil {
		return m.ObsoleteWriteRequest
	}
	return nil
}

func (m *Piece) GetTenantId() string {
	if m != nil {
		return m.TenantId
	}
	return ""
}

func (m *Piece) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Piece) GetCreatedAtMs() int64 {
	if m != nil {
		return m.CreatedAtMs
	}
	return 0
}

func (m *Piece) GetSnappyEncoded() bool {
	if m != nil {
		return m.SnappyEncoded
	}
	return false
}

type WriteResponse struct {
}

func (m *WriteResponse) Reset()      { *m = WriteResponse{} }
func (*WriteResponse) ProtoMessage() {}
func (*WriteResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{2}
}
func (m *WriteResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WriteResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WriteResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WriteResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteResponse.Merge(m, src)
}
func (m *WriteResponse) XXX_Size() int {
	return m.Size()
}
func (m *WriteResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteResponse.DiscardUnknown(m)
}

var xxx_messageInfo_WriteResponse proto.InternalMessageInfo

type WriteRequest struct {
	Piece       *Piece `protobuf:"bytes,1,opt,name=piece,proto3" json:"piece,omitempty"`
	PartitionId int32  `protobuf:"varint,2,opt,name=partition_id,json=partitionId,proto3" json:"partition_id,omitempty"`
}

func (m *WriteRequest) Reset()      { *m = WriteRequest{} }
func (*WriteRequest) ProtoMessage() {}
func (*WriteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff993cce43359ffa, []int{3}
}
func (m *WriteRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *WriteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_WriteRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *WriteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WriteRequest.Merge(m, src)
}
func (m *WriteRequest) XXX_Size() int {
	return m.Size()
}
func (m *WriteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_WriteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_WriteRequest proto.InternalMessageInfo

func (m *WriteRequest) GetPiece() *Piece {
	if m != nil {
		return m.Piece
	}
	return nil
}

func (m *WriteRequest) GetPartitionId() int32 {
	if m != nil {
		return m.PartitionId
	}
	return 0
}

func init() {
	proto.RegisterType((*Segment)(nil), "ingest.Segment")
	proto.RegisterType((*Piece)(nil), "ingest.Piece")
	proto.RegisterType((*WriteResponse)(nil), "ingest.WriteResponse")
	proto.RegisterType((*WriteRequest)(nil), "ingest.WriteRequest")
}

func init() { proto.RegisterFile("ingest.proto", fileDescriptor_ff993cce43359ffa) }

var fileDescriptor_ff993cce43359ffa = []byte{
	// 412 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xcf, 0x8a, 0xd4, 0x40,
	0x10, 0xc6, 0xd3, 0xce, 0xce, 0x38, 0xdb, 0x33, 0x51, 0x68, 0xd6, 0x25, 0xac, 0xd0, 0xc4, 0xc8,
	0x42, 0x4e, 0x99, 0x65, 0x04, 0x8f, 0xc2, 0x8a, 0x1e, 0x16, 0x14, 0x24, 0x82, 0x82, 0x97, 0xd0,
	0x49, 0xca, 0xd8, 0x68, 0xba, 0xdb, 0xee, 0x5a, 0xd4, 0x9b, 0x8f, 0xe0, 0x63, 0xf8, 0x28, 0x1e,
	0xc7, 0xdb, 0x1e, 0x9d, 0xcc, 0xc5, 0xe3, 0x3e, 0x82, 0x4c, 0x3a, 0x23, 0xbb, 0x7b, 0xea, 0xaa,
	0xdf, 0x47, 0xfd, 0xf9, 0xaa, 0xe9, 0x5c, 0xaa, 0x06, 0x1c, 0x66, 0xc6, 0x6a, 0xd4, 0x6c, 0xe2,
	0xb3, 0xa3, 0x93, 0x46, 0xe2, 0x87, 0xf3, 0x32, 0xab, 0x74, 0xbb, 0x68, 0xac, 0x78, 0x2f, 0x94,
	0x58, 0xb4, 0xb2, 0x95, 0x76, 0x61, 0x3e, 0x36, 0x3e, 0x32, 0xa5, 0x7f, 0x7d, 0x65, 0x72, 0x42,
	0x6f, 0xbf, 0x86, 0xa6, 0x05, 0x85, 0xec, 0x98, 0x4e, 0x8c, 0x84, 0x0a, 0x5c, 0x44, 0xe2, 0x51,
	0x3a, 0x5b, 0x86, 0xd9, 0x30, 0xe3, 0xd5, 0x96, 0xe6, 0x83, 0x98, 0xfc, 0x26, 0x74, 0xdc, 0x13,
	0xf6, 0x82, 0x1e, 0xea, 0xd2, 0xe9, 0x4f, 0x80, 0x50, 0x7c, 0xb1, 0x12, 0xa1, 0xb0, 0xf0, 0xf9,
	0x1c, 0x1c, 0x46, 0x24, 0x26, 0xe9, 0x6c, 0x79, 0x98, 0x55, 0xda, 0x22, 0x7c, 0x35, 0x65, 0xf6,
	0x76, 0x2b, 0xe7, 0x5e, 0xcd, 0x0f, 0x76, 0x55, 0x57, 0x29, 0xbb, 0x4f, 0xf7, 0x11, 0x94, 0x50,
	0x58, 0xc8, 0x3a, 0xba, 0x15, 0x93, 0x74, 0x3f, 0x9f, 0x7a, 0x70, 0x56, 0x33, 0x46, 0xf7, 0x6a,
	0x81, 0x22, 0x1a, 0xc5, 0x24, 0x9d, 0xe7, 0x7d, 0xcc, 0x12, 0x1a, 0x56, 0x16, 0x04, 0x42, 0x5d,
	0x08, 0x2c, 0x5a, 0x17, 0xed, 0xc5, 0x24, 0x1d, 0xe5, 0xb3, 0x01, 0x9e, 0xe2, 0x4b, 0xc7, 0x8e,
	0xe9, 0x1d, 0xa7, 0x84, 0x31, 0xdf, 0x0a, 0x50, 0x95, 0xae, 0xa1, 0x8e, 0x26, 0x31, 0x49, 0xa7,
	0x79, 0xe8, 0xe9, 0x73, 0x0f, 0x93, 0xbb, 0x34, 0x1c, 0x76, 0x71, 0x46, 0x2b, 0x07, 0xc9, 0x1b,
	0x3a, 0xbf, 0xb6, 0xdc, 0x43, 0x3a, 0xee, 0xed, 0x0f, 0xce, 0x6e, 0x9c, 0xc6, 0x6b, 0xec, 0x01,
	0x9d, 0x1b, 0x61, 0x51, 0xa2, 0xd4, 0x6a, 0x67, 0x62, 0x9c, 0xcf, 0xfe, 0xb3, 0xb3, 0x7a, 0xf9,
	0x8c, 0xd2, 0xbe, 0xef, 0x69, 0xb3, 0xbd, 0xf8, 0x63, 0x3a, 0xee, 0x33, 0x76, 0xb0, 0xeb, 0x77,
	0x75, 0xe8, 0xd1, 0xbd, 0x1b, 0x74, 0xd8, 0x2d, 0x78, 0xfa, 0x64, 0xb5, 0xe6, 0xc1, 0xc5, 0x9a,
	0x07, 0x97, 0x6b, 0x4e, 0xbe, 0x77, 0x9c, 0xfc, 0xec, 0x38, 0xf9, 0xd5, 0x71, 0xb2, 0xea, 0x38,
	0xf9, 0xd3, 0x71, 0xf2, 0xb7, 0xe3, 0xc1, 0x65, 0xc7, 0xc9, 0x8f, 0x0d, 0x0f, 0x56, 0x1b, 0x1e,
	0x5c, 0x6c, 0x78, 0xf0, 0x6e, 0xea, 0xbb, 0x99, 0xb2, 0x9c, 0xf4, 0x7f, 0xff, 0xe8, 0x5f, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x37, 0xf0, 0x7e, 0x44, 0x45, 0x02, 0x00, 0x00,
}

func (this *Segment) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Segment)
	if !ok {
		that2, ok := that.(Segment)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if len(this.Pieces) != len(that1.Pieces) {
		return false
	}
	for i := range this.Pieces {
		if !this.Pieces[i].Equal(that1.Pieces[i]) {
			return false
		}
	}
	return true
}
func (this *Piece) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Piece)
	if !ok {
		that2, ok := that.(Piece)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.ObsoleteWriteRequest.Equal(that1.ObsoleteWriteRequest) {
		return false
	}
	if this.TenantId != that1.TenantId {
		return false
	}
	if !bytes.Equal(this.Data, that1.Data) {
		return false
	}
	if this.CreatedAtMs != that1.CreatedAtMs {
		return false
	}
	if this.SnappyEncoded != that1.SnappyEncoded {
		return false
	}
	return true
}
func (this *WriteResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*WriteResponse)
	if !ok {
		that2, ok := that.(WriteResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	return true
}
func (this *WriteRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*WriteRequest)
	if !ok {
		that2, ok := that.(WriteRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.Piece.Equal(that1.Piece) {
		return false
	}
	if this.PartitionId != that1.PartitionId {
		return false
	}
	return true
}
func (this *Segment) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 5)
	s = append(s, "&ingestpb.Segment{")
	if this.Pieces != nil {
		s = append(s, "Pieces: "+fmt.Sprintf("%#v", this.Pieces)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *Piece) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 9)
	s = append(s, "&ingestpb.Piece{")
	if this.ObsoleteWriteRequest != nil {
		s = append(s, "ObsoleteWriteRequest: "+fmt.Sprintf("%#v", this.ObsoleteWriteRequest)+",\n")
	}
	s = append(s, "TenantId: "+fmt.Sprintf("%#v", this.TenantId)+",\n")
	s = append(s, "Data: "+fmt.Sprintf("%#v", this.Data)+",\n")
	s = append(s, "CreatedAtMs: "+fmt.Sprintf("%#v", this.CreatedAtMs)+",\n")
	s = append(s, "SnappyEncoded: "+fmt.Sprintf("%#v", this.SnappyEncoded)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *WriteResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 4)
	s = append(s, "&ingestpb.WriteResponse{")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *WriteRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&ingestpb.WriteRequest{")
	if this.Piece != nil {
		s = append(s, "Piece: "+fmt.Sprintf("%#v", this.Piece)+",\n")
	}
	s = append(s, "PartitionId: "+fmt.Sprintf("%#v", this.PartitionId)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringIngest(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// WriteAgentClient is the client API for WriteAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WriteAgentClient interface {
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
}

type writeAgentClient struct {
	cc *grpc.ClientConn
}

func NewWriteAgentClient(cc *grpc.ClientConn) WriteAgentClient {
	return &writeAgentClient{cc}
}

func (c *writeAgentClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, "/ingest.WriteAgent/Write", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WriteAgentServer is the server API for WriteAgent service.
type WriteAgentServer interface {
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
}

// UnimplementedWriteAgentServer can be embedded to have forward compatible implementations.
type UnimplementedWriteAgentServer struct {
}

func (*UnimplementedWriteAgentServer) Write(ctx context.Context, req *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}

func RegisterWriteAgentServer(s *grpc.Server, srv WriteAgentServer) {
	s.RegisterService(&_WriteAgent_serviceDesc, srv)
}

func _WriteAgent_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriteAgentServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ingest.WriteAgent/Write",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriteAgentServer).Write(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _WriteAgent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ingest.WriteAgent",
	HandlerType: (*WriteAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _WriteAgent_Write_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ingest.proto",
}

func (m *Segment) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Segment) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Segment) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Pieces) > 0 {
		for iNdEx := len(m.Pieces) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Pieces[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintIngest(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Piece) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Piece) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Piece) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.SnappyEncoded {
		i--
		if m.SnappyEncoded {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x30
	}
	if m.CreatedAtMs != 0 {
		i = encodeVarintIngest(dAtA, i, uint64(m.CreatedAtMs))
		i--
		dAtA[i] = 0x20
	}
	if len(m.Data) > 0 {
		i -= len(m.Data)
		copy(dAtA[i:], m.Data)
		i = encodeVarintIngest(dAtA, i, uint64(len(m.Data)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.TenantId) > 0 {
		i -= len(m.TenantId)
		copy(dAtA[i:], m.TenantId)
		i = encodeVarintIngest(dAtA, i, uint64(len(m.TenantId)))
		i--
		dAtA[i] = 0x12
	}
	if m.ObsoleteWriteRequest != nil {
		{
			size, err := m.ObsoleteWriteRequest.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIngest(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *WriteResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WriteResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WriteResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	return len(dAtA) - i, nil
}

func (m *WriteRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *WriteRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *WriteRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.PartitionId != 0 {
		i = encodeVarintIngest(dAtA, i, uint64(m.PartitionId))
		i--
		dAtA[i] = 0x10
	}
	if m.Piece != nil {
		{
			size, err := m.Piece.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintIngest(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintIngest(dAtA []byte, offset int, v uint64) int {
	offset -= sovIngest(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Segment) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Pieces) > 0 {
		for _, e := range m.Pieces {
			l = e.Size()
			n += 1 + l + sovIngest(uint64(l))
		}
	}
	return n
}

func (m *Piece) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ObsoleteWriteRequest != nil {
		l = m.ObsoleteWriteRequest.Size()
		n += 1 + l + sovIngest(uint64(l))
	}
	l = len(m.TenantId)
	if l > 0 {
		n += 1 + l + sovIngest(uint64(l))
	}
	l = len(m.Data)
	if l > 0 {
		n += 1 + l + sovIngest(uint64(l))
	}
	if m.CreatedAtMs != 0 {
		n += 1 + sovIngest(uint64(m.CreatedAtMs))
	}
	if m.SnappyEncoded {
		n += 2
	}
	return n
}

func (m *WriteResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func (m *WriteRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Piece != nil {
		l = m.Piece.Size()
		n += 1 + l + sovIngest(uint64(l))
	}
	if m.PartitionId != 0 {
		n += 1 + sovIngest(uint64(m.PartitionId))
	}
	return n
}

func sovIngest(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIngest(x uint64) (n int) {
	return sovIngest(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Segment) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForPieces := "[]*Piece{"
	for _, f := range this.Pieces {
		repeatedStringForPieces += strings.Replace(f.String(), "Piece", "Piece", 1) + ","
	}
	repeatedStringForPieces += "}"
	s := strings.Join([]string{`&Segment{`,
		`Pieces:` + repeatedStringForPieces + `,`,
		`}`,
	}, "")
	return s
}
func (this *Piece) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Piece{`,
		`ObsoleteWriteRequest:` + strings.Replace(fmt.Sprintf("%v", this.ObsoleteWriteRequest), "WriteRequest", "mimirpb.WriteRequest", 1) + `,`,
		`TenantId:` + fmt.Sprintf("%v", this.TenantId) + `,`,
		`Data:` + fmt.Sprintf("%v", this.Data) + `,`,
		`CreatedAtMs:` + fmt.Sprintf("%v", this.CreatedAtMs) + `,`,
		`SnappyEncoded:` + fmt.Sprintf("%v", this.SnappyEncoded) + `,`,
		`}`,
	}, "")
	return s
}
func (this *WriteResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&WriteResponse{`,
		`}`,
	}, "")
	return s
}
func (this *WriteRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&WriteRequest{`,
		`Piece:` + strings.Replace(this.Piece.String(), "Piece", "Piece", 1) + `,`,
		`PartitionId:` + fmt.Sprintf("%v", this.PartitionId) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringIngest(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Segment) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Segment: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Segment: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Pieces", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Pieces = append(m.Pieces, &Piece{})
			if err := m.Pieces[len(m.Pieces)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Piece) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Piece: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Piece: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObsoleteWriteRequest", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ObsoleteWriteRequest == nil {
				m.ObsoleteWriteRequest = &mimirpb.WriteRequest{}
			}
			if err := m.ObsoleteWriteRequest.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TenantId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.TenantId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Data = append(m.Data[:0], dAtA[iNdEx:postIndex]...)
			if m.Data == nil {
				m.Data = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CreatedAtMs", wireType)
			}
			m.CreatedAtMs = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CreatedAtMs |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SnappyEncoded", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.SnappyEncoded = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WriteResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WriteResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WriteResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *WriteRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: WriteRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: WriteRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Piece", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthIngest
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthIngest
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Piece == nil {
				m.Piece = &Piece{}
			}
			if err := m.Piece.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PartitionId", wireType)
			}
			m.PartitionId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PartitionId |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipIngest(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthIngest
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipIngest(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIngest
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIngest
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthIngest
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthIngest
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowIngest
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipIngest(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthIngest
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthIngest = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIngest   = fmt.Errorf("proto: integer overflow")
)
