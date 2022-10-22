ifneq (,$(wildcard ./.env))
    include .env
    export
endif

#--------------------------------------------------------------------------------------------------
##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

#--------------------------------------------------------------------------------------------------
##@ Tools

DEV_DEPENDENCIES := $(shell go list -f "{{range .Imports}}{{.}} {{end}}" dev/tools.go)
.PHONY: install-dev-tools
install-dev-tools: ## Install useful development tools.
	@go install -mod=readonly $(DEV_DEPENDENCIES)

PROTO_DEPENDENCIES := $(shell go list -f "{{range .Imports}}{{.}} {{end}}" grpc/dev/tools.go)
.PHONY: install-proto-tools
install-proto-tools: ## Install tools to generate Protobuf/gRPC code.
	@go install -mod=readonly $(PROTO_DEPENDENCIES)

#--------------------------------------------------------------------------------------------------
##@ Code Generation

.PHONY: grpc
grpc: ## Generate client and server code for the gRPC interfaces.
	rm -rf $(CURDIR)/grpc/gen
	cd $(CURDIR)/grpc/schema && buf generate

#--------------------------------------------------------------------------------------------------
##@ Testing

.PHONY: compose
compose: ## Start local dependencies (Kafka) via Docker Compose.
	@docker compose --file dev/docker-compose.yml up --wait

.PHONY: compose-down
compose-down: ## Remove all running local dependencies.
	@docker compose --file dev/docker-compose.yml down
