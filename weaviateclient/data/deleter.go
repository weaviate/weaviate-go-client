package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/paragons"
	"net/http"
)

type Deleter struct {
	connection   *connection.Connection
	uuid         string
	semanticKind paragons.SemanticKind
}

func (deleter *Deleter) WithID(uuid string) *Deleter {
	deleter.uuid = uuid
	return deleter
}

func (deleter *Deleter) WithKind(semanticKind paragons.SemanticKind) *Deleter {
	deleter.semanticKind = semanticKind
	return deleter
}

func (deleter *Deleter) Do(ctx context.Context) error {
	path := fmt.Sprintf("/%v/%v", deleter.semanticKind, deleter.uuid)
	responseData, err := deleter.connection.RunREST(ctx, path, http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if responseData.StatusCode == 204 {
		return nil
	}
	return clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}