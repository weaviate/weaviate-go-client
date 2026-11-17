# Contribution Guide

Thanks for looking into contributing to Weaviate Go client!
Please read the entire guide before starting to work on your contribution.


## Development

Start by cloning the repository:

```sh
git clone https://github.com/weaviate/weaviate-go-client.git .
```

Before you continue, make sure to have Go installed locally.
Once that's done, use our build tool to install the remaining dependecies.

```sh
go run ./cmd/build onboard
```

That's it! You should be able to use `go test` and `go generate` commands to run tests and generate REST/gRPC stubs.
Run the test suite to verify that your environment is complete.

### Unit tests

Weaviate Go Client is not only used by the OSS community. It is at the core of our cloud and infrastructure code.
It needs to be reliable, which is why we write a lot of unit tests. Run the command below with `-v` flag to see them.

```sh
go test ./...
```

You'll notice that the test run completes in only a few seconds. This is possible because each test only engages a small
part of the system at a time, and avoids doing expensive operations or sleeping whenever possible.

### Generate REST models and gRPC stubs

These come bundled with the rest of the source, so normally you wouldn't need to generate these.
However, if your local files went out of sync, they can be re-generated using this command:

```go
go generate ./...
```

### Updating API contracts

> [!INFO] This section is only relevant to developers who are part of the [Weaviate org](https://github.com/weaviate).

The protobuf files and the OpenAPI schema in `./api/proto` and `./api/rest` are the source of truth for code generation.
To keep them in sync with their upstreams in [`weaviate/weaviate`](https://github.com/weaviate/weaviate), we validate those in CI (see `./github/workflows/contracts.yaml`).
When (it's a _when_ not an _if_) these contracts get update, you'll be notified by the pipeline, in which case you'll want to fetch the latest versions from upstream
and open a PR with the updated versions.

```sh
go run ./cmd/build contracts
```

Run `go run ./cmd/build --help` to see other tools available to you.

## Submitting the pull request

### Contributor License Agreement

Contributions to Weaviate Go client must be accompanied by a Contributor License Agreement. You (or your employer) retain the copyright to your contribution; this simply gives us permission to use and redistribute your contributions as part of Weaviate Go client. Go to [this page](https://www.semi.technology/playbooks/misc/contributor-license-agreement.html) to read the current agreement.

The process works as follows:

- You contribute by opening a [pull request](#pull-request).
- If your account has no CLA, a DocuSign link will be added as a comment to the pull request.

### How we use Gitflow

The `main` branch is what is released and developed currently.
It is a protected branch, so we use [Gitflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow) to coordinate our development.
Follow these steps when contributing to the project:

- Create a feature branch following this format: `feature/YOUR-FEATURE-NAME`. Your feature must start from the tip of the `main` branch.
- When adding changes, always the the relevant issue in your git commit, a.k.a smart commits: `gh-100: Knit more sweaters.`
- Create a pull request into the `main` branch once your work is ready for review.

a pull request without smart commits, the pull request will be [squashed into](https://blog.github.com/2016-04-01-squash-your-commits/) one git commit.

### Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms.
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v2.0%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

