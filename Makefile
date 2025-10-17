# Generating code
gen:
	@echo "Generating code..."
	@cd config && sqlc generate
	@echo "Code generated successfully"

# Install dependencies
deps: gen
	@echo "Installing dependencies..."
	@go generate ./...
	@go mod download
	@go mod tidy

# Build the application
ARTIFACT_VERSION ?= 0.0.0-local
build: gen deps
	@echo "Building..."
	@go build \
	-o ./bin/end-user-info \
	./cmd/end-user-info

# Run the application
run: build
	@echo "Running binary..."
	./bin/end-user-info

# Lint the application
lint:
	@echo "Linting..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run --tests=false --disable-all --timeout=2m -p error

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin

 help:
	@echo "Available commands:"
	@echo "  deps    		- Install dependencies"
	@echo "  build   		- Build the application"
	@echo "  build-image		- Build the docker image"
	@echo "  run     		- Run the application"
	@echo "  lint    		- Lint the application"
	@echo "  test    		- Test the application"
	@echo "  clean   		- Clean the binary"

.DEFAULT_GOAL := help
.PHONY: help build build-image run lint test clean
