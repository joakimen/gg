ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

BIN=./bin/clone
MAINPRG=./cmd/clone

.PHONY: all
all: build

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...

.PHONY: build
build: fmt lint vet
	go build -o $(BIN) $(MAINPRG)

.PHONY: run
run: fmt vet
	@go run .

bump-tag:
	@command -v svu >/dev/null 2>&1 || { echo >&2 "svu is not installed. Aborting."; exit 1; }
	@next_tag=$$(svu next) && git tag $$next_tag
	git push && git push --tags

.PHONY: clean
clean:
	rm -rf $(BIN)
