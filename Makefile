BINARY=kmgmt

.DEFAULT_GOAL := help

test: ## Run test coverage
	go test -v -cover -covermode=atomic ./...

kmgmt: ## Build the main 'kmgmt' binary
	go build -o ${BINARY} main.go

unittest: ## Run unit tests
	go test -short  ./...

clean: ## Clean the working dir and it's compiled binary
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

run: ## Compile and run the main program
	go run cmd/kmgmt/main.go

list: ## Print the current module's dependencies.
	go list -m all

tidy: ## Remove unused dependencies
	go mod tidy

lint-prepare: ## Install the golangci linter
	@echo "Installing golangci-lint" 
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

lint: ## Run the golangci linter on source code
	./bin/golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...

help: ## Print help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: clean install unittest build run lint-prepare lint list tidy help