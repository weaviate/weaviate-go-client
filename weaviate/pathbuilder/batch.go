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
	if comp.TenantKey != "" {
		queryParams.Set("tenant_key", comp.TenantKey)
	}
	if len(queryParams) > 0 {
		path = fmt.Sprintf("%s?%v", path, queryParams.Encode())
	}

	return path
}

func BatchReferences(comp Components) string {
	path := "/batch/references"
	if comp.ConsistencyLevel != "" {
		pathParams := url.Values{}
		pathParams.Set("consistency_level", comp.ConsistencyLevel)
		path = fmt.Sprintf("%s?%v", path, pathParams.Encode())
	}
	return path
}
