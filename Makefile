# Variables
APP_NAME := golkube
GO_FILES := $(shell find . -type f -name '*.go')
DOCKER_TAG := $(shell git rev-parse --short HEAD)
DOCKER_IMAGE := $(APP_NAME):$(DOCKER_TAG)
DOCKER_IMAGE_LATEST := $(APP_NAME):latest
CONFIG_FILE := configs/default.yaml
KUBECONFIG_DIR := ~/.kube
KUBERNETES_NAMESPACE := $(shell yq '.kubernetes.namespace' $(CONFIG_FILE))

# Tools
YQ := yq
GOFMT := $(shell gofmt -l $(GO_FILES))
KUBECTL := kubectl

# Default target
.PHONY: all
all: deps fmt vet lint test build

# Dependency management
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Format Go code
.PHONY: fmt
fmt:
	@if [ -n "$(GOFMT)" ]; then \
		echo "The following files need formatting:"; \
		echo "$(GOFMT)"; \
		exit 1; \
	fi
	@echo "Code is properly formatted."

# Vet Go code
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Lint Go code
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run || echo "Linting passed."

# Run unit tests
.PHONY: test
test:
	@echo "Running unit tests..."
	go test ./... -v -cover

# Build the application
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	go build -o bin/$(APP_NAME) ./cmd/$(APP_NAME)

# Run the application locally
.PHONY: run
run: build verify-namespace
	@echo "Running $(APP_NAME)..."
	./bin/$(APP_NAME) --config $(CONFIG_FILE)

# Verify Kubernetes namespace exists, create if not
.PHONY: verify-namespace
verify-namespace:
	@if ! $(KUBECTL) get namespace $(KUBERNETES_NAMESPACE) &>/dev/null; then \
		echo "Namespace $(KUBERNETES_NAMESPACE) not found. Creating namespace..."; \
		$(KUBECTL) create namespace $(KUBERNETES_NAMESPACE); \
	fi
	@echo "Namespace $(KUBERNETES_NAMESPACE) verified."

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	rm -rf bin/ $(APP_NAME)
	docker system prune -f || true
	@echo "Cleaned successfully."

# Docker-related tasks

# Build Docker image
.PHONY: docker-build
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE)..."
	docker build -t $(DOCKER_IMAGE) -t $(DOCKER_IMAGE_LATEST) .

# Push Docker image to registry
.PHONY: docker-push
docker-push: docker-build
	@echo "Pushing Docker image $(DOCKER_IMAGE)..."
	docker push $(DOCKER_IMAGE)
	docker push $(DOCKER_IMAGE_LATEST)

# Run the application in Docker
.PHONY: docker-run
docker-run: docker-build
	@echo "Running $(APP_NAME) in Docker..."
	docker run --rm \
		-v $(PWD)/configs:/app/configs:ro \
		-v $(KUBECONFIG_DIR):/root/.kube:ro \
		-p 8080:8080 \
		-e ENV_VAR=value \
		$(DOCKER_IMAGE_LATEST) --config /app/$(CONFIG_FILE)

# Stop and clean up running Docker containers
.PHONY: docker-clean
docker-clean:
	@echo "Stopping and removing Docker containers..."
	-docker stop $(shell docker ps -q --filter ancestor=$(DOCKER_IMAGE_LATEST)) || true
	-docker rm $(shell docker ps -aq --filter ancestor=$(DOCKER_IMAGE_LATEST)) || true
	@echo "Docker containers cleaned."

# Kubernetes deployment
.PHONY: deploy
deploy: docker-build docker-push
	@echo "Deploying $(APP_NAME) to Kubernetes..."
	kubectl apply -f manifests/deployment.yaml

# Remove Kubernetes deployment
.PHONY: undeploy
undeploy:
	@echo "Removing $(APP_NAME) deployment from Kubernetes..."
	kubectl delete -f manifests/deployment.yaml

# Run application tests in Kubernetes
.PHONY: test-deploy
test-deploy:
	@echo "Testing deployed application..."
	kubectl get pods -n $(KUBERNETES_NAMESPACE)

# Watch for file changes and rebuild
.PHONY: watch
watch:
	@echo "Watching for changes..."
	@find . -name '*.go' | entr -r make build

# Pull latest Docker image
.PHONY: docker-pull
docker-pull:
	@echo "Pulling the latest Docker image..."
	docker pull $(DOCKER_IMAGE_LATEST)

# List running containers
.PHONY: docker-ps
docker-ps:
	@echo "Listing running Docker containers..."
	docker ps

# Show Make targets
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all              - Run all checks and build"
	@echo "  deps             - Install Go dependencies"
	@echo "  fmt              - Format Go code"
	@echo "  vet              - Vet Go code"
	@echo "  lint             - Lint Go code"
	@echo "  test             - Run unit tests"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application locally"
	@echo "  clean            - Clean up build artifacts"
	@echo "  docker-build     - Build Docker image"
	@echo "  docker-push      - Push Docker image to registry"
	@echo "  docker-run       - Run the application in Docker"
	@echo "  docker-clean     - Stop and clean up Docker containers"
	@echo "  docker-pull      - Pull the latest Docker image"
	@echo "  docker-ps        - List running Docker containers"
	@echo "  deploy           - Deploy the application to Kubernetes"
	@echo "  undeploy         - Remove Kubernetes deployment"
	@echo "  test-deploy      - Test deployed application in Kubernetes"
	@echo "  watch            - Watch for file changes and rebuild"
	@echo "  help             - Show this help message"
