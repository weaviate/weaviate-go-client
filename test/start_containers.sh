#!/bin/bash

docker-compose -f test/docker-compose-azure.yml up -d
docker-compose -f test/docker-compose-okta-cc.yml up -d
docker-compose -f test/docker-compose-okta-users.yml up -d
docker-compose -f test/docker-compose-wcs.yml up -d
docker-compose -f test/docker-compose.yml up -d
