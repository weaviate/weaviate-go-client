package query

type Client struct {
	gRPC any // gRPCClient

	NearVector NearVectorFunc
}

func NewClient(gRPC any) *Client {
	return &Client{
		gRPC:       gRPC,
		NearVector: nearVector,
	}
}

type CommonOptions struct {
	Limit            *int
	Offset           *int
	AutoLimit        *int
	After            string
	ReturnProperties []string
	IncludeVectors   []string
}

// LimitOption sets the `limit` parameter.
type LimitOption int

var _ NearVectorOption = (*LimitOption)(nil)

func WithLimit(l int) LimitOption {
	return LimitOption(l)
}

// OffsetOption sets the `limit` parameter.
type OffsetOption int

var _ NearVectorOption = (*OffsetOption)(nil)

func WithOffset(l int) OffsetOption {
	return OffsetOption(l)
}

// AutoLimitOption sets the `limit` parameter.
type AutoLimitOption int

var _ NearVectorOption = (*AutoLimitOption)(nil)

func WithAutoLimit(l int) AutoLimitOption {
	return AutoLimitOption(l)
}

type Result struct {
	Objects []map[string]any
}
