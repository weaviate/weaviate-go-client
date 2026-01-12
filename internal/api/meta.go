package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

var IsLiveRequest = transports.StaticEndpoint(http.MethodGet, "/.well-known/live")

var GetInstanceMetadataRequest = transports.StaticEndpoint(http.MethodGet, "/meta")

type GetInstanceMetadataResponse struct {
	Hostname           string
	Version            string
	Modules            map[string]any
	GRPCMaxMessageSize int
}

var _ json.Unmarshaler = (*GetInstanceMetadataResponse)(nil)

// UnmarshalJSON implements json.Unmarshaler.
func (r *GetInstanceMetadataResponse) UnmarshalJSON(data []byte) error {
	var meta rest.Meta
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}
	*r = GetInstanceMetadataResponse{
		Hostname:           meta.Hostname,
		Version:            meta.Version,
		Modules:            meta.Modules,
		GRPCMaxMessageSize: meta.GrpcMaxMessageSize,
	}
	return nil
}
