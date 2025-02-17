#!/usr/bin/env bash

set -eou pipefail

export WEAVIATE_VERSION=$1

function wait(){
  MAX_WAIT_SECONDS=60
  ALREADY_WAITING=0

  echo "Waiting for $1"
  while true; do
    if curl -s $1 > /dev/null; then
      break
    else
      if [ $? -eq 7 ]; then
        echo "$1 is not up yet. (waited for ${ALREADY_WAITING}s)"
        if [ $ALREADY_WAITING -gt $MAX_WAIT_SECONDS ]; then
          echo "Weaviate did not start up in $MAX_WAIT_SECONDS."
          exit 1
        else
          sleep 2
          let ALREADY_WAITING=$ALREADY_WAITING+2
        fi
      fi
    fi
  done

  echo "Weaviate is up and running!"
}

docker compose -f test/docker-compose-azure.yml up -d
docker compose -f test/docker-compose-okta-cc.yml up -d
docker compose -f test/docker-compose-okta-users.yml up -d
docker compose -f test/docker-compose-wcs.yml up -d
docker compose -f test/docker-compose-cluster.yml up -d
docker compose -f test/docker-compose.yml up -d
docker compose -f test/docker-compose-rbac.yml up -d

wait "http://localhost:8080"
wait "http://localhost:8081"
wait "http://localhost:8082"
wait "http://localhost:8083"
wait "http://localhost:8085"
wait "http://localhost:8087"
wait "http://localhost:8088"
wait "http://localhost:8089"

echo "All containers running"
