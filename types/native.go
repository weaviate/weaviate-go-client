package types

import "github.com/weaviate/weaviate-go-client/v6/internal/api"

// NOTE(dyma): Below we are doing what's called "type aliasing".
// While it's not meant to be used for exposing internal types,
// we will allow ourselves this liberty this once.
//
// When the user passes object properties as a struct via [data.Encode],
// mapstructure will internally encode GeoCoordinates before we hand the
// data over to 'internal/api'.
//
// However, if they pass a properties map like this instead:
// 	map[string]any{"location": types.GeoCoordinates{
// 		Latitude:  -37.815389,
// 		Longitude: 144.970806,
// 	}}
// 'internal/api' will see a public type it doesn't know how to handle.
// Checking every map value before handing over is probably wasteful: it will
// have to happen on every insert, and most of the data will have no coordinates.
// Same for PhoneNumber.
//
// See: https://go.dev/blog/alias-names

// GeoCoordinates describe a coordinate pair.
type GeoCoordinates = api.GeoCoordinates
