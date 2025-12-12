package internal

import "github.com/weaviate/weaviate-go-client/v6/types"

type RequestDefaults struct {
	CollectionName   string
	Tenant           string
	ConsistencyLevel types.ConsistencyLevel
}
