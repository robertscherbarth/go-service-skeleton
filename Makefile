run: ## Run the service
	go run .

build: ## Build the service
	go build -o main .

fmt: ## Run go fmt against code
	go fmt ./pkg/... .

vet: ## Run go vet against code
	go vet ./pkg/... .

test: ## Runs the tests
	go test -cover -race -short -v $(shell go list ./... | grep -v /vendor/ )

swagger: ## Generates swagger documentation
	GO111MODULE=off go get -u github.com/swaggo/swag/cmd/swag
	swag init -o ./pkg/docs

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: test build run fmt vet swagger