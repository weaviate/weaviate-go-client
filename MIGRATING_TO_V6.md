# Migrate your application from Go client v5 to v6

Weaviate's [GraphQL API is deprecated](https://weaviate.io/blog/graphql-api-deprecation).
We strongly recommend upgrading every application entry point to
`github.com/weaviate/weaviate-go-client/v6`: services, import commands, scheduled
jobs and maintenance tools, not just query paths.

See the official [GraphQL migration guide](https://docs.weaviate.io/weaviate/api/graphql/migration)
for the retirement schedule, query mappings and
[features without an equivalent](https://docs.weaviate.io/weaviate/api/graphql/migration#graphql-features-without-an-equivalent).

## Before changing code

- Keep the existing embedding provider, model, source-property selection,
  collection configuration and ID derivation. An SDK upgrade is not a reason to
  switch vectorizers or recreate data.
- Make both HTTP and gRPC endpoints reachable. Client construction and
  `client.IsReady` check REST readiness only; exercise a real data/query call
  to verify gRPC. See the [Go connection and compatibility requirements](https://docs.weaviate.io/weaviate/client-libraries/go).

The v6 import is at the module root, without a `/weaviate` suffix:
`weaviate "github.com/weaviate/weaviate-go-client/v6"`.

## Native v6 operations

Use the [official Go documentation](https://docs.weaviate.io/weaviate/client-libraries/go)
and its linked operation guides for complete examples and authentication.

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

Keep the loader's UUID namespace, hash and source key. Omitting `data.Object.UUID`
generates a random UUID. `Data.Insert` overwrites an existing object rather than
merging a partial property patch, so send complete property maps. Handle per-object
failures reported as `data.InsertError`; this is a non-streaming batch call.

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

### Example: preserve a self-hosted model2vec configuration

Keep your application's provider; this plain [model2vec example](https://docs.weaviate.io/weaviate/model-providers/model2vec/embeddings#requirements)
illustrates the configuration mapping for a fresh collection with the same
properties and inference service/model. Replace the v5 `models.Class` field:

```go
Vectorizer: "text2vec-model2vec",
```

with an explicit named vector on v6's `collections.Collection`:

```go
Vectors: map[string]collections.VectorConfig{
    "default": {Vectorizer: model2vec.Text2Vec{}},
},
```

Import `model2vec` from `github.com/weaviate/weaviate-go-client/v6/modules/model2vec`.

This leaves eligible text-property selection implicit. Adding an explicit property
list can change vectorization, including empty or missing text behavior.
Carry over custom module/property/index settings rather than assuming this
plain example matches them. Named vectors change the schema JSON and raw
vector-response layout.

Do not unconditionally add an explicit `"default"` target to legacy queries.
`"default"` is the example's chosen name, not a compatibility switch. Select query
targets to match the collection's actual vector configuration, especially with
[multiple named vectors](https://docs.weaviate.io/weaviate/search/hybrid#named-vectors).

## Test with GraphQL disabled

Use the [official GraphQL-disabled check](https://docs.weaviate.io/weaviate/api/graphql/migration#test-with-graphql-disabled)
on a non-production instance. V5 schema and batch operations can still work
without GraphQL, so passing this check does not prove that the application has
migrated completely to v6.
