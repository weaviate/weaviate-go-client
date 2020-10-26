package misc

import (
	"context"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type MetaGetter struct {
	connection *connection.Connection
}

func (mg *MetaGetter) Do(ctx context.Context) (*models.Meta, error) {

	responseData, responseErr := mg.connection.RunREST(ctx, "/meta", http.MethodGet, nil)
	err := clienterrors.CheckResponnseDataErrorAndStatusCode(responseData, responseErr, 200)
	if err != nil {
		return nil, err
	}
	var meta models.Meta
	parseErr := responseData.DecodeBodyIntoTarget(&meta)
	return &meta, parseErr
}
