package batch

import (
	"context"

	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/ssb"
	"github.com/weaviate/weaviate-go-client/v6/internal/api/transport/stream"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

type Task struct{ t *ssb.Task }

// ID returns object ID, if the task adds an object,
// and a beacon if it adds a reference.
func (t *Task) ID() string {
	return t.t.ID()
}

// Wait blocks until the task completes. If it fails, Wait return a non-nil error.
//
// Every task is guaranteed to complete before [Client.Close] returns;
// until then it may take an arbitrary amount of retries or reconnects
// for the data to be successfully written to the server. Most of the time
// you will call Wait to collect task results after the stream is closed.
func (t *Task) Wait() error {
	<-t.t.Done()
	return t.t.Err()
}

type (
	Config ssb.ClientConfig
	Option func(*ssb.ClientConfig)
)

func WithRetryFunc(f func(string, int, error) bool) Option {
	return func(c *ssb.ClientConfig) {
		c.CanRetry = ssb.CanRetryFunc(f)
	}
}

// Retry each task up to n times.
func WithRetryTimes(n int) Option {
	return func(c *ssb.ClientConfig) {
		c.CanRetry = func(_ string, retries int, _ error) bool {
			return retries < n
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

func (c *Client) Object(ctx context.Context, o *data.Object) (*Task, error) {
	data := o.ToAPI()
	return c.add(ctx, ssb.Data{Object: &data})
}

func (c *Client) Reference(ctx context.Context, ref data.Reference) (*Task, error) {
	data := ref.ToAPI()
	return c.add(ctx, ssb.Data{Reference: &data})
}

func (c *Client) add(ctx context.Context, d ssb.Data) (*Task, error) {
	t, err := c.protocol.Add(ctx, d)
	if err != nil {
		return nil, err
	}
	return &Task{t: t}, nil
}

func (c *Client) Close() error {
	dev.AssertNotNil(c.protocol, "protocol")
	return c.protocol.Close()
}
