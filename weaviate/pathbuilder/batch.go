package pathbuilder

import "fmt"

func BatchObjects(comp Components) string {
	path := "/batch/objects"
	if comp.ConsistencyLevel != "" {
		path = fmt.Sprintf("%s?consistency_level=%v", path, comp.ConsistencyLevel)
	}
	return path
}

func BatchReferences(comp Components) string {
	path := "/batch/references"
	if comp.ConsistencyLevel != "" {
		path = fmt.Sprintf("%s?consistency_level=%v", path, comp.ConsistencyLevel)
	}
	return path
}
