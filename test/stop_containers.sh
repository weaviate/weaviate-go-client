#!/usr/bin/env bash

set -eou pipefail

docker-compose -f test/docker-compose.yml down --remove-orphans
docker-compose -f test/docker-compose-azure.yml down --remove-orphans
docker-compose -f test/docker-compose-okta-cc.yml down --remove-orphans
docker-compose -f test/docker-compose-okta-users.yml down --remove-orphans
docker-compose -f test/docker-compose-wcs.yml down --remove-orphans
docker-compose -f test/docker-compose-cluster.yml down --remove-orphans

echo "All containers stopped"
