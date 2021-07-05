run: ## Run the service
	go run ./cmd/service/.

build: ## Build the service
	go build -o main ./cmd/service/.

fmt: ## Run go fmt against code
	go fmt ./internal/... ./cmd/...

vet: ## Run go vet against code
	go vet ./internal/... ./cmd/...

test: ## Runs the tests
	go test -cover -race -short -v $(shell go list ./... | grep -v /vendor/ )

openapi_http: ## generate open api resources
	GO111MODULE=off go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen
	oapi-codegen -generate types -o internal/users/ports/openapi_types.gen.go -package ports api/openapi/users.yml
	oapi-codegen -generate chi-server,spec -o internal/users/ports/openapi_api.gen.go -package ports api/openapi/users.yml

help: ## Shows the help
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: test build run fmt vet openapi_http
