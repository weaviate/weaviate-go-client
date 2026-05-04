package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

var (
	IsLiveRequest  = transports.StaticEndpoint(http.MethodGet, "/.well-known/live")
	IsReadyRequest = transports.StaticEndpoint(http.MethodGet, "/.well-known/ready")
)

var (
	GetInstanceMetadataRequest                  = transport.GetInstanceMetadataRequest
	_                          json.Unmarshaler = (*GetInstanceMetadataResponse)(nil)
)

type GetInstanceMetadataResponse transport.GetInstanceMetadataResponse

func (r *GetInstanceMetadataResponse) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*transport.GetInstanceMetadataResponse)(r))
}
