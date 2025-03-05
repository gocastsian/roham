#! /bin/bash

docker compose \
--env-file ./deploy/.env \
--project-directory . \
-f ./deploy/rohom/development/traefik-compose.yml \
-f ./deploy/user/development/docker-compose.yaml \
 "$@"
