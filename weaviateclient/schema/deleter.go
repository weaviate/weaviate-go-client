package schema

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/models"
	"net/http"
)

type ClassDeleter struct {
	connection *connection.Connection
	semanticKind models.SemanticKind
	className string
}

func (cd *ClassDeleter) WithClassName(className string) *ClassDeleter {
	cd.className = className
	return cd
}

func (cd *ClassDeleter) WithKind(semanticKind models.SemanticKind) *ClassDeleter {
	cd.semanticKind = semanticKind
	return cd
}

func (cd *ClassDeleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/schema/%v/%v", cd.semanticKind, cd.className)
	responseData, err := cd.connection.RunREST(ctx, path, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if responseData.StatusCode == 200 {
		return nil
	}
	return clienterrors.NewUnexpectedStatusCodeError(responseData.StatusCode, string(responseData.Body))
}

