package api

import (
	"encoding/json"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
	proto "github.com/weaviate/weaviate/grpc/generated/protocol/v1"

	"github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/rest"
	"github.com/weaviate/weaviate-go-client/v6/internal/transports"
)

type Tenant struct {
	Name   string
	Status TenantStatus
}

var _ json.Marshaler = (*Tenant)(nil)

type TenantStatus rest.TenantActivityStatus

const (
	TenantStatusActive     = TenantStatus(rest.ACTIVE)
	TenantStatusCold       = TenantStatus(rest.COLD)
	TenantStatusFreezing   = TenantStatus(rest.FREEZING)
	TenantStatusFrozen     = TenantStatus(rest.FROZEN)
	TenantStatusHot        = TenantStatus(rest.HOT)
	TenantStatusInactive   = TenantStatus(rest.INACTIVE)
	TenantStatusOffloaded  = TenantStatus(rest.OFFLOADED)
	TenantStatusOffloading = TenantStatus(rest.OFFLOADING)
	TenantStatusOnloading  = TenantStatus(rest.ONLOADING)
	TenantStatusUnfreezing = TenantStatus(rest.UNFREEZING)
)

// CreateTenantsRequest creates new tenants in the collection.
type CreateTenantsRequest struct {
	transports.BaseEndpoint

	Collection string
	Tenants    []Tenant
}

var _ transports.Endpoint = (*CreateTenantsRequest)(nil)

func (CreateTenantsRequest) Method() string  { return http.MethodPost }
func (r *CreateTenantsRequest) Path() string { return "/schema/" + r.Collection + "/tenants" }
func (r *CreateTenantsRequest) Body() any    { return r.Tenants }

// UpdateTenantsRequest sets new statuses for the existing tenants.
type UpdateTenantsRequest struct {
	transports.BaseEndpoint

	Collection string
	Tenants    []Tenant
}

var _ transports.Endpoint = (*UpdateTenantsRequest)(nil)

func (UpdateTenantsRequest) Method() string  { return http.MethodPut }
func (r *UpdateTenantsRequest) Path() string { return "/schema/" + r.Collection + "/tenants" }
func (r *UpdateTenantsRequest) Body() any    { return r.Tenants }

// DeleteTenantsRequest batch deletes tenants by their name.
type DeleteTenantsRequest struct {
	transports.BaseEndpoint

	Collection string
	Tenants    []string
}

var _ transports.Endpoint = (*DeleteTenantsRequest)(nil)

func (DeleteTenantsRequest) Method() string  { return http.MethodDelete }
func (r *DeleteTenantsRequest) Path() string { return "/schema/" + r.Collection + "/tenants" }
func (r *DeleteTenantsRequest) Body() any    { return r.Tenants }

func (t *Tenant) MarshalJSON() ([]byte, error) {
	return json.Marshal(rest.Tenant{
		Name:           t.Name,
		ActivityStatus: rest.TenantActivityStatus(t.Status),
	})
}

// GetTenantsRequest retrieves tenant information for selected tenants.
// Leave Tenants field unset to fetch all tenants.
// Use with [GetTenantsResponse].
type GetTenantsRequest struct {
	Collection string
	Tenants    []string
}

var (
	_ transport.Message[proto.TenantsGetRequest, proto.TenantsGetReply] = (*GetTenantsRequest)(nil)
	_ transport.MessageMarshaler[proto.TenantsGetRequest]               = (*GetTenantsRequest)(nil)
)

func (r *GetTenantsRequest) Method() transport.MethodFunc[proto.TenantsGetRequest, proto.TenantsGetReply] {
	return proto.WeaviateClient.TenantsGet
}

func (r *GetTenantsRequest) Body() transport.MessageMarshaler[proto.TenantsGetRequest] {
	return r
}

func (r *GetTenantsRequest) MarshalMessage() (*proto.TenantsGetRequest, error) {
	req := &proto.TenantsGetRequest{
		Collection: r.Collection,
	}

	if len(r.Tenants) > 0 {
		req.Params = &proto.TenantsGetRequest_Names{
			Names: &proto.TenantNames{
				Values: r.Tenants,
			},
		}
	}
	return req, nil
}

type GetTenantsResponse []Tenant

var _ transport.MessageUnmarshaler[proto.TenantsGetReply] = (*GetTenantsResponse)(nil)

func (r *GetTenantsResponse) UnmarshalMessage(reply *proto.TenantsGetReply) error {
	dev.AssertNotNil(reply, "reply")
	if len(reply.Tenants) == 0 {
		return nil
	}

	*r = make([]Tenant, len(reply.Tenants))
	for i, t := range reply.Tenants {
		var status TenantStatus
		switch t.ActivityStatus {
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_HOT:
			status = TenantStatusHot
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_COLD:
			status = TenantStatusCold
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_FROZEN:
			status = TenantStatusFrozen
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_UNFREEZING:
			status = TenantStatusUnfreezing
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_FREEZING:
			status = TenantStatusFreezing
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_ACTIVE:
			status = TenantStatusActive
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_INACTIVE:
			status = TenantStatusInactive
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_OFFLOADED:
			status = TenantStatusOffloaded
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_OFFLOADING:
			status = TenantStatusOffloading
		case proto.TenantActivityStatus_TENANT_ACTIVITY_STATUS_ONLOADING:
			status = TenantStatusOnloading
		}

		(*r)[i] = Tenant{
			Name:   t.Name,
			Status: status,
		}
	}
	return nil
}
