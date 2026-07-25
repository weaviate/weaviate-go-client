package batch

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport/stream"
)

type Task struct{ *ssb.Task }

func (t *Task) Wait() error {
	<-t.Done()
	return t.Err()
}

func NewClient(ctx context.Context, t stream.Transport, rd api.RequestDefaults) (*Client, error) {
	tad := stream.NewAdapter(t, rd)
	c, err := ssb.NewClient(ssb.ClientConfig{
		Context:   ctx,
		Transport: tad,
		// TODO(dyma): fill out the rest of the config
	})
	if err != nil {
		return nil, err
	}
	return &Client{
		protocol: c,
	}, nil
}

type Client struct {
	protocol *ssb.Client
}
