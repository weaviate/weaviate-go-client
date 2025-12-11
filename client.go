package weaviate

import (
	"github.com/weaviate/weaviate-go-client/v5/collections"
)

func NewClient() (*Client, error) {
	return &Client{
		Collections: *collections.NewClient("tramway"),
	}, nil
}

type Client struct {
	Collections collections.Client
}
