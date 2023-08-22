# Weaviate go client  <img alt='Weaviate logo' src='https://raw.githubusercontent.com/weaviate/weaviate/19de0956c69b66c5552447e84d016f4fe29d12c9/docs/assets/weaviate-logo.png' width='180' align='right' />

A go native client for weaviate.

## Usage

In order to get the go client v4 issue this command:

```bash
$ go get github.com/weaviate/weaviate-go-client/v4@v4.x.x
```

where `v4.x.x` is the desired go client v4 version, for example `v4.10.0`

Add dependency to your `go.mod`:

```go
require github.com/weaviate/weaviate-go-client/v4 v4.10.0
```

Connect to Weaviate on `localhost:8080` and fetch meta information

```go
package main

import (
  "context"
  "fmt"

  client "github.com/weaviate/weaviate-go-client/v4/weaviate"
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

- [Documentation](https://weaviate.io/developers/weaviate/current/client-libraries/go.html).

## Support

- [Stackoverflow for questions](https://stackoverflow.com/questions/tagged/weaviate).
- [Github for issues](https://github.com/weaviate/weaviate-go-client/issues).

## Contributing

- [How to Contribute](https://github.com/weaviate/weaviate-go-client/blob/main/CONTRIBUTE.md).

## Build Status

[![Build Status](https://github.com/weaviate/weaviate-go-client/actions/workflows/.github/workflows/tests.yaml/badge.svg?branch=main)](https://github.com/weaviate/weaviate-go-client/actions/workflows/.github/workflows/tests.yaml)
