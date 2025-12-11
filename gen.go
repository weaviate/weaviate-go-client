package weaviate

//go:generate go tool github.com/go-swagger/go-swagger/cmd/swagger generate client -f contracts/schema.json -c rest -t weaviate/internal/clients/rest -A weaviate;
