#!/usr/bin/env bash

set -eou pipefail

export WEAVIATE_VERSION=$1

# Removes containers in the default docker-compose.yml
# and _any other containers_ not defined in that file.
docker compose down --remove-orphans

echo "All containers stopped"
