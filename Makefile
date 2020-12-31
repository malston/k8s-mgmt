NAME ?= kmgmt
OUTPUT = ./bin/$(NAME)
GO_SOURCES = $(shell find . -type f -name '*.go')
GOBIN ?= $(shell go env GOPATH)/bin
VERSION ?= $(shell cat VERSION)
GITSHA = $(shell git rev-parse HEAD)
GITDIRTY = $(shell git diff --quiet HEAD || echo "dirty")
LDFLAGS_VERSION = -X github.com/malston/k8s-mgmt/pkg/cli.cli_name=$(NAME) \
				  -X github.com/malston/k8s-mgmt/pkg/cli.cli_version=$(VERSION) \
				  -X github.com/malston/k8s-mgmt/pkg/cli.cli_gitsha=$(GITSHA) \
				  -X github.com/malston/k8s-mgmt/pkg/cli.cli_gitdirty=$(GITDIRTY)

.PHONY: all
all: build test verify-goimports ## Build, test, verify source formatting

.PHONY: clean
clean: ## Delete build output
	rm -rf bin/
	rm -rf dist/

.PHONY: build
build: $(OUTPUT) ## Build the main 'kmgmt' binary

.PHONY: test
test: ## Run the tests
	go test ./...

.PHONY: install
install: build ## Copy build to GOPATH/bin
	cp $(OUTPUT) $(GOBIN)

.PHONY: coverage
coverage: ## Run the tests with coverage and race detection
	go test -v --race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: check-goimports
check-goimports: ## Checks if goimports is installed
	@which goimports > /dev/null || (echo goimports not found: issue \"GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports\" && false)

.PHONY: goimports
goimports: check-goimports ## Runs goimports on the project
	@goimports -w pkg cmd

.PHONY: verify-goimports
verify-goimports: check-goimports ## Verifies if all source files are formatted correctly
	@goimports -l pkg cmd | (! grep .) || (echo above files are not formatted correctly. please run \"make goimports\" && false)

$(OUTPUT): $(GO_SOURCES) VERSION
	go build -o $(OUTPUT) -ldflags "$(LDFLAGS_VERSION)" ./cmd/$(NAME)

.PHONY: release
release: $(GO_SOURCES) VERSION ## Cross-compile kmgmt for various operating systems
	@mkdir -p dist
	GOOS=darwin   GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     ./cmd/$(NAME) && tar -czf dist/$(NAME)-darwin-amd64.tgz  -C bin . && rm -f $(OUTPUT)
	GOOS=linux    GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT)     ./cmd/$(NAME) && tar -czf dist/$(NAME)-linux-amd64.tgz   -C bin . && rm -f $(OUTPUT)
	GOOS=windows  GOARCH=amd64 go build -ldflags "$(LDFLAGS_VERSION)" -o $(OUTPUT).exe ./cmd/$(NAME) && zip -rj  dist/$(NAME)-windows-amd64.zip    bin   && rm -f $(OUTPUT).exe

.PHONY: lint-prepare
lint-prepare: ## Install the golangci linter
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: lint
lint: ## Run the golangci linter on source code
	./bin/golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...

.PHONY: tidy
tidy: ## Remove unused dependencies
	go mod tidy

.PHONY: list
list: ## Print the current module's dependencies.
	go list -m all

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
