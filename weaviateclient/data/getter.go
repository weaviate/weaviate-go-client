package data

import (
	"context"
	"fmt"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/clienterrors"
	"github.com/semi-technologies/weaviate-go-client/weaviateclient/connection"
	"github.com/semi-technologies/weaviate/entities/models"
	"net/http"
)

type ActionGetter struct {
	connection *connection.Connection
	uuid string
}

type ThingGetter struct {
	connection *connection.Connection
	uuid string
}

func (getter *ActionGetter) WithID(uuid string) *ActionGetter {
	getter.uuid = uuid
	return getter
}

func (getter *ThingGetter) WithID(uuid string) *ThingGetter {
	getter.uuid = uuid
	return getter
}

func (getter *ActionGetter) Do(ctx context.Context) (*models.Action, error) {
	path := fmt.Sprintf("/actions/%v", getter.uuid)
	responseData, err := getter.connection.RunREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode == 200 {
		var action models.Action
		decodeErr := responseData.DecodeBodyIntoTarget(&action)
		return &action, decodeErr
	}
	return nil, clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}

func (getter *ThingGetter) Do(ctx context.Context) (*models.Thing, error) {
	path := fmt.Sprintf("/things/%v", getter.uuid)
	responseData, err := getter.connection.RunREST(ctx, path, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode == 200 {
		var thing models.Thing
		decodeErr := responseData.DecodeBodyIntoTarget(&thing)
		return &thing, decodeErr
	}
	return nil, clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
}



