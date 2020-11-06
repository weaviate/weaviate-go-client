package classifications

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviate/except"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
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
	err := except.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var classification models.Classification
	parseErr := responseData.DecodeBodyIntoTarget(&classification)
	return &classification, parseErr
}
