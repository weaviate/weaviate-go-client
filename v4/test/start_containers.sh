#!/bin/bash

docker-compose -f v4/test/docker-compose-azure.yaml up -d
docker-compose -f v4/test/docker-compose-okta.yaml up -d
docker-compose -f v4/test/docker-compose-wcs.yaml up -d
docker-compose -f v4/test/docker-compose.yaml up -d
