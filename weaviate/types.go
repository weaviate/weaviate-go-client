package weaviate

import "github.com/weaviate/weaviate-go-client/v6/types"

type wrappedVector types.Vector

// ToBytes implements internal.Vector.
func (v *wrappedVector) ToBytes() []byte {
	panic("unimplemented")
}

func (v *wrappedVector) IsMulti() bool {
	return len((*types.Vector)(v).Multi) > 0
}

func (v *wrappedVector) ToFloat32() []float32 {
	return (*types.Vector)(v).Single
}

func (v *wrappedVector) ToFloat32Multi() [][]float32 {
	return (*types.Vector)(v).Multi
}
