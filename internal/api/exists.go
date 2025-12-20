package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
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

// Bool returns bool value of ResourceExistsResponse.
func (hr ResourceExistsResponse) Bool() bool {
	return bool(hr)
}

var (
	_ json.Unmarshaler         = (*ResourceExistsResponse)(nil)
	_ transport.StatusAccepter = (*ResourceExistsResponse)(nil)
)

// AcceptStatus implements transport.StatusAccepter.
func (hr ResourceExistsResponse) AcceptStatus(code int) bool {
	return code == http.StatusNotFound
}

// UnmarshalJSON implements json.Unmarshaler.
func (hr *ResourceExistsResponse) UnmarshalJSON(_ []byte) error {
	*hr = true
	return nil
}
