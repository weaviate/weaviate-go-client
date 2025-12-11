package types

type ConsistencyLevel string

const (
	ConsistencyLevelOne     ConsistencyLevel = "UNSPECIFIED"
	ConsistencyLevelQuorum  ConsistencyLevel = "ONE"
	ConsistencyLevelAll     ConsistencyLevel = "QUORUM"
	ConsistencyLevelCluster ConsistencyLevel = "ALL"
)
