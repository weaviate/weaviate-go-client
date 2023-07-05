package pathbuilder

import (
	"fmt"
	"net/url"
)

func BatchObjects(comp Components) string {
	path := "/batch/objects"
	queryParams := url.Values{}

	if comp.ConsistencyLevel != "" {
		queryParams.Set("consistency_level", comp.ConsistencyLevel)
	}
	if comp.Tenant != "" {
		queryParams.Set("tenant", comp.Tenant)
	}
	if len(queryParams) > 0 {
		path = fmt.Sprintf("%s?%v", path, queryParams.Encode())
	}

	return path
}

func BatchReferences(comp Components) string {
	path := "/batch/references"
	pathParams := url.Values{}

	if comp.ConsistencyLevel != "" {
		pathParams.Set("consistency_level", comp.ConsistencyLevel)
	}
	if len(pathParams) > 0 {
		path = fmt.Sprintf("%s?%v", path, pathParams.Encode())
	}

	return path
}
