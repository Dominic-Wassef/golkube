# Variables
APP_NAME := golkube
GO_FILES := $(shell find . -type f -name '*.go')
GOFMT := $(shell gofmt -l $(GO_FILES))
DOCKER_IMAGE := $(APP_NAME):latest
CONFIG_FILE := configs/default.yaml

# Default target
.PHONY: all
all: fmt vet build

# Format Go code
.PHONY: fmt
fmt:
	@if [ -n "$(GOFMT)" ]; then \
		echo "The following files need formatting:"; \
		echo "$(GOFMT)"; \
		exit 1; \
	fi

# Vet Go code
.PHONY: vet
vet:
	go vet ./...

# Build the application
.PHONY: build
build:
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

# Run the application
.PHONY: run
run: build
	./bin/$(APP_NAME) --config $(CONFIG_FILE)

# Test the application
.PHONY: test
test:
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/*

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Push Docker image
.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE)

# Lint Go code
.PHONY: lint
lint:
	golangci-lint run ./...

# Generate code (e.g., mocks)
.PHONY: generate
generate:
	go generate ./...

# Install dependencies
.PHONY: deps
deps:
	go mod tidy
	go mod download

# Update dependencies
.PHONY: update-deps
update-deps:
	go get -u ./...
	go mod tidy

# Run in Docker
.PHONY: docker-run
docker-run: docker-build
	docker run --rm -v $(PWD)/configs:/app/configs -p 8080:8080 $(DOCKER_IMAGE) --config /app/$(CONFIG_FILE)
