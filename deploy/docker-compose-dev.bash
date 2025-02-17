#! /bin/bash

docker compose \
--env-file ./deploy/.env \
--project-directory . \
 "$@"

#TODO complete this