package request

import "github.com/weaviate/weaviate-go-client/v6/types"

type Defaults struct {
	CollectionName   string
	Tenant           string
	ConsistencyLevel types.ConsistencyLevel
}
