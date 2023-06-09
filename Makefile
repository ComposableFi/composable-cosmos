#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.0.0-rc8
BUILDDIR ?= $(CURDIR)/build
HTTPS_GIT := https://github.com/notional-labs/composable-testnet/.git

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

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc cleveldb
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=centauri \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=centaurid \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" 

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif


all: install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/centaurid

build:
	go build $(BUILD_FLAGS) -o bin/centaurid ./cmd/centaurid

docker-build-debug:
	@DOCKER_BUILDKIT=1 docker build -t centauri:debug -f Dockerfile .

###############################################################################
###                                  Proto                                  ###
###############################################################################

protoVer=0.11.6
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
containerProtoGen=proto-gen-$(protoVer)
containerProtoFmt=proto-fmt-$(protoVer)

proto-all: proto-format proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoGen}$$"; then docker start -a $(containerProtoGen); else docker run --name $(containerProtoGen) -v $(CURDIR):/workspace --workdir /workspace $(protoImageName) \
		sh ./scripts/protocgen.sh; fi

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${containerProtoFmt}$$"; then docker start -a $(containerProtoFmt); else docker run --name $(containerProtoFmt) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./ -not -path "./third_party/*" -name "*.proto" -exec clang-format -i {} \; ; fi

proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against $(HTTPS_GIT)#branch=main

.PHONY: proto-all proto-gen proto-format proto-lint proto-check-breaking 

###############################################################################
###                                 Localnet                                ###
###############################################################################
localnet-keys:
	./tests/localcentauri/scripts/add_keys.sh

localnet-init: localnet-clean localnet-build

localnet-build:  
	@chmod -R +x tests/localcentauri/
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker compose -f tests/localcentauri/docker-compose.yml build

localnet-start:  
	@STATE="" docker compose -f tests/localcentauri/docker-compose.yml up

localnet-start-with-state:	
	@STATE=-s docker compose -f tests/localcentauri/docker-compose.yml up

localnet-startd:
	@STATE="" docker compose -f tests/localcentauri/docker-compose.yml up -d

localnet-startd-with-state:
	@STATE=-s docker compose -f tests/localcentauri/docker-compose.yml up -d

localnet-stop:
	@STATE="" docker compose -f tests/localcentauri/docker-compose.yml down

localnet-clean:
	@rm -rfI $(HOME)/.centauri/

localnet-state-export-init: localnet-clean localnet-state-export-build 

localnet-state-export-build:
	@chmod -R +x tests/localcentauri/
	@DOCKER_BUILDKIT=1 COMPOSE_DOCKER_CLI_BUILD=1 docker compose -f tests/localcentauri/state_export/docker-compose.yml build

localnet-state-export-start:
	@docker compose -f tests/localcentauri/state_export/docker-compose.yml up

localnet-state-export-startd:
	@docker compose -f tests/localcentauri/state_export/docker-compose.yml up -d

localnet-state-export-stop:
	@docker compose -f tests/localcentauri/state_export/docker-compose.yml down

###############################################################################
###                             Interchain test                             ###
###############################################################################

# Executes start chain tests via interchaintest
ictest-start-cosmos:
	cd tests/interchaintest && go test -race -v -run TestStartCentauri .

# Executes start chain tests via interchaintest
ictest-start-polkadot:
	cd tests/interchaintest && go test -timeout=25m -race -v -run TestPolkadotComposableChainStart .

# Executes IBC tests via interchaintest
ictest-ibc:
	cd tests/interchaintest && go test -timeout=25m -race -v -run TestCentauriPicassoIBCTransfer .

# Executes all tests via interchaintest after compling a local image as centauri:debug
ictest-all: ictest-start-cosmos ictest-start-polkadot ictest-ibc

# Executes push wasm client tests via interchaintest
ictest-push-wasm:
	cd tests/interchaintest && go test -race -v -run TestPushWasmClientCode .

.PHONY: ictest-start-cosmos ictest-start-polkadot ictest-ibc ictest-push-wasm ictest-all
