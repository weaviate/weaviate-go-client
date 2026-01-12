package api

import (
	"encoding/json"
	"net/http"

	proto "github.com/weaviate/weaviate-go-client/v6/internal/api/gen/proto/v1"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

// ResourceExistsResponse is true if the requested resource exists.
//
// Weaviate does not support HEAD requests, so in order to check
// existence of some resource (object, collection, RBAC role)
// the client has to GET that resource instead. Unmarshaling the
// response body of that request is unnecessary, as we are only
// interested in a simple yes/no answer.
// A request that returns HTTP 404 is a "no" and any other
// successful response is a yes.
//
// ResourceExistsResponse imitates HEAD semantics for any HTTP
// and some gRPC requests:
//
//   - HTTP: accepts HTTP 404 code and update itself to true when its
//     UnmarshalJSON method is called.
//   - gRPC Search: updates itself to true if the result set is not empty.
//
// Example:
//
//	func SongExists(ctx context.Context, resourceID string) (bool, error) {
//		req := api.GetSongRequest(ctx, resourceID)
//		var resp api.ResourceExistsResponse
//		if err := transport.Do(ctx, req); err != nil {
//			return false, err
//		}
//		return resp.Bool(), nil
//	}
type ResourceExistsResponse bool

// TODO(dyma): unify with transports.BooleanResponse

// Bool returns bool value of ResourceExistsResponse.
func (exists ResourceExistsResponse) Bool() bool {
	return bool(exists)
}

var (
	_ json.Unmarshaler                      = (*ResourceExistsResponse)(nil)
	_ transports.StatusAccepter             = (*ResourceExistsResponse)(nil)
	_ MessageUnmarshaler[proto.SearchReply] = (*ResourceExistsResponse)(nil)
)

// AcceptStatus implements transport.StatusAccepter.
func (exists ResourceExistsResponse) AcceptStatus(code int) bool {
	return code == http.StatusNotFound
}

// UnmarshalJSON implements json.Unmarshaler.
func (exists *ResourceExistsResponse) UnmarshalJSON(_ []byte) error {
	*exists = true
	return nil
}

// UnmarshalMessage implements transport.MessageUnmarshaler.
func (exists *ResourceExistsResponse) UnmarshalMessage(r *proto.SearchReply) error {
	dev.Assert(exists != nil, "unmarshal called with nil receiver")
	dev.Assert(r != nil, "search reply is nil")
	*exists = len(r.Results) > 0
	return nil
}
