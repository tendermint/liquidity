#!/usr/bin/make -f

PACKAGES_NOSIMULATION=$(shell go list ./...)
BINDIR ?= $(GOPATH)/bin

export GO111MODULE = on

all: tools lint test

include contrib/devtools/Makefile


########################################
### Dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download
.PHONY: go-mod-cache

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

########################################
### Testing

SIM_NUM_BLOCKS ?= 50
SIM_BLOCK_SIZE ?= 50
SIM_COMMIT ?= true

test: test-unit
test-all: test-unit test-race test-cover

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly $(PACKAGES_NOSIMULATION)

test-race:
	@VERSION=$(VERSION) go test -mod=readonly -race $(PACKAGES_NOSIMULATION)

.PHONY: test test-all test-unit test-race

.PHONY: \

lint:
	$(BINDIR)/golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -d -s
	go mod verify
.PHONY: lint

format:
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs misspell -w
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "*.pb.go" | xargs goimports -w -local github.com/tendermint/liquidity
.PHONY: format

proto-all: proto-tools proto-gen

proto-gen:
	@./scripts/protocgen.sh
