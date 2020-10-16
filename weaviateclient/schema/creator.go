package schema

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	clientModels "github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type ClassCreator struct {
	connection   *connection.Connection
	class        *models.Class
	semanticKind clientModels.SemanticKind
}

func (cc *ClassCreator) WithClass(class *models.Class) *ClassCreator {
	cc.class = class
	return cc
}

func (cc *ClassCreator) WithKind(semanticKind clientModels.SemanticKind) *ClassCreator {
	cc.semanticKind = semanticKind
	return cc
}

func (cc *ClassCreator) Do(ctx context.Context) error {
	path := "/schema/"+string(cc.semanticKind)
	responseData, err := cc.connection.RunREST(ctx, path, http.MethodPost, cc.class)
	if err != nil {
		return err
	}
	if responseData.StatusCode == 200 {
		return nil
	}
	return clienterrors.NewUnexpectedStatusCodeError(responseData.StatusCode, string(responseData.Body))
}

