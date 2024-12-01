# Variables
APP_NAME := golkube
GO_FILES := $(shell find . -type f -name '*.go')
DOCKER_TAG := $(shell git rev-parse --short HEAD)
DOCKER_IMAGE := $(APP_NAME):$(DOCKER_TAG)
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
run: build verify-namespace
	./bin/$(APP_NAME) --config $(CONFIG_FILE)

# Verify Kubernetes namespace
.PHONY: verify-namespace
verify-namespace:
	kubectl get namespace $(shell yq '.kubernetes.namespace' $(CONFIG_FILE)) || \
	kubectl create namespace $(shell yq '.kubernetes.namespace' $(CONFIG_FILE))

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Push Docker image
.PHONY: docker-push
docker-push:
	docker push $(DOCKER_IMAGE)

# Run in Docker
.PHONY: docker-run
docker-run: docker-build
	docker run --rm -v $(PWD)/configs:/app/configs -p 8080:8080 $(DOCKER_IMAGE) --config /app/$(CONFIG_FILE)
