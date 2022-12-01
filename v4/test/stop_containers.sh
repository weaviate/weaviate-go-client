#!/bin/bash

docker-compose -f v4/test/docker-compose.yaml down --remove-orphans
docker-compose -f v4/test/docker-compose-azure.yaml down --remove-orphans
docker-compose -f v4/test/docker-compose-okta.yaml down --remove-orphans
docker-compose -f v4/test/docker-compose-wcs.yaml down --remove-orphans