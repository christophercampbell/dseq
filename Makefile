# Variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Architecture detection
ARCH := $(shell arch)
ifeq ($(ARCH),x86_64)
	ARCH = amd64
else ifeq ($(ARCH),aarch64)
	ARCH = arm64
endif

# Go variables
GOBASE := $(shell pwd)
GOOS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOBIN := $(GOBASE)/build
GOCMD := $(GOBASE)/main.go
GOBINARY := dseq

# Build variables
GOENVVARS := GOBIN=$(GOBIN) CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=$(GOOS)

# Docker variables
DOCKER_IMAGE := dseq
DOCKER_TAG ?= latest
DOCKER_PLATFORMS := linux/amd64,linux/arm64

# CometBFT variables
COMETBFT := $(shell command -v cometbft 2> /dev/null)
COMETBFT_VERSION := v0.38.0-rc3

# Load test variables
NODES ?= localhost:26657,localhost:26660,localhost:26662,localhost:26664
CONCURRENCY ?= 3
REQUESTS ?= 100

# Build targets
.PHONY: build
build: ## Build the binary for current platform
	$(GOENVVARS) go build $(LDFLAGS) -o $(GOBIN)/$(GOBINARY) $(GOCMD)

.PHONY: docker-build
docker-build: ## Build multi-architecture Docker image locally
	@echo "Building multi-architecture Docker image locally..."
	@if ! docker buildx inspect multiarch-builder >/dev/null 2>&1; then \
		echo "Creating new buildx builder..."; \
		docker buildx create --name multiarch-builder --use; \
	fi
	docker buildx build --platform $(DOCKER_PLATFORMS) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--load .

.PHONY: docker-publish
docker-publish: ## Build and publish multi-architecture Docker image
	@echo "Building and publishing multi-architecture Docker image..."
	@if ! docker buildx inspect multiarch-builder >/dev/null 2>&1; then \
		echo "Creating new buildx builder..."; \
		docker buildx create --name multiarch-builder --use; \
	fi
	docker buildx build --platform $(DOCKER_PLATFORMS) \
		-t $(DOCKER_IMAGE):$(DOCKER_TAG) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		--push .

.PHONY: build-docker-local
build-docker-local: ## Build local Docker image
	@cd networks/local && make

# Test targets
.PHONY: test
test: ## Run tests
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: test-coverage
test-coverage: test ## Generate test coverage report
	go tool cover -html=coverage.txt -o coverage.html

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	rm -rf $(GOBIN)
	find . -name "*.test" -type f -delete
	find . -name "coverage.txt" -type f -delete
	find . -name "coverage.html" -type f -delete

# Network targets
.PHONY: start
start: stop build-docker-local ## Run a 4-node testnet locally
	@if ! [ -f $(GOBIN)/node0/config/genesis.json ]; then \
		cometbft testnet --config networks/local/localnode/config-template.toml --o $(GOBIN) --starting-ip-address 192.167.10.2; \
	fi
	docker-compose up -d

.PHONY: stop
stop: ## Stop docker nodes
	docker-compose down

# Load testing targets
.PHONY: load
load: ## Run load test with custom parameters (make load NODES=localhost:26657 REQUESTS=10 CONCURRENCY=1)
	go run main.go load --nodes $(NODES) -r $(REQUESTS) -c $(CONCURRENCY)

.PHONY: load-quick
load-quick: ## Quick load test (10 requests, 3 concurrent)
	@$(MAKE) load REQUESTS=10 CONCURRENCY=3

.PHONY: load-medium
load-medium: ## Medium load test (100 requests, 9 concurrent)
	@$(MAKE) load REQUESTS=100 CONCURRENCY=9

.PHONY: load-heavy
load-heavy: ## Heavy load test (1000 requests, 30 concurrent)
	@$(MAKE) load REQUESTS=1000 CONCURRENCY=30

# Monitoring targets
.PHONY: checksum
checksum: ## Compare node sequence files
	@echo "Comparing sequence files..."
	@for i in {0..3}; do \
		echo "Node $$i: $$(md5sum "$(GOBIN)/node$${i}/dseq.bin" | cut -d' ' -f1)"; \
	done

.PHONY: read-all
read-all: ## Monitor all nodes using multitail
	multitail -l "$(GOBIN)/$(GOBINARY) read --node localhost:6900" \
		-l "$(GOBIN)/$(GOBINARY) read --node localhost:6901" \
		-l "$(GOBIN)/$(GOBINARY) read --node localhost:6902" \
		-l "$(GOBIN)/$(GOBINARY) read --node localhost:6903"

# Development tools
.PHONY: lint
lint: ## Run linters
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

# Help target
.PHONY: help
help: ## Print this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

