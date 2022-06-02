#!/bin/bash

set -euo pipefail

VERSION=${1-}
REQUIRED_TOOLS="sed git basename pwd"

if test -z "$VERSION"; then
  echo "Missing version parameter. Usage: $0 VERSION"
  exit 1
fi

if case $VERSION in v*) false;; esac; then
  VERSION="v$VERSION"
fi

for tool in $REQUIRED_TOOLS; do
  if ! hash "$tool" 2>/dev/null; then
    echo "This script requires '$tool', but it is not installed."
    exit 1
  fi
done

if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "Cannot prepare release, a release for $VERSION already exists"
  exit 1
fi

DIR=$(pwd -P)
LIB_VER=$(basename $DIR)

sed -i '' "s/^$ go get github.com\/semi-technologies\/weaviate-go-client\/.*/$ go get github.com\/semi-technologies\/weaviate-go-client\/$LIB_VER@$LIB_VER.x.x/" README.md
sed -i '' "s/^where \`.*/where \`$LIB_VER.x.x\` is the desired go client v4 version, for example \`$VERSION\`/" README.md
sed -i '' "s/^require github.com\/semi-technologies\/weaviate-go-client\/.*/require github.com\/semi-technologies\/weaviate-go-client\/$LIB_VER $VERSION/" README.md
sed -i '' "s/^  client \"github.com\/semi-technologies\/weaviate-go-client.*/  client \"github.com\/semi-technologies\/weaviate-go-client\/$LIB_VER\/weaviate\"/" README.md

git commit -a -m "Release $VERSION version"

git tag -a "$VERSION" -m "$VERSION"
