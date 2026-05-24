start-all:
	./deploy/docker-compose-dev.bash --profile traefik --profile user --profile filer --profile vectorlayer up --build
start-filer-dev:
	./deploy/docker-compose-dev.bash --profile filer up

# Testing commands
test:
	@echo "Running all tests..."
	go test -v -count=1 ./...

test-short:
	@echo "Running short tests..."
	go test -short -v -count=1 ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -count=1 -coverprofile=coverage.out ./...
	@echo "Coverage report generated: coverage.out"
	@echo "View with: go tool cover -html=coverage.out"

test-coverage-html:
	@echo "Generating HTML coverage report..."
	go test -v -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -race -v -count=1 ./...

test-unit:
	@echo "Running unit tests..."
	go test -v -count=1 -short ./...

test-integration:
	@echo "Running integration tests..."
	go test -v -count=1 -tags=integration ./...

test-e2e:
	@echo "Running E2E tests..."
	@echo "Note: Ensure services are running first"
	go test -v -count=1 ./filer/e2e/...

# Test specific packages
test-user:
	@echo "Running user service tests..."
	go test -v -count=1 ./userapp/...

test-filer:
	@echo "Running filer service tests..."
	go test -v -count=1 ./filer/...

test-vectorlayer:
	@echo "Running vectorlayer service tests..."
	go test -v -count=1 ./vectorlayerapp/...

test-pkg:
	@echo "Running package tests..."
	go test -v -count=1 ./pkg/...

# Clean test artifacts
test-clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html
	find . -name "*.test" -delete
