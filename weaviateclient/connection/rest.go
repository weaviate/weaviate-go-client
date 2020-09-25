package connection

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const apiVersion = "v1"

type Connection struct {
	basePath string
	httpClient *http.Client
}

func NewConnection(scheme string, host string) *Connection {
	return &Connection{
		basePath: scheme + "://" + host + "/" + apiVersion,
		httpClient: &http.Client{},
	}
}

func (con *Connection) addHeaderToRequest(request *http.Request) {
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
}

func (con *Connection) createRequest(ctx context.Context, path string, restMethod string) (*http.Request, error) {
	url := con.basePath+path

	request, err := http.NewRequest(restMethod, url, nil)
	if err != nil {
		return nil, err
	}
	con.addHeaderToRequest(request)
	request.WithContext(ctx)
	return request, nil
}

func (con *Connection) RunREST(ctx context.Context, path string, restMethod string) (*ResponseData, error) {
	request, requestErr := con.createRequest(ctx, path, restMethod)
	if requestErr != nil {
		return nil, requestErr
	}
	response, responseErr := con.httpClient.Do(request)
	if responseErr != nil {
		return nil, responseErr
	}

	defer response.Body.Close()
	body, bodyErr := ioutil.ReadAll(response.Body)
	if bodyErr != nil {
		return nil, bodyErr
	}

	return &ResponseData{
		Body:       body,
		StatusCode: response.StatusCode,
	}, nil
}

// ResponseData encapsulation of the http request body and status
type ResponseData struct {
	Body []byte
	StatusCode int
}

// DecodeBodyIntoTarget unmarshall body into target var
// successful if err is nil
func (rd *ResponseData) DecodeBodyIntoTarget(target interface{}) error {
	err := json.Unmarshal(rd.Body, target)
	if err != nil {
		return err
	}
	return nil
}

