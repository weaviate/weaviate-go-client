package paragons

import "github.com/semi-technologies/weaviate/entities/models"

// OpenIDConfiguration of weaviate
type OpenIDConfiguration struct {
	// The Location to redirect to
	Href string `json:"href,omitempty"`
	// OAuth Client ID
	ClientID string `json:"clientId,omitempty"`
}

// SchemaDump Contains all semantic types and respective classes of the schema
type SchemaDump struct {
	Things  *models.Schema `json:"things"`
	Actions *models.Schema `json:"actions"`
}

// ThingsBatchRequestBody wrapping things to a batch
type ThingsBatchRequestBody struct {
	Fields []string        `json:"fields"`
	Things []*models.Thing `json:"things"`
}

// ActionsBatchRequestBody wrapping actions to a batch
type ActionsBatchRequestBody struct {
	Fields  []string         `json:"fields"`
	Actions []*models.Action `json:"actions"`
}
