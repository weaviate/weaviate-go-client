package types

// ConsistencyLevel defines supported read / write consistency level.
type ConsistencyLevel string

const (
	ConsistencyLevelOne     ConsistencyLevel = "UNSPECIFIED"
	ConsistencyLevelQuorum  ConsistencyLevel = "ONE"
	ConsistencyLevelAll     ConsistencyLevel = "QUORUM"
	ConsistencyLevelCluster ConsistencyLevel = "ALL"
)
