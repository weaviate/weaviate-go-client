#!/usr/bin/env sh

# I sincerely wish there were no shell scripts in this repo
# but this is easier to read than a one-line go:generate.

set -e

OUT=$(pwd)
MOD=github.com/weaviate/weaviate-go-client/v6/internal/api/internal/gen/proto

rm -rf v1
mkdir -p v1

(
    cd $(git rev-parse --show-toplevel)

    bin/protoc \
      --go_out=module=$MOD:$OUT \
      --go-grpc_out=module=$MOD:$OUT \
      api/proto/v1/*.proto
)
