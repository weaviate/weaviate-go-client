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

type (
	Config ssb.ClientConfig
	Option func(*ssb.ClientConfig)
)

func WithRetryFunc(f func(string, int, error) bool) Option {
	return func(c *ssb.ClientConfig) {
		c.RetryFunc = ssb.RetryFunc(f)
	}
}

func WithRetryTimes(n int) Option {
	return func(c *ssb.ClientConfig) {
		c.RetryFunc = func(_ string, retries int, _ error) bool {
			return n < retries
		}
	}
}

type ReconnectPolicy ssb.ReconnectPolicy

func WithReconnectPolicy(rp ReconnectPolicy) Option {
	return func(c *ssb.ClientConfig) {
		c.Reconnect = ssb.ReconnectPolicy(rp)
	}
}

const (
	queueSize   = 1 << 8
	batchSize   = 1 << 10
	reconnLimit = 5
)

func NewClient(ctx context.Context, t stream.Transport, rd api.RequestDefaults, options ...Option) *Client {
	conf := ssb.ClientConfig{
		Context:   ctx,
		Transport: stream.NewAdapter(t, rd),
		QueueSize: queueSize,
		BatchSize: batchSize,
		Reconnect: ssb.ReconnectPolicy{
			Limit:     reconnLimit,
			DelayFunc: ssb.ExponentialDelay,
		},
	}
	for _, opt := range options {
		opt(&conf)
	}
	return &Client{
		protocol: ssb.NewClient(conf),
	}
}

type Client struct {
	protocol *ssb.Client
}
