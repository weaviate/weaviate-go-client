package types

type Map map[string]any

type Properties interface {
	Map | any
}

type Object[P Properties] struct {
	UUID               string
	Properties         P
	Vectors            Vectors
	CreationTimeUnix   *int64
	LastUpdateTimeUnix *int64
}
