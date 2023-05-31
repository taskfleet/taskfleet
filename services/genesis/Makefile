ifneq (,$(wildcard ./.env))
    include .env
    export
endif

#--------------------------------------------------------------------------------------------------
##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: config
config: check-service.env.yaml ## Update the service configuration to reference machine images.
	@cat $(CURDIR)/config/service.env.yaml \
		| yq ".template.aws.boot.amis[0].config.owner = \"$(shell aws sts get-caller-identity | jq -r .Account)\"" \
		| yq ".template.gcp.boot.imageLinks[0].config = \"$(shell gcloud compute images describe taskfleet-instance-manager-testserver --format json | jq -r .selfLink)\"" \
		| tee $(CURDIR)/config/service.env.yaml

#--------------------------------------------------------------------------------------------------
##@ Development

.PHONY: run-service
run-service: check-service.env.yaml ## Run the instance manager service.
	go run cmd/service/main.go --config config/service.env.yaml

.PHONY: grpcui
grpcui: ## Run a grpcui server to interactively play with the instance manager service.
	@grpcui \
		-import-path $(CURDIR)/../../grpc/schema \
		-proto $(CURDIR)/../../grpc/schema/genesis/v1/service.proto \
		-plaintext \
		localhost:5404


#--------------------------------------------------------------------------------------------------

define ERROR_MESSAGE
!!! ERROR !!!

Missing configuration file 'config/service.env.yaml'!
Generate the file by running 'make test-infrastructure'. This will create required cloud
infrastructure and fill in the template at 'config/service.template.yaml'.
endef

export ERROR_MESSAGE
check-service.env.yaml:
	@if [ ! -f config/service.env.yaml ]; then \
		echo "$$ERROR_MESSAGE"; \
		false; \
	fi
