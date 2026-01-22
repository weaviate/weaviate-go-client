package api

// Version is the server API version this package defines requests for.
//
// For REST, it's the http://my-weaviate/v1.
// For gRPC, it's the stubs generated from v1 protobufs.
const Version = "v1"

// RequestDefaults holds options which are shared by
// all requests done by the same "collection handle".
type RequestDefaults struct {
	CollectionName   string
	Tenant           string
	ConsistencyLevel ConsistencyLevel
}

type ConsistencyLevel string

const (
	consistencyLevelUndefined                  = ""
	ConsistencyLevelOne       ConsistencyLevel = "ONE"
	ConsistencyLevelQuorum    ConsistencyLevel = "QUORUM"
	ConsistencyLevelAll       ConsistencyLevel = "ALL"
)
