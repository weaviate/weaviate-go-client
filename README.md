# weaviate-go-client  <img alt='Weaviate logo' src='https://raw.githubusercontent.com/weaviate/weaviate/19de0956c69b66c5552447e84d016f4fe29d12c9/docs/assets/weaviate-logo.png' width='180' align='right' />

[![Tests](https://github.com/weaviate/weaviate-go-client/actions/workflows/tests.yaml/badge.svg?branch=v6)](https://github.com/weaviate/weaviate-go-client/actions/workflows/tests.yaml?query=branch%3Av6)
[![Go Reference](https://pkg.go.dev/badge/github.com/weaviate/weaviate-go-client/v6.svg)](https://pkg.go.dev/github.com/weaviate/weaviate-go-client/v6)

An idiomatic Go client for [Weaviate](https://github.com/weaviate/weaviate), the open-source
AI-native vector database.

> [!IMPORTANT]
> **This is the v6 pre-release line — the current tag is `v6.0.0-beta.1`.** v6 is a ground-up
> rewrite and its API can still change between pre-releases.
> **[v5](https://github.com/weaviate/weaviate-go-client/tree/main) remains the stable line for
> production use.**
>
> v6 does not yet cover the whole v5 feature set. Missing so far, among other things: generative
> (RAG) queries; reranking; sorting; standalone keyword (BM25) search, which exists only as the
> keyword half of hybrid search; near-object search; near-image and other multimodal search;
> fetching an object by ID; partial (merge) updates; deleting or replacing references; and changing
> an existing collection's configuration. Aggregation covers `OverAll` and `NearVector` only. On the
> configuration side there is one vector index (HFresh — no HNSW, flat or dynamic), one compression
> scheme (RQ — no PQ, BQ or SQ), and three vectorizers (`text2vec-weaviate`, `text2vec-model2vec`
> and self-provided vectors — no OpenAI, Cohere, Ollama and so on). This list is not exhaustive;
> check [pkg.go.dev](https://pkg.go.dev/github.com/weaviate/weaviate-go-client/v6) for what a
> package actually exposes.
>
> Some exported methods are known to be broken at `v6.0.0-beta.1`: `Data.Replace`,
> `Data.DeleteSelected`, `Tenants.Get`, and `Query.Hybrid` when given a `NearVector`. They compile
> and appear in the package reference, but fail or panic at runtime. Searches also fail
> intermittently with `invalid UUID (got 15 bytes)` when a returned object's server-assigned ID
> happens to begin with a zero byte — roughly one ID in 256; retrying, or assigning your own IDs,
> avoids it. Bug reports are very welcome — please
> [open an issue](https://github.com/weaviate/weaviate-go-client/issues).

## Requirements

- Go 1.25.8 or newer for the client itself. The programs under [`example/`](example) declare Go
  1.26.1, so building those needs a 1.26 toolchain.
- A Weaviate instance, either on [Weaviate Cloud](https://docs.weaviate.io/cloud) or
  [self-hosted](https://docs.weaviate.io/deploy/installation-guides/docker-installation). The REST
  contract vendored in [`api/rest`](api/rest) is generated from Weaviate `1.39.0`, and this
  pre-release has been exercised against Weaviate `1.38.0`; there is no declared minimum supported
  server version yet.
- Both the REST and the gRPC endpoint of that instance must be reachable. The client uses gRPC for
  search, aggregation, inserts and batch, and REST for schema, object replace/delete and
  administrative calls.

## Installation

```sh
go get github.com/weaviate/weaviate-go-client/v6@v6.0.0-beta.1
```

Pinning the version is deliberate. There is no stable `v6.0.0` tag yet, so a bare
`go get github.com/weaviate/weaviate-go-client/v6` resolves to the newest pre-release today too —
but it will move you onto `v6.0.0` without warning the moment that ships. Pin the version if you
would rather take that upgrade when you choose to.

The client lives at the module root and its package is named `weaviate`:

```go
import weaviate "github.com/weaviate/weaviate-go-client/v6"
```

## Getting started

The program below connects to Weaviate Cloud, creates a collection whose vectors are produced by
[Weaviate Embeddings](https://docs.weaviate.io/cloud/embeddings), imports three objects, and runs a
vector search over them.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	weaviate "github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/data"
	embeddings "github.com/weaviate/weaviate-go-client/v6/modules/weaviate"
	"github.com/weaviate/weaviate-go-client/v6/query"
)

type Movie struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Genre       string `json:"genre"`
}

func main() {
	ctx := context.Background()

	// 1. Connect. WEAVIATE_HOST is the cluster hostname without a scheme,
	//    e.g. "my-cluster.c0.europe-west3.gcp.weaviate.cloud".
	client, err := weaviate.NewWeaviateCloud(ctx,
		os.Getenv("WEAVIATE_HOST"),
		os.Getenv("WEAVIATE_API_KEY"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// 2. Create a collection. Vectors are named; here a single "default"
	//    vector is built from the description by Weaviate Embeddings.
	movies, err := client.Collections.Create(ctx, collections.Collection{
		Name: "Movie",
		Properties: []collections.Property{
			{Name: "title", DataType: collections.DataTypeText},
			{Name: "description", DataType: collections.DataTypeText},
			{Name: "genre", DataType: collections.DataTypeText},
		},
		Vectors: map[string]collections.VectorConfig{
			"default": {
				Vectorizer: embeddings.Text2Vec{
					Properties: []string{"description"},
					Model:      embeddings.SnowflakeArcticEmbedMv1_5,
					Dimensions: 256,
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 3. Insert objects. Insert is variadic: pass one object or many.
	inserted, err := movies.Data.Insert(ctx,
		&data.Object{Properties: data.MustEncode(&Movie{
			Title:       "The Matrix",
			Genre:       "Science Fiction",
			Description: "A hacker discovers that reality is a simulation and joins the war against the machines running it.",
		})},
		&data.Object{Properties: data.MustEncode(&Movie{
			Title:       "Spirited Away",
			Genre:       "Animation",
			Description: "A young girl is trapped in a world of spirits and must find a way to free her parents.",
		})},
		&data.Object{Properties: data.MustEncode(&Movie{
			Title:       "Arrival",
			Genre:       "Science Fiction",
			Description: "A linguist is recruited to communicate with extraterrestrial visitors and learns to perceive time differently.",
		})},
	)
	if err != nil {
		log.Fatal(err)
	}
	for id, msg := range inserted.Errors {
		log.Printf("object %s was rejected: %s", id, msg)
	}

	// 4. Search. The collection has exactly one vector, so the query resolves
	//    to it; with several vectors, name one with the Target field.
	found, err := movies.Query.NearText(ctx, query.NearText{
		Concepts:       []string{"simulated reality"},
		Limit:          2,
		ReturnMetadata: query.ReturnMetadata{Distance: true},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 5. Decode the results into your own type and print them.
	var results []query.Object[Movie]
	if err := query.Decode(found, &results); err != nil {
		log.Fatal(err)
	}
	for _, obj := range results {
		fmt.Printf("%-16s %-16s distance=%.4f\n",
			obj.Properties.Title, obj.Properties.Genre, *obj.Metadata.Distance)
	}
}
```

`Collections.Create` fails if the collection already exists, so delete it with
`client.Collections.Delete(ctx, "Movie")` — or pick a different name — before running this a second
time.

### Running against a local instance

The program above needs Weaviate Cloud, because Weaviate Embeddings is a Cloud service. Swapping in
`NewLocal` is not enough on its own: a self-hosted instance has no vectorizer configured, and v6
rejects a `VectorConfig{}` that does not name one. Use the `selfprovided` vectorizer and supply the
vectors yourself. This program runs as-is against
`docker run -p 8080:8080 -p 50051:50051 semitechnologies/weaviate`:

```go
package main

import (
	"context"
	"fmt"
	"log"

	weaviate "github.com/weaviate/weaviate-go-client/v6"
	"github.com/weaviate/weaviate-go-client/v6/collections"
	"github.com/weaviate/weaviate-go-client/v6/data"
	"github.com/weaviate/weaviate-go-client/v6/modules/selfprovided"
	"github.com/weaviate/weaviate-go-client/v6/query"
	"github.com/weaviate/weaviate-go-client/v6/types"
)

func main() {
	ctx := context.Background()

	client, err := weaviate.NewLocal(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Every named vector must declare a vectorizer — an empty VectorConfig{} is
	// rejected. selfprovided means "I supply the vectors myself".
	movies, err := client.Collections.Create(ctx, collections.Collection{
		Name:       "Movie",
		Properties: []collections.Property{{Name: "title", DataType: collections.DataTypeText}},
		Vectors: map[string]collections.VectorConfig{
			"default": {Vectorizer: selfprovided.Vectorizer},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// The collection's vector name, the name on each inserted vector and the
	// query target must all match.
	for title, vec := range map[string][]float32{
		"The Matrix":    {1, 0, 0},
		"Spirited Away": {0, 1, 0},
		"Arrival":       {0.9, 0.1, 0},
	} {
		if _, err := movies.Data.Insert(ctx, &data.Object{
			Properties: map[string]any{"title": title},
			Vectors:    []types.Vector{{Name: "default", Single: vec}},
		}); err != nil {
			log.Fatal(err)
		}
	}

	found, err := movies.Query.NearVector(ctx, query.NearVector{
		Target:         &types.Vector{Name: "default", Single: []float32{1, 0, 0}},
		Limit:          2,
		ReturnMetadata: query.ReturnMetadata{Distance: true},
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, obj := range found.Objects {
		fmt.Printf("%-16v distance=%.4f\n", obj.Properties["title"], *obj.Metadata.Distance)
	}
}
```

### Connecting

```go
// A local instance on the default ports: REST on localhost:8080, gRPC on localhost:50051.
client, err := weaviate.NewLocal(ctx)

// Weaviate Cloud, authenticated with an API key. Pass the cluster hostname only.
client, err := weaviate.NewWeaviateCloud(ctx, "my-cluster.weaviate.cloud", apiKey)

// Anything else: set the REST and gRPC endpoints yourself.
client, err := weaviate.NewClient(ctx,
	weaviate.WithScheme("https"),
	weaviate.WithHTTPHost("weaviate.example.com"),
	weaviate.WithHTTPPort(443),
	weaviate.WithGRPCHost("grpc.weaviate.example.com"),
	weaviate.WithGRPCPort(443),
	weaviate.WithAPIKey(apiKey),
)
```

A client owns a gRPC channel, so always `defer client.Close()`. `client.IsReady(ctx)` reports whether
the instance is serving.

API keys and OIDC tokens are sent as bearer credentials over gRPC, which requires transport-level
security, so authenticate against a TLS (`https`) endpoint. Passing `WithAPIKey`, `WithBearerToken`
or any other token source to a plaintext `http` endpoint currently fails while the client is being
constructed.

### More examples

The programs under [`example/`](example) are runnable end to end and double as journey tests in CI:

- [`example/basic`](example/basic) — collection setup, inserts, near-text and near-vector search,
  and grouped aggregation.
- [`example/batch`](example/batch) — streaming (server-side) batch import via `handle.Batch(ctx)`.
- [`example/rbac`](example/rbac) — roles, database users, OIDC groups and permissions.

See [`example/README.md`](example/README.md) for how to point them at a cluster.

## What changed in v6

v6 is a rewrite and is not source-compatible with v5. The API changed shape in five main ways:

- **Collections first.** Take a handle with `client.Collections.Use("Movie")` (or keep the one
  returned by `Collections.Create`), then reach data, search, aggregation and tenants through it:
  `handle.Data`, `handle.Query`, `handle.Aggregate`, `handle.Tenants`. Tenant and consistency level
  are bound once on the handle with `collections.WithTenant` / `collections.WithConsistencyLevel`
  instead of being repeated on every call.
- **Context first, no `.Do()`.** Every call is `f(ctx, params) (result, error)`. Request parameters
  are plain structs, so there is no builder chain and no terminating `.Do(ctx)`.
- **Named vectors by default.** A collection declares `Vectors map[string]VectorConfig`, and a search
  picks one with a `Target`. When a collection has exactly one vector, the target may be omitted and
  the server resolves it.
- **Grouped sub-clients.** Cluster-wide operations hang off the client as fields: `Collections`,
  `Roles`, `Users`, `Groups`, `Backup`, `Cluster`, `Replication` and `Alias`.
- **Typed results.** Results arrive as `query.Object[map[string]any]`; `query.Decode` (and
  `query.DecodeGrouped`) converts them into a slice of your own struct type using generics.

This is not a migration guide — see the [Weaviate documentation](https://docs.weaviate.io/weaviate)
for the concepts behind these calls.

## Documentation

- [Package reference on pkg.go.dev](https://pkg.go.dev/github.com/weaviate/weaviate-go-client/v6)
- [Weaviate documentation](https://docs.weaviate.io/weaviate)
- [Quickstart](https://docs.weaviate.io/weaviate/quickstart)
- [Connect to Weaviate](https://docs.weaviate.io/weaviate/connections)
- [Manage collections](https://docs.weaviate.io/weaviate/manage-collections)
- [Manage data objects](https://docs.weaviate.io/weaviate/manage-objects)
- [Search](https://docs.weaviate.io/weaviate/search)

The [Go client library page](https://docs.weaviate.io/weaviate/client-libraries/go) still documents
v5; the v6 pages are in progress. Until they land, `pkg.go.dev` and the [`example/`](example)
programs are the closest thing to a v6 reference — read them alongside the list of gaps and broken
methods at the top of this file, because the package reference does not distinguish them.

## Community

- [GitHub issues](https://github.com/weaviate/weaviate-go-client/issues) — bugs and feature requests
  for this client.
- [Weaviate forum](https://forum.weaviate.io/) — questions and discussion.

## Contributing

Read [CONTRIBUTE.md](CONTRIBUTE.md) before you start; it covers the development setup, the code
generation steps, and how the API contracts are kept in sync with the server. This project follows a
[Code of Conduct](CODE_OF_CONDUCT.md).

## Tests

The suite is made of unit tests. It needs no running Weaviate instance and no extra tooling:

```sh
TESTKIT_NO_WITHONLY=1 go test ./...
```

Set `TESTKIT_NO_WITHONLY`: it turns a stray `testkit.WithOnly` into an error instead of silently
dropping the rest of a table. CI sets it, so a leftover `WithOnly` that passes locally without it
will fail there. CI also runs `go vet ./...`, `golangci-lint`, and the suite under the race detector
(`go test -v -race -cover ./...`).

Code generation and running the linter locally need a few binaries, which are installed into a
git-ignored `./bin`:

```sh
go run ./cmd/build onboard   # installs protoc, golangci-lint and an oapi-codegen shim
```

The programs under `example/` are journey tests that run against a live cluster:

```sh
export WEAVIATE_HOST=<cluster_hostname>
export WEAVIATE_API_KEY=<api_key>
cd example && go run ./basic
```

## License

BSD 3-Clause — see [LICENSE](LICENSE).
