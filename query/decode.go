package query

import (
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

// Decode decodes types.Map properties of query result objects into arbitrary Go structs.
func Decode[P types.Properties](r *Result) ([]Object[P], error) {
	out := make([]Object[P], len(r.Objects))
	for i, obj := range r.Objects {
		properties, err := internal.Decode[P](obj.Properties)
		if err != nil {
			return nil, fmt.Errorf("decode: %w", err)
		}
		out[i] = Object[P]{
			Object: types.Object[P]{
				UUID:               obj.UUID,
				Properties:         *properties,
				CreationTimeUnix:   obj.CreationTimeUnix,
				LastUpdateTimeUnix: obj.LastUpdateTimeUnix,
			},
			Metadata: obj.Metadata,
		}
	}
	return out, nil
}

// Decode decodes types.Map properties of grouped query result objects into arbitrary Go structs.
func DecodeGrouped[P types.Properties](r *GroupByResult) (map[string]Group[P], []GroupByObject[P], error) {
	groups := make(map[string]Group[P], len(r.Groups))
	objects := make([]GroupByObject[P], 0, len(r.Objects))
	for name, group := range r.Groups {
		for _, obj := range group.Objects {
			properties, err := internal.Decode[P](obj.Properties)
			if err != nil {
				return nil, nil, fmt.Errorf("decode grouped: %w", err)
			}
			objects = append(objects, GroupByObject[P]{
				BelongsToGroup: name,
				Object: Object[P]{
					Object: types.Object[P]{
						UUID:               obj.UUID,
						Properties:         *properties,
						CreationTimeUnix:   obj.CreationTimeUnix,
						LastUpdateTimeUnix: obj.LastUpdateTimeUnix,
					},
					Metadata: obj.Metadata,
				},
			})
		}

		// Create a view into the Objects slice rather than allocating a separate one.
		from, to := len(objects)-len(group.Objects), len(objects)-1
		groups[name] = Group[P]{
			Name:        name,
			MinDistance: group.MinDistance,
			MaxDistance: group.MaxDistance,
			Size:        group.Size,
			Objects:     objects[from:to],
		}
	}
	return groups, objects, nil
}
