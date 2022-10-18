# Application

VERSION := 1.0
DEVICES_SIMULATOR_API_IMAGE_NAME := simulator-api

ARCH :=  $(shell ./get-go-arch.sh)

KIND_CLUSTER := devices-simulator-cluster

# Generate vendor
tidy:
	go mod tidy
	go mod vendor

# Linter and Test.
lint:
	golangci-lint version
	golangci-lint linters
	golangci-lint run --fix

# Test.
test:
	go test -coverprofile=profile.cov ./... -p 2
	go tool cover -func profile.cov
	go vet ./...
	gofmt -l .

# Run Local.
run-simulator-api:
	go run app/services/simulator-api/main.go -f "configFile=config.yaml"

# Build Local.
build-simulator-api:
	go build app/services/simulator-api/main.go

# Swagger documentation.
swagger:
	swag init --dir "app/services/simulator-api/handlers"  --output "app/services/simulator-api/docs"  --generalInfo handlers.go  --parseDependency true

# Build docker.
all: simulator-api

simulator-api:
	DOCKER_BUILDKIT=1 docker build \
		-f deploy/docker/simulator-api.dockerfile  \
		-t $(DEVICES_SIMULATOR_API_IMAGE_NAME):$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg GOARCH=$(ARCH) \
		.

# Kind Created && Delete.
kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config deploy/k8s/kind/kind-config.yaml
	kubectl create namespace myc-devices-simulator-cluster-system
	kubectl config set-context --current --namespace=myc-devices-simulator-cluster-system
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

# Kind Status service.
kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-service:
	kubectl get pods -o wide --watch -n device-simulator-system

kind-ingress-check:
	kubectl wait --namespace ingress-nginx \
	  --for=condition=ready pod \
	  --selector=app.kubernetes.io/component=controller \
	  --timeout=90s

# Kind Load && Apply app.
kind-load: kind-load-simulator-api
kind-load-simulator-api:
	kind load docker-image $(DEVICES_SIMULATOR_API_IMAGE_NAME):$(VERSION) --name $(KIND_CLUSTER)

kind-apply-simulator-api:
	kustomize build deploy/k8s/kind/simulator-api-pod | kubectl apply -f -

kind-restart:
	kubectl rollout restart deployment simulator-api -n device-simulator-system

kind-update: all kind-load kind-restart
kind-update-apply: all kind-load kind-apply-simulator-api

# Kind logs && describe.
kind-logs-simulator-api:
	kubectl logs -l app=simulator-api --all-containers=true -f --tail=100 -n device-simulator-system

kind-describe-simulator-api:
	kubectl describe pod -l app=simulator-api -n device-simulator-system