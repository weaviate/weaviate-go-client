# weaviate-go-client  <img width="215" height="250" alt="Weaviate Gopher" align='right' src="https://github.com/user-attachments/assets/807e11d7-7194-4652-be9f-6ce8eef8f80f" />

An idiomatic Go client for [Weaviate](https://github.com/weaviate/weaviate).

## Installation and usage

```sh
go get github.com/weaviate/weaviate-go-client/v6@latest
```

Import the package into your project and create a client instance:

```go
package main

import "github.com/weaviate/weaviate-go-client/v6"

func main() {
    c, err := weaviate.NewLocal()
    if err != nil {
        panic(err)
    }
    defer c.Close()
}
```

See [examples](./example) for using Weaviate's core features and APIs. Refer to our [documentation](https://docs.weaviate.io/weaviate/client-libraries/go) to learn more about Weaviate and how to use it via the Go client.

## Contributing

See our [Contribute Guide](https://github.com/weaviate/weaviate-go-client/blob/main/CONTRIBUTE.md).
