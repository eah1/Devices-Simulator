# Application
VERSION := 1.0
ARCH :=  $(shell ./get-go-arch.sh)

KIND_CLUSTER := devices-simulator-cluster
DEVICES_SIMULATOR_API_IMAGE_NAME := simulator-api
DEVICES_SIMULATOR_QUEUE_IMAGE_NAME := simulator-queue

SCHEMA_DIR := business/db/schema
POSTGRES_URI := postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable&timezone=utc
POSTGRES_KIND_URI := postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&timezone=utc

# Generate vendor
tidy:
	go mod tidy
	go mod vendor

# Linter and Test.
lint:
	golangci-lint version
	golangci-lint linters
	golangci-lint run

# BBDD.
start-postgres-test:
	docker run --name postgresTest -e POSTGRES_PASSWORD=postgres -p 5430:5432  -d postgres
	echo "POSTGRES_URI=\"postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable\""
	sleep 3

stop-postgres-test:
	docker stop postgresTest
	docker rm postgresTest

# REDIS.
start-redis-test:
	docker run --name redisTest -p 6379:6379 -d redis
	sleep 3
stop-redis-test:
	docker stop redisTest
	docker rm redisTest

# RUN SERVICE FAKEMAILER
start-fakemailer-test:
	docker pull circutor/fakemailer:latest
	docker run --name fakemailer -p 25:25  -p 4301:4301 -d circutor/fakemailer:latest
stop-fakemailer-test:
	docker stop fakemailer
	docker rm fakemailer

# Goose Postgres.
goose-status:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_URI)" goose -dir "$(SCHEMA_DIR)" status
goose-up:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_URI)" goose -dir "$(SCHEMA_DIR)" up
goose-down:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_URI)" goose -dir "$(SCHEMA_DIR)" down

# Test.
test:
	MYC_DEVICES_SIMULATOR_DBPOSTGRES=$(POSTGRES_URI) go test -coverprofile=profile.cov ./... -p 2
	go tool cover -func profile.cov
	go vet ./...
	gofmt -l .

test-local: start-postgres-test start-redis-test start-fakemailer-test goose-up test stop-postgres-test stop-redis-test stop-fakemailer-test

# Run Local.
run-simulator-api:
	go run app/services/simulator-api/main.go -f "configFile=config.yaml"
run-simulator-queue:
	go run app/services/simulator-queue/main.go -f "configFile=config.yaml"

# Build Local.
build-simulator-api:
	go build app/services/simulator-api/main.go

# Swagger documentation.
swagger:
	swag init --dir "app/services/simulator-api/handlers"  --output "app/services/simulator-api/docs"  --generalInfo handlers.go  --parseDependency true

# Build docker.
all: simulator-api simulator-queue

simulator-api:
	DOCKER_BUILDKIT=1 docker build \
		-f deploy/docker/simulator-api.dockerfile  \
		-t $(DEVICES_SIMULATOR_API_IMAGE_NAME):$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg GOARCH=$(ARCH) \
		.
simulator-queue:
	DOCKER_BUILDKIT=1 docker build \
    	-f deploy/docker/simulator-queue.dockerfile  \
    	-t $(DEVICES_SIMULATOR_QUEUE_IMAGE_NAME):$(VERSION) \
    	--build-arg BUILD_REF=$(VERSION) \
    	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
    	--build-arg GOARCH=$(ARCH) \
    	.

##################          Kind created - delete        ##################
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

##################          Kind status service          ##################
kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-service:
	kubectl get pods -o wide --watch -n device-simulator-system
kind-status-db:
	kubectl get pods -o wide --watch --namespace=database-system

kind-ingress-check:
	kubectl wait --namespace ingress-nginx \
	  --for=condition=ready pod \
	  --selector=app.kubernetes.io/component=controller \
	  --timeout=90s

##################          Kind load          ##################
kind-load: kind-load-simulator-api kind-load-simulator-queue
kind-load-simulator-api:
	kind load docker-image $(DEVICES_SIMULATOR_API_IMAGE_NAME):$(VERSION) --name $(KIND_CLUSTER)
kind-load-simulator-queue:
	kind load docker-image $(DEVICES_SIMULATOR_QUEUE_IMAGE_NAME):$(VERSION) --name $(KIND_CLUSTER)

##################          Kind apply          ##################
kind-apply-db:
	kustomize build deploy/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
kind-apply-redis:
	kustomize build deploy/k8s/kind/redis-pod | kubectl apply -f -
	kubectl wait --namespace=redis-system --timeout=120s --for=condition=Available deployment/redis-pod
kind-apply-simulator-api:
	kustomize build deploy/k8s/kind/simulator-api-pod | kubectl apply -f -
kind-apply-simulator-queue:
	kustomize build deploy/k8s/kind/simulator-queue-pod | kubectl apply -f -

##################          Kind restart          ##################
kind-restart:
	kubectl rollout restart deployment simulator-api --namespace=device-simulator-system

kind-update: all kind-load kind-restart
kind-update-apply: all kind-load kind-apply-db kind-apply-redis kind-apply-simulator-api kind-apply-simulator-queue

##################           Kind logs          ##################
kind-logs-db:
	kubectl logs -l app=database --namespace=database-system --all-containers=true -f --tail=100
kind-logs-redis:
	kubectl logs -l app=redis --namespace=redis-system --all-containers=true -f --tail=100
kind-logs-simulator-api:
	kubectl logs -l app=simulator-api --all-containers=true -f --tail=100 --namespace=device-simulator-system
kind-logs-simulator-queue:
	kubectl logs -l app=simulator-queue --all-containers=true -f --tail=100 --namespace=device-simulator-system

##################           Kind describe          ##################
kind-describe-db:
	kubectl describe pod -l app=database --namespace=database-system
kind-describe-redis:
	kubectl describe pod -l app=redis --namespace=redis-system
kind-describe-simulator-api:
	kubectl describe pod -l app=simulator-api --namespace=device-simulator-system
kind-describe-simulator-queue:
	kubectl describe pod -l app=simulator-queue --namespace=device-simulator-system

##################           Kind DDBB          ##################
kind-connection-db:
	kubectl exec -it <pod> --namespace=database-system -- psql -h <host> -U <user> --password -p 5432 postgres

kind-db-migration:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_KIND_URI)" goose -dir "$(SCHEMA_DIR)" status
	GOOSE_DRIVER=postgres GOOSE_DBSTRING="$(POSTGRES_KIND_URI)" goose -dir "$(SCHEMA_DIR)" up
