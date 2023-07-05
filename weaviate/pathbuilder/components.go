package pathbuilder

import "github.com/weaviate/weaviate-go-client/v4/weaviate/db"

type Components struct {
	ID                string
	Class             string
	DBVersion         *db.VersionSupport
	ConsistencyLevel  string
	Tenant            string
	ReferenceProperty string
}
