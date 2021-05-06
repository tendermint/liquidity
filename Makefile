#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
BINDIR ?= $(GOPATH)/bin
DOCKER := $(shell which docker)
BUILDDIR ?= $(CURDIR)/build
SIMAPP = ./app

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=liquidityd \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=liquidity \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cosmos/cosmos-sdk/types.reDnmString=[a-zA-Z][a-zA-Z0-9/:]{2,127}

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: tools install lint

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
ifeq ($(OS),Windows_NT)
	go build $(BUILD_FLAGS) -o build/liquidityd.exe ./cmd/liquidityd
else
	go build $(BUILD_FLAGS) -o build/liquidityd ./cmd/liquidityd
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

install: go.sum
	go install $(BUILD_FLAGS) ./cmd/liquidityd

build-reproducible: go.sum
	$(DOCKER) rm latest-build || true
	$(DOCKER) run --volume=$(CURDIR):/sources:ro \
        --env TARGET_PLATFORMS='linux/amd64 darwin/amd64 linux/arm64' \
        --env APP=liquidityd \
        --env VERSION=$(VERSION) \
        --env COMMIT=$(COMMIT) \
        --env LEDGER_ENABLED=$(LEDGER_ENABLED) \
        --name latest-build cosmossdk/rbuilder:latest
	$(DOCKER) cp -a latest-build:/home/builder/artifacts/ $(CURDIR)/

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify
	@go mod tidy

clean:
	rm -rf build/

.PHONY: go-mod-cache clean

###############################################################################
###                              Documentation                              ###
###############################################################################

update-swagger-docs: statik
	$(BINDIR)/statik -src=client/docs/swagger-ui -dest=client/docs -f -m
	@if [ -n "$(git status --porcelain)" ]; then \
        echo "Swagger docs are out of sync";\
        exit 1;\
    else \
    	echo "Swagger docs are in sync";\
    fi

.PHONY: update-swagger-docs

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test: test-unit
test-all: test-unit test-race test-cover

test-unit: 
	@VERSION=$(VERSION) go test -mod=readonly -tags='norace' $(PACKAGES_NOSIMULATION)

test-race:
	@go test -mod=readonly -timeout 30m -race -coverprofile=coverage.txt -covermode=atomic -tags='ledger test_ledger_mock' ./...

test-cover:
	@go test -mod=readonly -timeout 30m -coverprofile=coverage.txt -covermode=atomic -tags='norace ledger test_ledger_mock' ./...

test-build: build
	@go test -mod=readonly -p 4 `go list ./cli_test/...` -tags=cli_test -v

benchmark:
	@go test -mod=readonly -bench=. ./...

.PHONY: test test-all test-unit test-race test-cover test-build

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=1000 -Commit=true -Period=0 -v -timeout 24h

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=1 -SimAppPkg=$(SIMAPP) -ExitOnFail 50 5 TestAppImportExport

test-sim-after-import: runsim
	@echo "Running application simulation-after-import. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) -ExitOnFail 1 1 TestAppSimulationAfterImport

.PHONY: \
test-sim-nondeterminism \
test-sim-import-export \
test-sim-after-import 

SIM_NUM_BLOCKS ?= 50
SIM_BLOCK_SIZE ?= 50
SIM_COMMIT ?= true

test-sim-benchmark:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$  \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h

test-sim-profile:
	@echo "Running application benchmark for numBlocks=$(SIM_NUM_BLOCKS), blockSize=$(SIM_BLOCK_SIZE). This may take awhile!"
	@go test -mod=readonly -benchmem -run=^$$ $(SIMAPP) -bench ^BenchmarkFullAppSimulation$$ \
		-Enabled=true -NumBlocks=$(SIM_NUM_BLOCKS) -BlockSize=$(SIM_BLOCK_SIZE) -Commit=$(SIM_COMMIT) -timeout 24h -cpuprofile cpu.out -memprofile mem.out

.PHONY: test-sim-profile test-sim-benchmark

###############################################################################
###                               Localnet                                  ###
###############################################################################

# Run a single testnet locally
localnet: 
	./scripts/localnet.sh

.PHONY: localnet

###############################################################################
###                                Linting                                  ###
###############################################################################

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

###############################################################################
###                                Protobuf                                 ###
###############################################################################

proto-all: proto-gen proto-swagger-gen update-swagger-docs

proto-gen:
	docker run --rm -v $(CURDIR):/workspace --workdir /workspace bharvest/liquidity-proto-gen sh ./scripts/protocgen.sh
	go mod tidy

proto-swagger-gen:
	docker run --rm -v $(CURDIR):/workspace --workdir /workspace bharvest/liquidity-proto-gen sh ./scripts/protoc-swagger-gen.sh

proto-js-gen:
	starport build --rebuild-proto-once

.PHONY: proto-all proto-gen proto-swagger-gen proto-js-gen
