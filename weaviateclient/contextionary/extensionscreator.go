package contextionary

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

// ExtensionCreator builder for the weaviate contextionary
type ExtensionCreator struct {
	connection *connection.Connection
	extension *models.C11yExtension
}

// WithConcept a new concept that should be added or an existing concept that should be changed
func (ec *ExtensionCreator) WithConcept(concept string) *ExtensionCreator {
	ec.extension.Concept = concept
	return ec
}

// WithDefinition for the concept
func (ec *ExtensionCreator) WithDefinition(definition string) *ExtensionCreator {
	ec.extension.Definition = definition
	return ec
}

// WithWeight this new concept should be considered over a preexisting one
func (ec *ExtensionCreator) WithWeight(weight float32) *ExtensionCreator {
	ec.extension.Weight = weight
	return ec
}

// Do create the concept
func (ec *ExtensionCreator) Do(ctx context.Context) error {
	if ec.extension.Weight > 1.0 || ec.extension.Weight < 0.0 {
		return fmt.Errorf("weight must be between 0.0 and 1.0")
	}
	responseData, responseErr := ec.connection.RunREST(ctx, "/c11y/extensions", http.MethodPost, ec.extension)
	return clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
}