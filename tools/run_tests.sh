#!/usr/bin/env bash

function main() {

  GO_MOD_PKG=$(head -n 1 go.mod | sed -r 's/module //g')
  LIB_VER=$(basename "$GO_MOD_PKG")
  PKG_TEST="weaviate-go-client/$LIB_VER/test/"

  # This script runs all non-benchmark tests if no CMD switch is given and the respective tests otherwise.
  run_all_tests=true
  run_unit_tests=false
  run_integration_tests=false
  run_auth_integration_tests=false
  run_deprecated_tests=false


  while [[ "$#" -gt 0 ]]; do
      case $1 in
          --unit-only) run_all_tests=false; run_unit_tests=true;;
          --integration-only) run_all_tests=false; run_integration_tests=true;;
          --auth-integration-only) run_all_tests=false; run_auth_integration_tests=true;;
          --deprecated-only) run_all_tests=false; run_deprecated_tests=true;;
          *) echo "Unknown parameter passed: $1"; exit 1 ;;
      esac
      shift
  done

  # Jump to root directory
  cd "$( dirname "${BASH_SOURCE[0]}" )"/.. || exit

  if  $run_unit_tests || $run_all_tests
  then
    echo_green "Run all unit tests..."
    go test -v ./weaviate/...
    echo_green "Unit tests successful"
  fi

  if $run_integration_tests || $run_all_tests
  then
    echo_green "Run all integration tests..."
    run_integration_tests  "$@"
    echo_green "integration tests successful"
  fi

  if $run_auth_integration_tests || $run_all_tests
  then
    echo_green "Run all auth integration tests..."
    ./test/start_containers.sh
    go test -v -count 1 -race test/auth/*.go
    ./test/stop_containers.sh
    echo_green "auth integration tests successful"
  fi

  if $run_deprecated_tests || $run_all_tests
  then
    echo_green "Run all deprecated tests..."
    go test -v -count 1 -race test_deprecated/*.go
    echo_green "deprecated tests successful"
  fi
  echo "Done!"
}

function run_integration_tests() {
  for pkg in $(go list ./... | grep "$PKG_TEST" | grep -v 'auth' ); do
      if ! go test -count 1 -race "$pkg"; then
        echo "Test for $pkg failed" >&2
        return 1
      fi
    done
}


function echo_green() {
  green='\033[0;32m'
  nc='\033[0m'
  echo -e "${green}${*}${nc}"
}


main "$@"