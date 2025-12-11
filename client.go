package weaviate

import (
	"github.com/weaviate/weaviate-go-client/v5/query"
)

func NewClient() (*Client, error) {
	return &Client{
		Query: *query.NewClient(nil),
	}, nil
}

type Client struct {
	Query query.Client
}
