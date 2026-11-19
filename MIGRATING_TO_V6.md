# Migrate your application from Go client v5 to v6

The v5 Go client makes heavy use of `GraphQL()` builders. However, Weaviate's [GraphQL API is deprecated](https://weaviate.io/blog/graphql-api-deprecation).

Thus, we strongly recommend upgrading every application entry point to
`github.com/weaviate/weaviate-go-client/v6`: services, import commands, scheduled
jobs and maintenance tools. That means every feature - rather than, for example,
updating queries alone and leaving ingestion on v5.

See the official [GraphQL migration guide](https://docs.weaviate.io/weaviate/api/graphql/migration)
for the retirement schedule, query mappings and
[features without an equivalent](https://docs.weaviate.io/weaviate/api/graphql/migration#graphql-features-without-an-equivalent).

## Recommendations:

### Before changing code

- Record representative responses, ordering, filters, pagination, counts and
  object IDs. Preserve those contracts rather than accepting new SDK defaults.
- Keep the existing embedding provider, model, source-property selection,
  collection configuration and ID derivation. An SDK upgrade is not a reason to
  switch vectorizers or recreate data.
- For the v6 client, use the version of as declared in its go.mod. For example **Go 1.26.10+** is specified in [v6.0.0](https://github.com/weaviate/weaviate-go-client/blob/v6.0.0/go.mod).
- Consider upgrading Weaviate to the latest version also, which you can find at: https://github.com/weaviate/weaviate/tags
- Make both HTTP and gRPC endpoints reachable. Client construction and
  `client.IsReady` check REST readiness only; exercise a real data/query call
  to verify gRPC. Follow the [Go connection requirements](https://docs.weaviate.io/weaviate/client-libraries/go).

Install v6 in the application module:

```sh
go get github.com/weaviate/weaviate-go-client/v6
```

The import is at the module root, without a `/weaviate` suffix:
`weaviate "github.com/weaviate/weaviate-go-client/v6"`.

## Migrate every application command

Use the [official Go documentation](https://docs.weaviate.io/weaviate/client-libraries/go)
and its linked operation guides for complete examples and authentication.
Individual documentation pages (https://docs.weaviate.io/weaviate/manage-collections/collection-operations, or https://docs.weaviate.io/weaviate/search/similarity for example) hold code examples.

These native APIs cover more than the running service:

| Operation | Native v6 API | What to preserve |
| --- | --- | --- |
| Search | `collection.Query.Hybrid`, `NearText`, `NearVector`, `BM25`, etc. | Ranking parameters, fusion, filters, pagination, target selection and response mapping |
| Aggregation | `collection.Aggregate.OverAll` and `collection.Aggregate.OverAll.GroupBy` | Grouping, counts, statistics, ordering and result limits |
| Collection management | `client.Collections.Exists`, `.Use`, `.Create` | Existing configuration; an explicit creation policy when absent |
| Import with stable IDs | `collection.Data.Insert` with explicit `data.Object.UUID` | Original UUID derivation and complete property maps |
| Count | `collection.Count(ctx)` | The same collection and tenant scope |

See [hybrid search](https://docs.weaviate.io/weaviate/search/hybrid),
[filters](https://docs.weaviate.io/weaviate/search/filters),
[aggregation](https://docs.weaviate.io/weaviate/search/aggregate) and
[object import](https://docs.weaviate.io/weaviate/manage-objects/import).

Keep the loader's original namespace, hash and source key when deriving UUIDs.
Pass the parsed UUID explicitly as `data.Object.UUID`; omitting it generates a
random UUID. `Data.Insert` uses gRPC `BatchObjects`: inserting an existing UUID
overwrites the object rather than merging a partial property patch, so send the
complete intended property map. Check errors,
including per-object failures reported as `data.InsertError`. This variadic
call is not streaming batching; choose suitable batch sizes for the workload.

## Reuse existing collections; configure new ones explicitly

Use `client.Collections.Exists(ctx, name)` to distinguish the two paths:

- **Exists:** obtain a handle with `client.Collections.Use(name)`. This does not
  create or convert the collection. Native v6 query and import APIs can use an
  existing legacy unnamed-vector collection unchanged.
- **Missing:** create the intended collection through `client.Collections.Create`,
  preserving the deployment's vectorizer, model, properties and index settings.

Specify vector settings explicitly. Server defaults changed in
[Weaviate 1.40](https://docs.weaviate.io/deploy/migration#vector-config-defaults):
a definition with no vector parameters creates a collection without a vector,
and `DEFAULT_VECTORIZER_MODULE` is ignored. Enabling a module is not the same
as selecting it for a collection.

### Worked example: an existing self-hosted model2vec configuration

This example uses the `text2vec-model2vec` as an illustrative example.

Keep your own provider, and adopt the pattern rather than adopting this example verbatim.

For this plain legacy configuration, the v5 declaration is:

```go
package migration

import (
    "context"

    "github.com/weaviate/weaviate/entities/models"
    "github.com/weaviate/weaviate-go-client/v5/weaviate"
)

// CreateLegacyCollection creates the original plain model2vec collection.
func CreateLegacyCollection(ctx context.Context, client *weaviate.Client, name string) error {
    return client.Schema().ClassCreator().WithClass(&models.Class{
        Class: name,
        Vectorizer: "text2vec-model2vec",
        Properties: []*models.Property{
            {Name: "body", DataType: []string{"text"}},
            {Name: "kind", DataType: []string{"text"}},
            {Name: "price", DataType: []string{"number"}},
        },
    }).Do(ctx)
}
```

For a fresh collection with the same properties and inference service/model,
the native v6 replacement is:

```go
package migration

import (
    "context"

    weaviate "github.com/weaviate/weaviate-go-client/v6"
    "github.com/weaviate/weaviate-go-client/v6/collections"
    "github.com/weaviate/weaviate-go-client/v6/modules/model2vec"
)

// CreateNamedCollection creates a fresh single-vector model2vec collection.
func CreateNamedCollection(ctx context.Context, client *weaviate.Client, name string) error {
    _, err := client.Collections.Create(ctx, collections.Collection{
        Name: name,
        Properties: []collections.Property{
            {Name: "body", DataType: collections.DataTypeText},
            {Name: "kind", DataType: collections.DataTypeText},
            {Name: "price", DataType: collections.DataTypeNumber},
        },
        Vectors: map[string]collections.VectorConfig{
            "default": {Vectorizer: model2vec.Text2Vec{}},
        },
    })
    return err
}
```

You can thusly replace the v5 collection creation code with the "legacy" vector,
with the v6 code using a named vector with "default" vector, and the vectorizer
explicitly set.

Both leave eligible text-property selection implicit. Adding an explicit property
list can change vectorization, including empty or missing text behavior.
Carry over custom module/property/index settings rather than assuming this
plain example matches them.

This translation was checked with v5.7.4, v6.0.0 and Weaviate 1.39.5: the shown
configuration preserved observed vectorization and query behavior. The schema
JSON and raw vector-response layout differ; check those separately if your
application consumes them. This is not a guarantee for arbitrary vectorizers,
custom settings or every query.

On that server, the same v6 queries with targets omitted worked for an unchanged
legacy collection and a fresh collection with exactly one named vector. Do not
unconditionally add an explicit `"default"` target to legacy queries. Multiple
named vectors need an appropriate target-selection policy. `"default"` is the
example's chosen name, not a compatibility switch.

## Verify completion

1. Inspect every application package and command, not just the service:
   `go list -f '{{.ImportPath}}: {{join .Imports " "}}' ./...`.
2. Replace remaining application-owned v5 calls/imports and remove the v5 module
   requirement. Stable UUIDs, ingestion, collection management and counts do not
   require retaining v5. Remove obsolete direct server schema-model imports;
   keep v6's legitimate generated-protocol dependency.
3. Run `go mod tidy` after migrating source, then build/test every command.
   Diagnose actual dependency errors rather than assuming they belong to v6.
4. Compare the recorded responses, ranking, filters, pagination and aggregates.
   Import twice and compare object IDs, full properties and counts.
5. Exercise the complete application with GraphQL disabled on a non-production
   instance, following the [official migration check](https://docs.weaviate.io/weaviate/api/graphql/migration#test-with-graphql-disabled).
   Successful queries alone do not establish a complete SDK upgrade.
