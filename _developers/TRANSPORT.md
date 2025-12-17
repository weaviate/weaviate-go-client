# Transport

Transport layer is responsible for executing requests to the Weaviate REST and gRPC endpoints.
The key requirement for the transport layer is _transparency_: the caller (public layer, the `client`) should be
entirely oblivious to the method (REST / gRPC) by which a request is executed, the _how?_.

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
package _

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
package _

type GRPCRequest interface {
    Body() any // generated Go stubs do not implement any shared interface, no need to impose one here either.
}
```

Extracting the common parts (`Body() any`) and polishing up the names here's what we end up with:

```go
package _

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

The `Request` interface is a general form of the request and `Endpoint` is a narrower kind, tailored to a REST API.
To correctly hand off the request, our transport implementation can do a simple switch:

```go
switch req := req.(type) {
case Endpoint:
    http.execute(req)
case *api.SearchRequest:
    gRPC.Search(req) // generated proto.WeaviateClient has a dedicated Search() method
case *api.AggregateRequest:
    gRPC.Aggregate(req) // generated proto.WeaviateClient has a dedicated Aggregate() method
default:
    panic("unknown request type")
}
```

The actual implementation is even _shorter_; this example shows it in a bit more detail.

Before we wrap up let's quickly look at response handling.


### Love, `dest`, and Robots

Previously we mentioned that the public layer of our client is only interested in data -- input and output, not behavior.
In _structs_, not _interfaces_. The fact that something like `api.CreateBackupRequest` implements `Endpoint` is only relevant in the transport layer itself.
Which is why we shouldn't try and come up with an `Response` interface, there's _simply no use for it_. We want data!

In our design the caller owns the response struct, and the transport layer receives it as an opaque pointer `dest any`, much like in `json.Unmarshal`.
Reading a REST response into such a pointer is trivial using the same `json.Unmarshal` API. In case of gRPC, the transport needs to cast it back to
the appropriate response type before writing to it.


```go
func execute(req Request, dest any) {
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
    return func(req Request, dest any) {
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
- `Request` defines the common shape of a request object.
- Transport returns responses to the caller through an opaque pointer `dest any`.
- The API surface of the interface should be kept to a minimum.

At this point we can practically touch the interface. And it `Go`es like this:

```go
package internal

type Transport interface {
	Do(ctx context.Context, req Request, dest any) error
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
```

### Why is this better than the `v5` transport?

As a refresher, the transport layer in `v5` essentially looks like this:

```go
type Transport interface {
    RunREST(ctx context.Context, hostAndPath string, restMethod string, requestBody interface{}) (*ResponseData, error)
    RunRESTExternal(ctx context.Context, hostAndPath string, restMethod string, requestBody interface{}) (*ResponseData, error)
    Search(ctx context.Context, req *proto.SearchRequest) (*pb.SearchReply, error)

    // Not supported in v5:
    BatchObjects(ctx context.Context, req *proto.BatchObjectsRequest) (*pb.BatchObjectsReply, error)
    BatchReferences(ctx context.Context, req *proto.BatchReferencesRequest) (*pb.BatchReferencesReply, error)
    BatchDelete(ctx context.Context, req *proto.BatchDeleteRequest) (*pb.BatchDeleteReply, error)
    TenantsGet(ctx context.Context, req *proto.TenantsGetRequest) (*pb.TenantsGetReply, error)
    Aggregate(ctx context.Context, req *proto.AggregateRequest) (*pb.AggregateReply, error)
}
```

I say "essentially" because these methods do not exist in a single interface and part of them (gRPC requests) does not exist at all.
Looking at the `Search` method, I think I'd be safe to assume `v5` would have them if it were to support other gRPC calls.

The key deficiency of this design is that it results in _so much more code_ than necessary.

First, consider that `v5.Transport` has 8 methods, where `v6.Transport` only has 1.
In practice the `v6.Transport` implementation must actually have 3 methods: a public `Do()` and private `doREST()` and `doGRPC()`.
Even so, it has 5 methods less. Most importantly though -- these 3 methods can service _any number of requests_. On the other hand,
if tomorrow `backup` API was migrated to gRPC, the `v5.Transport` surface would grow to **14 methods** to accomodate the 6 requests in `backup`.

Next, consider what using the `v5.Transport` looks like ([data.Creator.Do](https://github.com/weaviate/weaviate-go-client/blob/cf0624dc4258b33ac8b5c1139a64aec6b87e005b/weaviate/data/creator.go#L76)):

```go
// EDIT: PayloadObject() and buildPath() are inlined for the purposes of this example.
// Neither is re-used in any other place so I think it's fair to include them in full.
func (creator *Creator) Do(ctx context.Context) (*ObjectWrapper, error) {
    var err error
	var responseData *connection.ResponseData

	// object, _ := creator.PayloadObject()
    object := models.Object{
		Class:      creator.className,
		Properties: creator.propertySchema,
		Vector:     creator.vector,
		Vectors:    creator.vectors,
		Tenant:     creator.tenant,
	}
	if creator.uuid != "" {
		object.ID = strfmt.UUID(creator.uuid)
	}

	// path := creator.buildPath()
    path := "/objects"
	pathParams := url.Values{}

	if creator.consistencyLevel != "" {
		pathParams.Set("consistency_level", creator.consistencyLevel)
	}

	if len(pathParams) > 0 {
		path = fmt.Sprintf("%s?%v", path, pathParams.Encode())
	}

	responseData, err = creator.connection.RunREST(ctx, path, http.MethodPost, object)
	respErr := except.CheckResponseDataErrorAndStatusCode(responseData, err, 200)
	if respErr != nil {
		return nil, respErr
	}

	var resultObject models.Object
	parseErr := responseData.DecodeBodyIntoTarget(&resultObject)
	return &ObjectWrapper{
		Object: &resultObject,
	}, parseErr
}
```

Remember, this is the API layer. In the API layer, why do we know that `consistecy_level` belongs into a URL query and not in the path parameters?
Why do we need to know that `200` is the "OK" code  and not `201`? Why does the transport return what is essentially _raw bytes_ in a wrapper
that we need to decode in a separate step? Why should we care if it's a REST or a gRPC request? Also, should we call `RunREST` or `RunRESTExternal` and what is the difference between the two?

That's a lot of questions to answer in the API layer. And even after answering them and doing all this work we still return a `models.Object` to the user, albeit wrapped in an `ObjectWrapper { Object *models.Object }`! If we wanted to return a nice public-facing struct, we'd be _easily_ looking at another 20-30 lines of code, because `models.Object` has such niceties as `AdditionalProperties map[string]any` for metadata and `Vectors map[string]any` for named vectors.

Not to mention that the function above is 38 lines of code (+20-30 if we don't return `models.Object`). That's a lot of lines, considering we might need to recreate them across 20-25(?) requests. We need to write that code, test that code, and maintain that code.
And the more code we write and the more questions we have to answer while doing that, the greater is the opportunity for making a mistake.

To summarize, our `v5` transport is a thin wrapper that isn't doing any of the heavy lifting it should be doing. Not because it the wrapper is thin, but because the abstraction is poorly designed. The proposed Transport interface **drastically** reduces the amount of code we need to write now and in the future, when the API evolves, by providing the right abstractions.

To exemplify, here's what the `Insert()` method looks like in `v6` with the proposed transport design:

```go
func (c *Client) Insert(ctx context.Context, options ...InsertOption) (*types.Object[types.Properties], error) {
	var ir insertRequest
	InsertOptions(options).Apply(&ir) // Build the request

	req := &api.InsertObjectRequest{
        RequestDefaults:    c.defaults, // tenant + consistency level, set once per collection handle
		UUID:               ir.UUID,
		Properties:         ir.Properties,
		Vectors:            ir.Vectors,
	}

	var resp api.InsertObjectResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("insert object: %w", err)
	}

	return &types.Object{
        UUID:               resp.UUID,
        Properties:         resp.Properties,
        Vectors:            types.Vectors(resp.Vectors),
        CreationTimeUnix:   resp.CreationTimeUnix,
        LastUpdateTimeUnix: resp.LastUpdateTimeUnix,
    }, nil
}
```

That's 24 lines, plus we return a user-facing `types.Object`. In the reference implementation The [code executing this request](https://github.com/weaviate/weaviate-go-client/blob/cf0624dc4258b33ac8b5c1139a64aec6b87e005b/weaviate/data/creator.go#L16) is ~60 lines long all told. Implementing `Insert` in `v5` takes roughly the same amount of code as implementing `Insert` _and_ the underlying transport which will be reused for all other requests.


## Discussion

1. We could re-use the `v5` transport code as is, copy-pasting most of it.

That doesn't change the _test_ and _maintain_ parts -- we still need to do that.
Also at this point we should treat `v5` transport as an external dependency that we're adding to the project and ask ourselves:

- Does it solve our problem? How much time does it save us if any?
- Are we willing to maintain that amount of code?
- Can we do better? In less time?

I firmly believe that the answer to the last 2 question is yes.

2. Why not introduce generic parameters for `req` and `dest` to avoid type-casting altogether?

Because that generic parameter would need to be defined at the `interface` level:

```go
interface Transport[Req any, Resp any] interface{
    Do(context.Context, Req, Resp)
}
```
and provided at transport instantiation!, which is only done once and not per-request.

A `Transport[api.SearchRequest, api.SearchResponse]` cannot be used to create a backup or run an aggregation query.

This is crucial point to understand. By introducing generic `[Req, Resp]` arguments into our interface we are forced
to write a _**separate** implementation for **every** single request-response pair_. Instead of 1 well-written and
well-tested `Do()` method we would have to maintain 30-40(?) separate implementations. And that number will keep growing
with each new request we add to the client.

3. `panic("unknown request type")` -- what's that all about?

[Assertions detect programmer errors](https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/TIGER_STYLE.md#safety). Since transport layer is internal, the only way an error can be introduced is through a developer (our, not user) mistake. Should a mistake like this happen, a simple test we've mentioned before will catch that long before this code hits anyone's production server. This might be a hard sell, so I'm happy to return an error there as well.

To make sure panics never manifest in user's code we wrap `panic()` into an `Assert(bool)` that can be disabled via a build flag.

> [!NOTE]
> You can find a reference implementation of [`internal.Transport`](https://github.com/weaviate/weaviate-go-client/blob/dyma/v6/internal/transport.go) along with some more APIs and requests on [`dyma/v6`](https://github.com/weaviate/weaviate-go-client/tree/dyma/v6).

