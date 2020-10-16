package schema

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	clientModels "github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type PropertyCreator struct {
	connection   *connection.Connection
	semanticKind clientModels.SemanticKind
	className    string
	property     models.Property
}

func (pc *PropertyCreator) WithClassName(className string) *PropertyCreator {
	pc.className = className
	return pc
}

func (pc *PropertyCreator) WithProperty(property models.Property) *PropertyCreator {
	pc.property = property
	return pc
}

func (pc *PropertyCreator) WithKind(semanticKind clientModels.SemanticKind) *PropertyCreator {
	pc.semanticKind = semanticKind
	return pc
}

func (pc *PropertyCreator) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/%v/properties", string(pc.semanticKind), pc.className)
	responseData, err := pc.connection.RunREST(ctx, path, http.MethodPost, pc.property)
	if err != nil {
		return err
	}
	if responseData.StatusCode == 200 {
		return nil
	}
	return clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}