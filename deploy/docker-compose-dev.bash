#! /bin/bash

docker compose \
--env-file ./deploy/.env \
--project-directory . \
-f ./deploy/roham/development/traefik-compose.yaml \
-f ./deploy/user/development/docker-compose.yaml \
 "$@"
