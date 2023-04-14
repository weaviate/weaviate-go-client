package pathbuilder

import (
	"fmt"
	"net/url"
)

func BatchObjects(comp Components) string {
	path := "/batch/objects"
	if comp.ConsistencyLevel != "" {
		pathParams := url.Values{}
		pathParams.Set("consistency_level", comp.ConsistencyLevel)
		path = fmt.Sprintf("%s?%v", path, pathParams.Encode())
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
