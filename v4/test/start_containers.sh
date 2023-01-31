#!/bin/bash

docker-compose -f v4/test/docker-compose-azure.yml up -d
docker-compose -f v4/test/docker-compose-okta-cc.yml up -d
docker-compose -f v4/test/docker-compose-okta-users.yml up -d
docker-compose -f v4/test/docker-compose-wcs.yml up -d
docker-compose -f v4/test/docker-compose.yml up -d
