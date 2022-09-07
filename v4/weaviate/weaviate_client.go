package weaviate

import (
	"context"
	"net/http"

	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/batch"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/classifications"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/contextionary"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/data"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/graphql"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/misc"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/schema"
	"github.com/semi-technologies/weaviate-go-client/v4/weaviate/util"
)

// Config of the client endpoint
type Config struct {
	// Host of the weaviate instance; this is a mandatory field.
	Host string
	// Scheme of the weaviate instance; this is a mandatory field.
	Scheme string

	// ConnectionClient that will be used to execute http requests to the weaviate instance.
	//  If omitted a default will be used. The default is not able to handle authenticated requests.
	//
	//  To connect with an authenticated weaviate consider using the client from the golang.org/x/oauth2 module.
	ConnectionClient *http.Client

  // Headers added for every request
  Headers map[string]string
}

// Client implementing the weaviate API
// Every function represents one API group of weaviate and provides a set of functions and builders to interact with them.
//
// The client uses the original data models as provided by weaviate itself.
// All these models are provided in the sub module "github.com/semi-technologies/weaviate/entities/models"
type Client struct {
	connection        *connection.Connection
	misc              *misc.API
	schema            *schema.API
	data              *data.API
	batch             *batch.API
	c11y              *contextionary.API
	classifications   *classifications.API
	graphQL           *graphql.API
	version           string
	dbVersionProvider *util.DBVersionProvider
	dbVersionSupport  *util.DBVersionSupport
}

// New client from config
// Every function represents one API group of weaviate and provides a set of functions and builders to interact with them.
//
// The client uses the original data models as provided by weaviate itself.
// All these models are provided in the sub module "github.com/semi-technologies/weaviate/entities/models"
func New(config Config) *Client {
	con := connection.NewConnection(config.Scheme, config.Host, config.ConnectionClient, config.Headers)

	// some endpoints now require a className namespace.
	// to determine if this new convention is to be used,
	// we must check the weaviate server version
	getVersionFn := func() string {
		meta, err := misc.New(con, nil).MetaGetter().Do(context.Background())
		if err == nil {
			return meta.Version
		}
		return ""
	}

	dbVersionProvider := util.NewDBVersionProvider(getVersionFn)
	dbVersionSupport := util.NewDBVersionSupport(dbVersionProvider)

	client := &Client{
		connection:      con,
		misc:            misc.New(con, dbVersionProvider),
		schema:          schema.New(con),
		c11y:            contextionary.New(con),
		classifications: classifications.New(con),
		graphQL:         graphql.New(con),
		data:            data.New(con, dbVersionSupport),
		batch:           batch.New(con, dbVersionSupport),
	}

	return client
}

// Misc collection group for .well_known and root level API commands
func (c *Client) Misc() *misc.API {
	return c.misc
}

// Schema API group
func (c *Client) Schema() *schema.API {
	return c.schema
}

// Data API group including both things and actions
func (c *Client) Data() *data.API {
	return c.data
}

// Batch loading API group
func (c *Client) Batch() *batch.API {
	return c.batch
}

// C11y (contextionary) API group
func (c *Client) C11y() *contextionary.API {
	return c.c11y
}

// Classifications API group
func (c *Client) Classifications() *classifications.API {
	return c.classifications
}

// GraphQL API group
func (c *Client) GraphQL() *graphql.API {
	return c.graphQL
}
