#!/usr/bin/env bash
set -euo pipefail

proto_dir="$PWD/api/proto"
gen_dir="$PWD/internal/gen"
out_dir="$gen_dir/proto"

echo "Installing latest gRPC libs..."

if command -v brew >/dev/null 2>&1; then
    brew update && brew upgrade protobuf protolint
fi

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

echo "Generating Go protocol stubs..."

rm -rf $out_dir && mkdir -p $out_dir  && cd $gen_dir \
    && protoc \
      --proto_path=$proto_dir \
      --go_out=$out_dir \
      --go_opt=paths=source_relative \
      --go-grpc_out=$out_dir \
      --go-grpc_opt=paths=source_relative \
      $proto_dir/v1/*.proto

echo "Lint generated file headers..."

# cd - && sed -i '' '/versions:/, /source: .*/d' "$out_dir/**/*.go"

echo "Fix imports and format generated files with 'gofumpt'..."

goimports -w $out_dir
gofumpt -w $out_dir

echo "Done"
