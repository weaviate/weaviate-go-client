package query

import "github.com/weaviate/weaviate-go-client/v6/types"

func Scan[P types.Properties](r *Result) []types.Object[P] {
	return nil
}

func ScanGrouped[P types.Properties](r *GroupByResult) (map[string]Group[P], []GroupByObject[P]) {
	return nil, nil
}
