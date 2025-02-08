ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

BIN=./bin/gg
MAINPRG=./cmd/gg

.PHONY: all
all: build

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: fmt
fmt:
	gofumpt -l -w .
	goimports -w .

.PHONY: lint
lint:
	go vet ./...
	staticcheck ./...

.PHONY: build
build: fmt lint
	go build -o $(BIN) $(MAINPRG)

release:
	@command -v svu >/dev/null 2>&1 || { echo >&2 "svu is not installed. Aborting."; exit 1; }
	@next_tag=$$(svu next) && git tag $$next_tag && echo "git tag $$next_tag"
	git push --tags

.PHONY: clean
clean:
	rm -rf $(BIN)
