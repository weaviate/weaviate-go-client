package weaviate

import "github.com/weaviate/weaviate-go-client/v6/types"

// CollectionClient provides operations on a single collection
type CollectionClient struct {
	name        string
	client      *any // *Client
	tenant      string
	consistency types.ConsistencyLevel
}

// Query returns the query operations client
func (c *CollectionClient) Query() *any {
	return nil
}

// Data returns the data operations client
func (c *CollectionClient) Data() *DataClient {
	return &DataClient{collection: c}
}
