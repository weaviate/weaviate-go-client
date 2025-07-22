# Weaviate Go client  <img alt='Weaviate logo' src='https://raw.githubusercontent.com/weaviate/weaviate/19de0956c69b66c5552447e84d016f4fe29d12c9/docs/assets/weaviate-logo.png' width='180' align='right' />

A native Go client for [Weaviate](https://github.com/weaviate/weaviate).

## Usage

> :warning: `v5.0.0` cannot be installed due to a malformed `go.mod` file.
> Prefer `v5.0.1` or higher.

In order to get the Go client `v5` issue this command:

```bash
$ go get github.com/weaviate/weaviate-go-client/v5@v5.x.x
```

where `v5.x.x` is the desired Go client `v5` version, for example `v5.4.1`.

Add dependency to your `go.mod`:

```go
require github.com/weaviate/weaviate-go-client/v5 v5.4.1
```

Connect to Weaviate on `localhost:8080` and fetch meta information

```go
package main

import (
  "context"
  "fmt"

  client "github.com/weaviate/weaviate-go-client/v5/weaviate"
)

func main() {
  config := client.Config{
    Scheme: "http",
    Host:   "localhost:8080",
  }
  c, err := client.NewClient(config)
  if err != nil {
    fmt.Printf("Error occurred %v", err)
    return
  }
  metaGetter := c.Misc().MetaGetter()
  meta, err := metaGetter.Do(context.Background())
  if err != nil {
    fmt.Printf("Error occurred %v", err)
    return
  }
  fmt.Printf("Weaviate meta information\n")
  fmt.Printf("hostname: %s version: %s\n", meta.Hostname, meta.Version)
  fmt.Printf("enabled modules: %+v\n", meta.Modules)
}
```

## Documentation

- [Documentation](https://docs.weaviate.io/weaviate/client-libraries/go).

## Support

- [Stackoverflow for questions](https://stackoverflow.com/questions/tagged/weaviate).
- [Github for issues](https://github.com/weaviate/weaviate-go-client/issues).

## Contributing

- [How to Contribute](https://github.com/weaviate/weaviate-go-client/blob/main/CONTRIBUTE.md).

## Build Status

[![Build Status](https://github.com/weaviate/weaviate-go-client/actions/workflows/.github/workflows/tests.yaml/badge.svg?branch=main)](https://github.com/weaviate/weaviate-go-client/actions/workflows/.github/workflows/tests.yaml)
