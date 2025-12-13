package api

type RequestDefaults struct {
	CollectionName   string
	Tenant           string
	ConsistencyLevel ConsistencyLevel
}
