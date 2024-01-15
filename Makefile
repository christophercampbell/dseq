ARCH := $(shell arch)

ifeq ($(ARCH),x86_64)
	ARCH = amd64
else
	ifeq ($(ARCH),aarch64)
		ARCH = arm64
	endif
endif

GOBASE := $(shell pwd)
GOOS=$(shell uname -s  | tr '[:upper:]' '[:lower:]')
GOBIN := $(GOBASE)/build
GOCMD := $(GOBASE)/main.go
GOBINARY := dseq

GOENVVARS := GOBIN=$(GOBIN) CGO_ENABLED=0 GOARCH=$(ARCH) GOOS=$(GOOS)
DOCKERVARS := GOBIN=$(GOBIN) CGO_ENABLED=0 GOARCH=amd64 GOOS=linux

COMETBFT := $(shell command -v cometbft 2> /dev/null)
COMETBFT_VERSION := v0.38.0-rc3

build: build-linux ## build the binary
	$(GOENVVARS)  go build -o $(GOBIN)/$(GOBINARY) $(GOCMD)
.PHONY: build

build-linux:
	$(DOCKERVARS) go build -o $(GOBIN)/$(GOBINARY)-linux $(GOCMD)
.PHONY: build-linux

build-docker-local: build-linux ## build docker image
	@cd networks/local && make
.PHONY: build-docker

test: ## run tests
	go test ./...
.PHONY: test

clean: ## clean build artifacts
	rm -rf $(GOBIN)
.PHONY: clean

start: stop build-docker-local ## run a 4-node testnet locally
	@if ! [ -f build/node0/config/genesis.json ]; then cometbft testnet --config networks/local/localnode/config-template.toml --o ./build --starting-ip-address 192.167.10.2; fi
	docker-compose up -d
.PHONY: start

stop: ## stop docker nodes
	docker-compose down
.PHONY: stop

load-10: ## send 10 txs to nodes randomly with concurrency 1
	go run main.go load --nodes localhost:26657,localhost:26660,localhost:26662,localhost:26664 -r 10 -c 3
.PHONY: load-10

load-100: ## send 100 txs to nodes randomly with concurrency 3
	go run main.go load --nodes localhost:26657,localhost:26660,localhost:26662,localhost:26664 -r 100 -c 9
.PHONY: load-100

checksum: ## compare node sequence files
	for i in {0..3}; do md5sum "./build/node$${i}/dseq.bin" | cut -d' ' -f1; done
.PHONY: compare

read-all: ## multi tail consumers
	multitail -l "./build/dseq read --node localhost:6900" -l "./build/dseq read --node localhost:6901" -l "./build/dseq read --node localhost:6902" -l "./build/dseq read --node localhost:6903"
.PHONY: read


.PHONY: help
help: ## prints this help
		@grep -h -E '^[a-zA-Z0-9_-]*:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

