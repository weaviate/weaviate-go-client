package classifications

import (
	"context"
	"fmt"
	"net/http"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/except"
	"github.com/weaviate/weaviate/entities/models"
)

// Getter builder to retrieve a classification status object
type Getter struct {
	connection *connection.Connection
	withID     string
}

// WithID of the classification
func (g *Getter) WithID(uuid string) *Getter {
	g.withID = uuid
	return g
}

// Do get the classification
func (g *Getter) Do(ctx context.Context) (*models.Classification, error) {
	path := fmt.Sprintf("/classifications/%v", g.withID)
	responseData, responseErr := g.connection.RunREST(ctx, path, http.MethodGet, nil)
	err := except.CheckResponseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var classification models.Classification
	parseErr := responseData.DecodeBodyIntoTarget(&classification)
	return &classification, parseErr
}
