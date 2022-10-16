# Generate vendor
tidy:
	go mod tidy
	go mod vendor

# Test.
test:
	go test -coverprofile=profile.cov ./... -p 2
	go tool cover -func profile.cov
	go vet ./...
	gofmt -l .

# Run Local.
run-myc-devices-simulator:
	go run app/services/myc-devices-simulator/main.go -f "configFile=config.yaml"

# Build Local.
build-myc-devices-simulator:
	go build app/services/myc-devices-simulator/main.go

# Swagger documentation.
swagger:
	swag init --dir "app/services/myc-devices-simulator/handlers"  --output "app/services/myc-devices-simulator/docs"  --generalInfo handlers.go  --parseDependency true