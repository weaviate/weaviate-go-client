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

func (getter *ActionGetter) Do(ctx context.Context) ([]*models.Action, error) {
	responseData, err := getObjectList(ctx, "/actions", getter.uuid, getter.connection)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 200 {
		return nil, clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
	}
	if getter.uuid == "" {
		var actions models.ActionsListResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&actions)
		return actions.Actions, decodeErr
	}
	var action models.Action
	decodeErr := responseData.DecodeBodyIntoTarget(&action)
	return []*models.Action{&action}, decodeErr
}

func (getter *ThingGetter) Do(ctx context.Context) ([]*models.Thing, error) {
	responseData, err := getObjectList(ctx, "/things", getter.uuid, getter.connection)
	if err != nil {
		return nil, err
	}
	if responseData.StatusCode != 200 {
		return nil, clienterrors.NewUnexpectedStatusCodeErrorFromRESTResponse(responseData)
	}
	if getter.uuid == "" {
		var things models.ThingsListResponse
		decodeErr := responseData.DecodeBodyIntoTarget(&things)
		return things.Things, decodeErr
	}
	var thing models.Thing
	decodeErr := responseData.DecodeBodyIntoTarget(&thing)
	return []*models.Thing{&thing}, decodeErr
}

func getObjectList(ctx context.Context, basePath string, uuid string, con *connection.Connection) (*connection.ResponseData, error) {
	path := basePath
	if uuid != "" {
		path += fmt.Sprintf("/%v", uuid)
	}
	responseData, err := con.RunREST(ctx, path, http.MethodGet, nil)
	return responseData, err
}

