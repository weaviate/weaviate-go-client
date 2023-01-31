#!/bin/bash

docker-compose -f v4/test/docker-compose.yml down --remove-orphans
docker-compose -f v4/test/docker-compose-azure.yml down --remove-orphans
docker-compose -f v4/test/docker-compose-okta-cc.yml down --remove-orphans
docker-compose -f v4/test/docker-compose-okta-users.yml down --remove-orphans
docker-compose -f v4/test/docker-compose-wcs.yml down --remove-orphans