package api

const Version = "v1"

type RequestDefaults struct {
	CollectionName   string
	Tenant           string
	ConsistencyLevel ConsistencyLevel
}
