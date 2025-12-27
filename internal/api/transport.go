package api

import (
	"sync"

	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/transport"
)

// TransportFactory returns an internal.Transport instance for TransportConfig.
type TransportFactory func(TransportConfig) (internal.Transport, error)

var (
	tfMu sync.RWMutex                    // tfMu guards tf.
	tf   TransportFactory = newTransport // tf is the default TransportFactory.
)

// SetTransportFactory changes the transport factory for this package.
// This is a test helper and MUST NOT be used in the public layer.
func SetTransportFactory(newtf TransportFactory) {
	tfMu.Lock()
	defer tfMu.Unlock()
	tf = newtf
}

// GetTransportFactory returns the current transport factory.
// GetTransportFactory is exported for testing purposes in production.
// Use NewTransport to obtain a new transport.
func GetTransportFactory() TransportFactory {
	tfMu.RLock()
	defer tfMu.RUnlock()
	return tf
}

type TransportConfig transport.Config

// NewTransport returns internal.Transport for this API version.
func NewTransport(cfg TransportConfig) (internal.Transport, error) {
	newT := GetTransportFactory()
	return newT(cfg)
}

func newTransport(cfg TransportConfig) (internal.Transport, error) {
	cfg.Version = Version
	return transport.New(transport.Config(cfg))
}
