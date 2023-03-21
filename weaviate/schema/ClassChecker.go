package schema

import (
	"context"
)

// ClassDeleter builder to remove a class from weaviate
type ClassChecker struct {
	schemaAPI *API
	className string
}

// WithClassName defines the name of the class that should be checked
func (cd *ClassChecker) WithClassName(className string) *ClassChecker {
	cd.className = className
	return cd
}

// Do delete the class from the weaviate schema
func (cd *ClassChecker) Do(ctx context.Context) bool {
	_, err := cd.schemaAPI.ClassGetter().WithClassName(cd.className).Do(context.Background())
	return err == nil
}
