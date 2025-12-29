package weaviate

import (
	"github.com/weaviate/weaviate-go-client/v6/collections"
)

func NewClient() (*Client, error) {
	return &Client{
		Collections: *collections.NewClient(nil),
	}, nil
}

type Client struct {
	Collections collections.Client
}
