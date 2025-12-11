#!/bin/bash

# Check if dependencies are installed
check_dependencies() {
    local commands=("$@")
    for cmd in "${commands[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            echo "Error: $cmd is a required dependency and not found."
            exit 1
        fi
    done
}

check_dependencies curl tar realpath dirname grep

branchOrTag="${1:-main}"
dir="$(realpath "$( dirname "${BASH_SOURCE[0]}" )" )"
baseout="${dir}/../contracts"
out="${baseout}/proto"

rm -rf "${out}"
mkdir -p "${out}"

curl -LkSs "https://api.github.com/repos/weaviate/weaviate/tarball/${branchOrTag}" -o "${dir}/weaviate.tar.gz"
files="$(tar -tf "${dir}/weaviate.tar.gz" | grep -E 'grpc/proto/v1/[^\.]+\.proto$' | tr '\n' ' ')"
fileSchema="$(tar -tf "${dir}/weaviate.tar.gz" | grep -E 'openapi-specs/schema\.json$' | tr '\n' ' ')"

echo ${fileSchema}

# shellcheck disable=SC2086 # we want to pass multiple arguments to tar
tar --strip-components=3 -C "${out}" -xvf "${dir}/weaviate.tar.gz" ${files}
tar --strip-components=2 -C "${baseout}" -xvf "${dir}/weaviate.tar.gz" ${fileSchema}

rm "${dir}/weaviate.tar.gz"

echo "done"