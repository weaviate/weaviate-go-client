#!/bin/bash

for pkg in $(go list ./... | grep 'weaviate-go-client/v3/test'); do 
  if ! go test -v -count 1 -race "$pkg"; then 
    echo "Test for $pkg failed" >&2; false; exit
  fi
done
