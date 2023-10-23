package grpc

type Config struct {
	Enabled bool
	// Host of the weaviate instance; this is a mandatory field.
	Host string
	// Scheme of the weaviate instance; this is a mandatory field.
	Scheme string
}
