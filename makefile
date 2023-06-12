# ==============================================================================
# Define dependencies

KIND            := kindest/node:v1.27.2 # https://hub.docker.com/r/kindest/node/tags
POSTGRES        := postgres:15.3 # https://hub.docker.com/_/postgres

KIND_CLUSTER    := rosie-starter-cluster


# ==============================================================================
# Install dependencies

dev-brew:
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize

dev-docker:
	docker pull $(KIND)
	docker pull $(POSTGRES)

# ==============================================================================
# Building containers



# ==============================================================================
# Running from within k8s/kind

dev-up-local:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-down-local:
	kind delete cluster --name $(KIND_CLUSTER)

# ==============================================================================
# Modules support

tidy:
	rm -rf vendor
	go mod tidy
	go mod vendor
