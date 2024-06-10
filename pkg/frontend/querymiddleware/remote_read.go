// SPDX-License-Identifier: AGPL-3.0-only

package querymiddleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/prometheus/prometheus/prompb"
	"github.com/prometheus/prometheus/storage/remote"
)

func NewRemoteReadRoundTripper(next http.RoundTripper, _ Limits) http.RoundTripper {
	return &remoteReadRoundTripper{
		next: next,
	}
}

type remoteReadRoundTripper struct {
	next http.RoundTripper
}

func (r *remoteReadRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Parse the request without consuming the body so it can be forwarded to the next round tripper.
	if req.Body == nil {
		// This is not valid, let the next round tripper handle it.
		return r.next.RoundTrip(req)
	}

	// Inspired by ParseRequestFormWithoutConsumingBody from pkg/util/http.go

	// Close the original body reader. It's going to be replaced later in this function.
	origBody := req.Body
	defer func() { _ = origBody.Close() }()

	// Store the body contents, so we can read it multiple times.
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// TODO handle snappy compression.

	// Parse the request data.
	remoteReadRequest := prompb.ReadRequest{}
	if err := proto.Unmarshal(bodyBytes, &remoteReadRequest); err != nil {
		return nil, err
	}

	queries := remoteReadRequest.GetQueries()
	if len(queries) < 1 {
		// This is not valid, let the next round tripper handle it.
		return r.next.RoundTrip(req)
	}

	details := QueryDetailsFromContext(req.Context())
	if details == nil {
		return nil, fmt.Errorf("query details not found in context") // This should not happen as the context is set in the HTTPServer handler.
	}
	details.Params = make(map[string]string)
	query := queries[0]
	details.Params["start"] = fmt.Sprintf("%d", query.GetStartTimestampMs())
	details.Params["end"] = fmt.Sprintf("%d", query.GetEndTimestampMs())

	matchersStrings := make([]string, 0, len(query.Matchers))
	matchers, err := remote.FromLabelMatchers(query.Matchers)
	for _, m := range matchers {
		matchersStrings = append(matchersStrings, m.String())
	}

	details.Params["matchers"] = strings.Join(matchersStrings, ",")

	return r.next.RoundTrip(req)
}
