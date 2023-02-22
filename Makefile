#!/usr/bin/make -f


export GO111MODULE = on

BUILD_DIR ?= $(CURDIR)/build

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')

ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif


###############################################################################
###                                   All                                   ###
###############################################################################

all: lint build

###############################################################################
###                                Build flags                              ###
###############################################################################

LD_FLAGS = -X github.com/cheqd/did-resolver/cmd.version=$(VERSION) \
	-X github.com/cheqd/did-resolver/cmd.Commit=$(COMMIT)
BUILD_FLAGS :=  -ldflags '$(LD_FLAGS)'

ifeq ($(LINK_STATICALLY),true)
  LD_FLAGS += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS :=  -ldflags '$(LD_FLAGS)' -tags "$(build_tags)"

###############################################################################
###                                  Build                                  ###
###############################################################################

build: go.sum
	@echo "Building DID Resolver binary..."
	@go build -mod=readonly $(BUILD_FLAGS) -o build/did-resolver main.go
.PHONY: build

###############################################################################
###                                 Install                                 ###
###############################################################################

install: go.sum
	@echo "Installing DID Resolver dependencies..."
	@go install -mod=readonly
.PHONY: install

###############################################################################
###                               Go Version                                ###
###############################################################################

GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MIN_GO_MAJOR_VERSION = 1
MIN_GO_MINOR_VERSION = 18
GO_VERSION_ERROR = Golang version $(GO_MAJOR_VERSION).$(GO_MINOR_VERSION) is not supported, \
please update to at least $(MIN_GO_MAJOR_VERSION).$(MIN_GO_MINOR_VERSION)

go-version:
	@echo "Verifying go version..."
	@if [ $(GO_MAJOR_VERSION) -gt $(MIN_GO_MAJOR_VERSION) ]; then \
		exit 0; \
	elif [ $(GO_MAJOR_VERSION) -lt $(MIN_GO_MAJOR_VERSION) ]; then \
		echo $(GO_VERSION_ERROR); \
		exit 1; \
	elif [ $(GO_MINOR_VERSION) -lt $(MIN_GO_MINOR_VERSION) ]; then \
		echo $(GO_VERSION_ERROR); \
		exit 1; \
	fi

.PHONY: go-version

###############################################################################
###                               Go Modules                                ###
###############################################################################

go.sum: go.mod
	@echo "Ensuring app dependencies have not been modified..."
	go mod verify
	go mod tidy

verify:
	@echo "Verifying all go module dependencies..."
	@find . -name 'go.mod' -type f -execdir go mod verify \;

tidy:
	@echo "Cleaning up all go module dependencies..."
	@find . -name 'go.mod' -type f -execdir go mod tidy \;

.PHONY: verify tidy

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

lint:
	golangci-lint run --out-format=tab

lint-fix:
	golangci-lint run --fix --out-format=tab --issues-exit-code=0
.PHONY: lint lint-fix

format:
	find . -name '*.go' -type f -not -path "*.git*" | xargs gofmt -w -s
	find . -name '*.go' -type f -not -path "*.git*" | xargs misspell -w
	find . -name '*.go' -type f -not -path "*.git*" | xargs goimports -w -local github.com/cheqd/bdjuno
.PHONY: format

clean:
	rm -rf $(BUILD_DIR)
.PHONY: clean

###############################################################################
###                                Swagger                                  ###
###############################################################################

swagger:
	swag init -g main.go
.PHONY: swagger
