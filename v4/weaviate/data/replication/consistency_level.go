package replication

var ConsistencyLevel = struct {
	ALL    string
	ONE    string
	QUORUM string
}{
	ALL:    "ALL",
	ONE:    "ONE",
	QUORUM: "QUORUM",
}
