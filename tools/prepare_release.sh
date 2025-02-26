#!/bin/bash

set -euo pipefail

REQUIRED_TOOLS="sed git basename pwd"
for tool in $REQUIRED_TOOLS; do
  if ! hash "$tool" 2>/dev/null; then
    echo "This script requires '$tool', but it is not installed."
    exit 1
  fi
done

TARGET_VERSION=${1-}

if test -z "$TARGET_VERSION"; then
  echo "Missing version parameter. Usage: $0 VERSION"
  exit 1
fi

if case $TARGET_VERSION in v*) false;; esac; then
  TARGET_VERSION="v$TARGET_VERSION"
fi

if git rev-parse "$TARGET_VERSION" >/dev/null 2>&1; then
  echo "Cannot prepare release, a release for $TARGET_VERSION already exists"
  exit 1
fi

PACKAGE='github.com/weaviate/weaviate-go-client'
GO_MOD_PKG=$(head -n 1 go.mod | sed -r 's/module //g')

echo "Prepare release for $PACKAGE $TARGET_VERSION"

# Check if we are doing a major version release.
CURRENT_MAJOR_VERSION=$(echo $GO_MOD_PKG | sed -nE "s|$PACKAGE(/v([0-9]+))?|\2|p")
TARGET_MAJOR_VERSION=$(echo $TARGET_VERSION | cut -d'.' -f1 | tr -d 'v')

if [[ $TARGET_MAJOR_VERSION -lt $CURRENT_MAJOR_VERSION ]]; then
    echo "Cannot downgrade major version v$CURRENT_MAJOR_VERSION -> v$TARGET_MAJOR_VERSION"
    exit 1
elif [[ $TARGET_MAJOR_VERSION -gt $CURRENT_MAJOR_VERSION ]]; then
    echo "Major version changed (v$CURRENT_MAJOR_VERSION -> v$TARGET_MAJOR_VERSION), " \
        "updating module declaration and import paths"

    find . -type f \( -name "*.go" -o -name "go.mod" \) \
        -exec sed -i '' "s|$GO_MOD_PKG|$PACKAGE/v$TARGET_MAJOR_VERSION|g" {} \;

    echo "Building project... An error might indicate a malformed go.mod or unresolvable dependencies."
    go build ./...
    echo "OK"
fi

# `v4` -> `v4` (except quoted warnings starting with >)
# `v4.1.1` -> `v5.0.1` (except quoted warnings starting with >)
# v4.x.x -> v5.x.x
# github.com/weaviate/weaviate-go-client/v4 -> github.com/weaviate/weaviate-go-client/v5
# require github.com/weaviate/weaviate-go-client/v4 v4.1.1 -> github.com/weaviate/weaviate-go-client/v5 v5.0.1
sed -i '' \
    -e "/^>/!s/\`v${CURRENT_MAJOR_VERSION}\`/\`v${TARGET_MAJOR_VERSION}\`/g" \
    -e "/^>/!s/\`v[0-9]*.[0-9]*.[0-9]*\`/\`$TARGET_VERSION\`/g" \
    -e "s/v${CURRENT_MAJOR_VERSION}.x.x/v${TARGET_MAJOR_VERSION}.x.x/g" \
    -e "s|\($PACKAGE\)/v${CURRENT_MAJOR_VERSION}|\1/v${TARGET_MAJOR_VERSION}|g" \
    -e "/require/s/ v[0-9]*.[0-9]*.[0-9]*$/ $TARGET_VERSION/g" README.md

git commit -a -m "Release $TARGET_VERSION version"

git tag -a "$TARGET_VERSION" -m "$TARGET_VERSION"
