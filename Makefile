start-vectorlayer-app-dev:
	./deploy/docker-compose-dev.bash --profile vectorlayer up

start-user-app-dev:
	./deploy/docker-compose-dev.bash --profile user up

start-all-dev: ## Start all infra services and services
	./deploy/docker-compose-dev.bash --profile vectorlayer  --profile user  up -d

stop-all-dev: ## Stomp and remove all services
	./deploy/docker-compose-dev.bash --profile vectorlayer  --profile user  down

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[.a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
