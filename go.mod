module github.com/weaviate/weaviate-go-client/v5

go 1.24.10

retract (
	v5.4.0 // Missing entry for golang.org/x/oauth2 in go.sum
	v5.0.0 // Malformed go.mod declares v4 module path
)
