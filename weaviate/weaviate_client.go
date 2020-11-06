package weaviate

import (
	"github.com/semi-technologies/weaviate-go-client/weaviate/batch"
	"github.com/semi-technologies/weaviate-go-client/weaviate/classifications"
	"github.com/semi-technologies/weaviate-go-client/weaviate/connection"
	"github.com/semi-technologies/weaviate-go-client/weaviate/contextionary"
	"github.com/semi-technologies/weaviate-go-client/weaviate/data"
	"github.com/semi-technologies/weaviate-go-client/weaviate/graphql"
	"github.com/semi-technologies/weaviate-go-client/weaviate/misc"
	"github.com/semi-technologies/weaviate-go-client/weaviate/schema"
)

// Config of the client endpoint
type Config struct {
	Host   string
	Scheme string
}

// Client implementing the weaviate API
type Client struct {
	connection      *connection.Connection
	misc            *misc.API
	schema          *schema.API
	data            *data.API
	batch           *batch.API
	c11y            *contextionary.API
	classifications *classifications.API
	graphQL         *graphql.API
}

// New weaviate client from config
func New(config Config) *Client {
	con := connection.NewConnection(config.Scheme, config.Host)

	return &Client{
		connection:      con,
		misc:            misc.New(con),
		schema:          schema.New(con),
		data:            data.New(con),
		batch:           batch.New(con),
		c11y:            contextionary.New(con),
		classifications: classifications.New(con),
		graphQL:         graphql.New(con),
	}
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
