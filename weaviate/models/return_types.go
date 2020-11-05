package models

// OpenIDConfiguration of weaviate
type OpenIDConfiguration struct {
	// The Location to redirect to
	Href string `json:"href,omitempty"`
	// OAuth Client ID
	ClientID string `json:"clientId,omitempty"`
}

// SchemaDump Contains all semantic types and respective classes of the schema
type SchemaDump struct {
	Things  *Schema `json:"things"`
	Actions *Schema `json:"actions"`
}

// ThingsBatchRequestBody wrapping things to a batch
type ThingsBatchRequestBody struct {
	Fields []string        `json:"fields"`
	Things []*Thing `json:"things"`
}

// ActionsBatchRequestBody wrapping actions to a batch
type ActionsBatchRequestBody struct {
	Fields  []string         `json:"fields"`
	Actions []*Action `json:"actions"`
}
