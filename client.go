package weaviate

import (
	"github.com/weaviate/weaviate-go-client/v6/backup"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/internal"
)

func NewClient() (*Client, error) {
	t := internal.NewTransport()
	return &Client{
		Backup:      backup.NewClient(t),
		Collections: collections.NewClient(t),
	}, nil
}

type Client struct {
	Backup      *backup.Client
	Collections *collections.Client
}
