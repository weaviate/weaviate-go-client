# Weaviate Go Client v6 - Design Decisions

**Audience:** Internal stakeholders
**Purpose:** High-level overview of key architectural decisions through examples

## Overview

Goals: a native Go experience with:
- **Explicit, readable APIs**
- **Minimal boilerplate**
- **Type safety**

Here are some concrete examples illustrating the design choices.

---

## Single insert

```go
songs := client.Collections.Use("Songs")
songs.Data.Insert(ctx,
    data.WithProperties(Song{
        Title:  "Bohemian Rhapsody",
        Artist: "Queen",
        Year:   1975,
    }),
    data.WithVector(types.Vector{Single: []float32{0.1, 0.2, 0.3}}),
)
```

- **Functional, variadic options**
- **Flexible properties**: Accept a struct or `map[string]any`
- **Unified Vector type**: Single `Vector` type for both 1D (`Single`) and 2D (`Multi`) vectors

---

## NearVector query

```go
queryVector := types.Vector{Single: []float32{0.1, 0.2, 0.3}}
result, err := songs.Query.NearVector(ctx, queryVector,
    query.WithLimit(10),
    query.WithDistance(0.5),
    query.WithOffset(20),
)

// Access results
for _, obj := range result.Objects {
    fmt.Println(obj.UUID)                     // Direct access
    fmt.Println(obj.Vectors["text"].Single)   // Vectors always accessible
    title := obj.Properties["title"].(string) // Map-based by default
}
```

- **Variadic options**: Similar to insert
- **Map-based objects**: Default object returns are `map[string]any` for convenient prototyping
- **UUID and Vector fields**

---

## Opt-In generics

```go
result, err := songs.Query.NearVector(ctx, queryVector)
// result.Objects -> []WeaviateObject[map[string]any]

type Song struct {
    Title  string `json:"title"`
    Artist string `json:"artist"`
    Year   int    `json:"year"`
}

// Use `Scan` to convert to typed objects
typedObjects := query.Scan[Song](result)

for _, obj := range typedObjects {
    fmt.Println(obj.Properties.Title)  // No type assertion!
    fmt.Println(obj.Properties.Artist) // IDE autocomplete works
    fmt.Println(obj.UUID)
}
```

- **Additional function for type safety**: `Scan[T]()` converts map-based results to typed structs
    - Simplicity with optional type safety

---

## GroupBy queries

```go
// Standard vs grouped NearVector queries - function-as-receiver pattern
single, err := songs.Query.NearVector(ctx, vector, query.WithLimit(10))
groups, err := songs.Query.NearVector.GroupBy(ctx, vector, "category", query.WithLimit(10))
```

- **Function-as-receiver**: More common query prioritized, less common GroupBy as method
- **Shared options with different return types**

---

## Multiple Vector Formats

```go
// Default vector (unnamed)
songs.Data.Insert(ctx,
    data.WithProperties(data),
    data.WithVector(types.Vector{Single: []float32{0.1, 0.2, 0.3}}),
)

// Named single-dimensional vector
songs.Data.Insert(ctx,
    data.WithProperties(data),
    data.WithVector(types.Vector{
        Name:   "text_embedding",
        Single: []float32{0.1, 0.2, 0.3},
    }),
)

// Named multi-dimensional vector
songs.Data.Insert(ctx,
    data.WithProperties(imageData),
    data.WithVector(types.Vector{
        Name:  "colbert",
        Multi: [][]float32{
          {0.1, 0.2}, {0.3, 0.4}
        },
    }),
)

// Multiple vectors at once
songs.Data.Insert(ctx,
    data.WithProperties(data),
    data.WithVector([]types.Vector{
        {Name: "single_vec", Single: []float32{0.1, 0.2, 0.3}},
        {Name: "matrix_vec", Multi: [][]float32{
            {0.4, 0.5}, {0.6, 0.7}
        }},
    }
)
```

---

## Multi-Vector search

```go
singleVec := types.Vector{Name: "single_vec", Single: []float32{0.1, 0.2, 0.3}}
matrixVec := types.Vector{Name: "matrix_vec", Single: []float32{0.4, 0.5, 0.6}}

result, err := songs.Query.NearVector(ctx,
    query.Average(singleVec, matrixVec),
    query.WithLimit(10),
)

result, err := songs.Query.NearVector(ctx,
    query.ManualWeights(
        query.Target(singleVec, 0.7),
        query.Target(matrixVec, 0.3),
    ),
    query.WithLimit(10),
)
```

---

## Reuse Vectors from results

```go
result, err := songs.Query.NearVector(ctx, queryVector, query.WithLimit(10))

// Re-insert using returned vectors
for _, obj := range result.Objects {
    newID, err := songs.Data.Insert(ctx,
        data.WithProperties(obj.Properties),  // Reuse properties map
        data.WithVector(obj.Vectors),         // Reuse entire vector map
    )
}
```

---

### Pointers for optional values

Implication: Use pointers for optional fields in structs to distinguish between zero values and unset fields.

---


## Transport

Transport layer is responsible for executing requests to the Weaviate REST and gRPC endpoints.
The key requirement for the transport layer is _transparency_: the caller (public layer, the `client`) should be
entirely oblivious to the method (REST / gRPC, the _how?_) by which a request is executed.

A good abstraction must be:

- Flexible. Has to accomodate both REST and gRPC requests seamlessly.
- Robust. Even though Transport is not user-facing, we want to reduce the chance of a developer mistake.
- Minimal. A small API is easy to mock in tests, reason about, and maintain.


### Transparency: ignoring the _how?_

The easiest way to have the client ignore "how something should be done" is to let focus it entirely on "what should be done".

"What should be done" is invariably about _data transformation_ (passing inputs and receiving outputs) and we can enumerate all possible inputs and outputs thusly:

```go
package api

type (
    SearchRequest   struct{}
    SearchResponse  struct{}

    InsertObjectRequest     struct{}
    InsertObjectResponse    struct{}

    CreateBackupRequest  struct{}
    CreateBackupResponse struct{}

    // and so on...
)
```

> [!IMPORTANT]
> `api` is an **internal** package under "internal/api"

To execute a request the client has to fill out the request parameters, pass the request to the transport and then receive the output.
Something like this:

```go
package query

func (c *Client) NearVector() {
    var req *api.SearchRequest = &api.SearchRequest{...}
    var resp *api.SearchResponse = execute(req)
}
```

For all we care, `api.SearchRequest` can be delegated to a Postgres connector using `pg_vector` under the hood. Transparency!

Here's another example, this time with a "REST"-request:

```go
package backup

func (c *Client) CreateBackupRequest() {
    var req *api.CreateBackupRequest = &api.CreateBackupRequest{...}
    var resp *api.CreateBackupResponse = execute(req)
}
```

In both of the examples above `execute()` is some fixture that implements Transport interface.
Also worth noting is the fact that neither of the -Request structs has fields like `Endpoint` or `Uses127Api`
or anything related to the underlying transport for that matter -- only data pertinent to the request: `CollectionName`, `Tenant`, `Bucket`, etc.

A set of structs with parameters for each of the supported requests therefore is the outward-facing half of the transport layer.


### Implementation: acknowledging the _how?_

Let us ignore the actual shape of the Transport interface a bit longer and see how a hypothetical implementation
might go about fulfilling the requests above.

First, assume the said implementation holds both HTTP and gRPC transports as unexported dependencies.

```go
type transport struct { gRPC, http any }
```

Now it's task boils down to passing a request to the right handler based on... what?
To answer that question let's see what information is needed to "describe" each type of the request.

Rougly speaking, a **REST request** is described by its `method`, URL (`path` and `query` components), and `body`.
The interface below captures that:

```go
package api

type RESTRequest interface{
    Method() string
    Path() string
    Query() url.Values
    Body() any // does not need a specific shape, anything can be marshaled to JSON.
}
```

Each **gRPC request** has a corresponding `message`, for which we generate a Go stub.
The type of the message struct is therefore it's most concise discriminator. Put diffeerently,
the _body_ of a gRPC request uniquely describes it. We can encode this notion in an interface:

```go
package api

type GRPCRequest interface {
    Body() any // generated Go stubs do not implement any shared interface, no need to impose one here either.
}
```

Extracting the common parts (`Body() any`) and polishing up the names here's what we end up with:

```go
package api

// Request is anything that can have a body.
type Request interface {
    Body() any
}

// Endpoint describes a REST request.
type Endpoint interface {
    Request

    Method() string
    Path() string
    Query() url.Values

}
```

The `api.Request` interface is a general form of the request and `api.Endpoint` is a narrower kind, tailored to a REST API.
To correctly hand off the request, our transport implementation can do a simple switch:

```go
switch req := req.(type) {
case api.Endpoint:
    http.execute(req)
case *api.SearchRequest:
    gRPC.Search(req) // generated proto.WeaviateClient has a dedicated Search() method
case *api.AggregateRequest:
    gRPC.Aggregate(req) // generated proto.WeaviateClient has a dedicated Aggregate() method
default:
    panic("unknown request type")
}
```

If `panic` made you frown, see the **Discussion** below.

Before we wrap up let's quickly look at response handling.


### Love, `dest`, and Robots

Previously we mentioned that the public layer of our client is only interested in data -- input and output, not behavior.
In _structs_, not _interfaces_. The fact that something like `api.CreateBackupRequest` implements `api.Endpoint` is only relevant in the transport layer itself.
Which is why we shouldn't try and come up with an `api.Response` interface, there's _simply no use for it_. We want data!

In our design the caller owns the response struct, and the transport layer receives it as an opaque pointer `dest any`, much like in `json.Unmarshal`.
Reading a REST response into such a pointer is trivial using the same `json.Unmarshal` API. In case of gRPC, the transport needs to cast it back to
the appropriate response type before writing to it.


```go
func execute(req api.Request, dest any) {
    // REST
    json.Unmarshal(responseBody, dest)

    // gRPC, example for api.SearchRequest
    if resp, ok := dest.(*api.SearchResponse); ok {
        *resp = gRPC.search(req)
    }
}
```

In that type-cast above lies a tiny opportunity for a bug: what if someone passes `*api.AggregateRequest` with an `*api.SearchRequest`?
Luckily, this can be caught early with a simple test:

```go
func mockTransport(t *testing.T) Transport {
    return func(req api.Request, dest any) {
    switch (req).(type) {
        case *api.SearchRequest:
            require.IsType(t, (*api.SearchResponse)(nil), dest)
        case *api.AggregateRequst:
            require.IsType(t, (*api.AggregateResponse)(nil), dest)
        case *api.CreateBackupRequest:
        // and so on
        }
    }
}

// Now inject this interface into each client and call its public methods.
// The `dest` type is independent from the input parameters, so calling each method once
// is enough to assert that all of them use correct response types.
// A small price to pay for not having a redundant Response interface. What methods was it gonna have anyways?
transport := mockTransport(t)
query.NewClient(transport).NearVector()
backup.NewClient(transport).Create()
```

The above is only a rough sketch, but with a _minimal_ interface, mocking transport out is extremely cheap and can be done in a handful of lines.
Finally, if response is a pointer we get a neat property of passing a `nil` to tell transport we aren't interested in the response body and that it needn't bother unmarshaling it.

Alright! This being said, we now have everything we need to formulate the Transport interface.


### The Transport Interface

To summarize the points above:

- Transport interface makes execution transparent to the caller; the latter should focus on the "what" not the "how".
- `internal/api` package exports structs for all supported requests without mandating their execution.
- `api.Request` defines the common shape of a request object.
- Transport returns responses to the caller through an opaque pointer `dest any`.
- The API surface of the interface should be kept to a minimum.

At this point we can practically touch the interface. And it `Go`es like this:

```go
package internal

type Transport interface {
	Do(ctx context.Context, req api.Request, dest any) error
}
```

If this is underwhelming, then we've done well.
At a risk of repeating ourselves, here're the supporting docs for the `Do()` method:

```go
	// Do executes a request and populates the response object.
	// Response dest SHOULD be nil if no response is expected
	// and MUST be a non-nil pointer otherwise.
	//
	// To keep execution transparent to the caller, the request type
	// only enforces a minimal constraint -- a request is anything
    // that MAY have a body.
	//
	// The "internal/api" package defines structs for all
	// supported requests, which in turn implement api.Request.
    // The contract is that Transport is able to execute any
    // one of those requests.
    //
    // The transport is also able to execute any custom [api.Endpoint].
```

### Discussion

1. Why not introduce generic parameters for `req` and `dest` to avoid type-casting altogether?

Firstly, because that generic parameter would need to be defined at the `interface` level:

```go
interface Transport[Req any, Resp any] interface{
    Do(context.Context, Req, Resp)
}
```
and provided at transport instantiation!, which is only done once and not per-request.

A `Transport[api.SearchRequest, api.SearchResponse]` cannot be used to create a backup or run an aggregation query.

Secondly, because mismatching the request and response types is a _developer_ (our!) error.
And for every scenario where someone accidentaly passes a wrong response type as `dest`,
there exists an equally-likely scenario of someone passing a wrong generic parameter:

```go
c.transport.Do[api.SearchRequest, api.AggregateResponse](ctx, req, dest)
```
Not to mention that a developer from either scenario would have a really hard time converting `api.AggregateRespose` to `query.Response` (remember, types from the `api` package are never relayed to the user directly and are always re-packaged into a another user-facing struct).

To the best of my knowledge Go doesn't have a notion of conditional types (`T = X extends Y ? Z : never`).
Even if there was a way to do it by arranging several interfaces in some way, I'd prefer to write that simple test at this point (also see pt. 2).

2. `panic("unknown request type")` -- what's that all about?

[Assertions detect programmer errors](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety). Since transport layer is internal, the only way an error can be introduced is through a developer (our, not user) mistake. Should a mistake like this happen, a simple test we've mentioned before will catch that long before this code hits anyone's production server. This might be a hard sell, so I'm happy to return an error there as well.

> [!NOTE]
> You can find a reference implementation of [`internal.Transport`](https://github.com/weaviate/weaviate-go-client/blob/dyma/v6/internal/transport.go) along with some more APIs and requests on [`dyma/v6`](https://github.com/weaviate/weaviate-go-client/tree/dyma/v6).



