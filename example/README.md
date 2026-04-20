# example

This directory is a collection of examples of using the Go client in different scenarios.
Each example is structured as a standalone Go project that you can check out and run.

These examples also run in the CI and function as **journey tests** to make sure the client
requests are understood by the Weaviate server. To save valuable CI time, most if not all
examples should target a dedicated WCD cluster instead of spinning up a Docker container.

To run locally, [create a new sandbox cluster](https://docs.weaviate.io/cloud/manage-clusters/create) (or use an existing one) and configure the client via environment variables:

```sh
export WEAVIATE_HOST=<wcd_cluster_url>
export WEAVIATE_API_KEY=<api_key>
go run ./example/basic
```

