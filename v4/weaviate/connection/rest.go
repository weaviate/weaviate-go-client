package connection

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/weaviate/weaviate-go-client/v4/weaviate/fault"
	"golang.org/x/oauth2"
)

const apiVersion = "v1"

// Connection networking layer accessing weaviate using http requests
type Connection struct {
	basePath   string
	httpClient *http.Client
	headers    map[string]string
	doneCh     chan bool
}

func finalizer(c *Connection) {
	c.doneCh <- true
}

// NewConnection based on scheme://host
// if httpClient is nil a default client will be used
func NewConnection(scheme string, host string, httpClient *http.Client, headers map[string]string) *Connection {
	client := httpClient
	if client == nil {
		client = &http.Client{}
	}
	connection := &Connection{
		basePath:   scheme + "://" + host + "/" + apiVersion,
		httpClient: client,
		headers:    headers,
		doneCh:     make(chan bool),
	}

	// shutdown goroutine when connections is cleaned up
	runtime.SetFinalizer(connection, finalizer)

	transport, ok := connection.httpClient.Transport.(*oauth2.Transport)
	if ok {
		connection.startRefreshGoroutine(transport)
	}

	return connection
}

// startRefreshGoroutine starts a background goroutine that periodically refreshes the auth token.
// The oauth2 package only refreshes the Tokens on new http requests => if there is no request for the lifetime of
// the refresh token the client will become de-authenticated without this.
func (con *Connection) startRefreshGoroutine(transport *oauth2.Transport) {
	token, err := transport.Source.Token()
	if err != nil {
		log.Printf("Error during token refresh, getting token: %v", err)
		return
	}
	// there is no point in manual refreshing if there is no refresh token. Note that this is the default with client
	// credentials
	if token.RefreshToken == "" {
		return
	}

	timeToSleep := token.Expiry.Sub(time.Now()) - time.Second/10
	if timeToSleep <= 0 {
		return
	}
	ticker := time.NewTicker(timeToSleep)
	go func() {
		for {
			select {
			case <-con.doneCh:
				return
			case <-ticker.C:
				_, err = con.RunREST(context.TODO(), "/meta", http.MethodGet, nil)
				if err != nil {
					log.Printf("Error during token refresh, rest request: %v", err)
					return
				}
			}
		}
	}()

}

func (con *Connection) addHeaderToRequest(request *http.Request) {
	for k, v := range con.headers {
		request.Header.Add(k, v)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
}

func (con *Connection) marshalBody(body interface{}) (io.Reader, error) {
	if body == nil {
		return nil, nil
	}
	jsonBody, err := json.Marshal(body) // Create the JSON body
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonBody), nil
}

func (con *Connection) createRequest(ctx context.Context, path string,
	restMethod string, body interface{},
) (*http.Request, error) {
	url := con.basePath + path // Create the URL

	jsonBody, err := con.marshalBody(body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(restMethod, url, jsonBody)
	if err != nil {
		return nil, err
	}
	con.addHeaderToRequest(request)
	request.WithContext(ctx)
	return request, nil
}

// RunREST executes a http request
// path: expects a resource path e.g. `/schema/things`
// restMethod: as they are defined in constants in the *http* package
// Returns:
//
//	a response that may be parsed into a struct after the fact
//	error if there was a network issue
func (con *Connection) RunREST(ctx context.Context, path string,
	restMethod string, requestBody interface{},
) (*ResponseData, error) {
	request, requestErr := con.createRequest(ctx, path, restMethod, requestBody)
	if requestErr != nil {
		return nil, requestErr
	}
	response, responseErr := con.httpClient.Do(request)
	if responseErr != nil {
		return nil, responseErr
	}

	defer response.Body.Close()
	body, bodyErr := io.ReadAll(response.Body)
	if bodyErr != nil {
		return nil, bodyErr
	}

	return &ResponseData{
		Body:       body,
		StatusCode: response.StatusCode,
	}, nil
}

func (con *Connection) RunRESTExternal(ctx context.Context, hostAndPath string, restMethod string, requestBody interface{}) (*ResponseData, error) {
	jsonBody, err := con.marshalBody(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(restMethod, hostAndPath, jsonBody)
	if err != nil {
		return nil, err
	}
	con.addHeaderToRequest(request)
	request.WithContext(ctx)
	response, responseErr := con.httpClient.Do(request)
	if responseErr != nil {
		return nil, responseErr
	}

	defer response.Body.Close()
	body, bodyErr := io.ReadAll(response.Body)
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
	Body       []byte
	StatusCode int
}

// DecodeBodyIntoTarget unmarshall body into target var
// successful if err is nil
func (rd *ResponseData) DecodeBodyIntoTarget(target interface{}) error {
	err := json.Unmarshal(rd.Body, target)
	if err != nil {
		return &fault.WeaviateClientError{
			IsUnexpectedStatusCode: false,
			StatusCode:             -1,
			Msg:                    "failed to parse resonse data check DerivedFromError field for more information",
			DerivedFromError:       err,
		}
	}
	return nil
}
