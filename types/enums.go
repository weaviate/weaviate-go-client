package types

import (
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
)

// ConsistencyLevel defines supported read / write consistency level.
type ConsistencyLevel api.ConsistencyLevel

const (
	ConsistencyLevelOne    = ConsistencyLevel(api.ConsistencyLevelOne)
	ConsistencyLevelQuorum = ConsistencyLevel(api.ConsistencyLevelQuorum)
	ConsistencyLevelAll    = ConsistencyLevel(api.ConsistencyLevelAll)
)
