package pathbuilder

import "github.com/weaviate/weaviate-go-client/v5/weaviate/db"

type Components struct {
	ID                string
	Class             string
	DBVersion         *db.VersionSupport
	ConsistencyLevel  string
	Tenant            string
	ReferenceProperty string
}
