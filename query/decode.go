package query

import (
	"fmt"
	"slices"

	"github.com/weaviate/weaviate-go-client/v6/internal"
)

// Decode decodes map[string]any properties of [query.Result] objects into arbitrary Go structs.
func Decode[P any](r *Result, dest *[]Object[P]) error {
	*dest = slices.Grow(*dest, len(r.Objects))

	for i, obj := range r.Objects {
		if i > len(*dest)-1 {
			*dest = append(*dest, *new(Object[P]))
		}
		d := (*dest)[i]
		if err := decode(&obj, &d); err != nil {
			return err
		}
		(*dest)[i] = d
	}

	return nil
}

// DecodeGrouped decodes map[string]any properties of [query.GroupByResult] objects into arbitrary Go structs.
func DecodeGrouped[P any](r *GroupByResult, dest *[]GroupObject[P]) (map[string]Group[P], error) {
	groups := make(map[string]Group[P], len(r.Groups)) // TODO(dyma): use internal.MakeMap
	*dest = slices.Grow(*dest, len(r.Objects))

	var tail int
	for _, group := range r.Groups {
		for _, obj := range group.Objects {
			if tail > len(*dest)-1 {
				*dest = append(*dest, *new(GroupObject[P]))
			}
			d := (*dest)[tail]
			if err := decode(&obj.Object, &d.Object); err != nil {
				return nil, err
			}
			d.BelongsToGroup = group.Name
			(*dest)[tail] = d
			tail++
		}

		// Create a view into the Objects slice rather than allocating a separate one.
		from, to := tail-len(group.Objects), tail
		if len(group.Objects) == 0 {
			to = from
		}
		groups[group.Name] = Group[P]{
			Name:        group.Name,
			MinDistance: group.MinDistance,
			MaxDistance: group.MaxDistance,
			Size:        group.Size,
			Objects:     (*dest)[from:to],
		}
	}
	return groups, nil
}

func decode[P any](src *Object[map[string]any], dest *Object[P]) error {
	err := internal.Decode(src.Properties, &dest.Properties)
	if err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	dest.UUID = src.UUID
	dest.CreatedAt = src.CreatedAt
	dest.LastUpdatedAt = src.LastUpdatedAt
	dest.Metadata = src.Metadata
	return nil
}
