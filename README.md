# Devices-Simulator

Device and metric simulator

## Index

* [Arquitecture](#arquitecture)
* [Getting Started](#getting-started)
* [Development](#development)

## Arquitechture <a name="architecture"></a>

Devices-Simulator exposes the following services:

* simulator-api: The API server.

Each service has its own:

- Docker Image
- Kubernetes Deployment

Which can be found under `deployments` directory.

Internally, Devices-Simulator has the following components:

* business/usecase: The use case bussines application.
* business/core: The core business logic.
* business/db: The persistence logic and database schema

## Getting Started <a name="getting-started"></a>

First of all, ensure that you environment has all the following applications installed

Required (mandatory):

- Go 1.19.1 [releases notes](https://tip.golang.org/doc/go1.19)
- [Makefile](https://www.gnu.org/software/make/manual/make.html)
- [Docker](https://www.docker.com/products/docker-desktop/)
- [Golanci](https://golangci-lint.run/usage/install/)
- [Goose](https://pressly.github.io/goose/installation/)

To update docs (optional):

- [swag](https://github.com/swaggo/swag)

To use kind (optional):

- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

### Local

Start a database that listens on port 5430 (kind database is listening on 5432)

    make start-postgres-test

Update schema:

    make goose-up POSTGRES_URI=MYC_DEVICES_SIMULATOR

Run tests:

* First we need to make sure we have exported the following environment variables:

      export MYC_DEVICES_SIMULATOR_DBPOSTGRES=postgres://postgres:postgres@localhost:5430/postgres?sslmode=disable

* Now we can run the test assigning the previous variables:

      make test

Set a a configuration file:

    cp config.yaml.sample config.yaml

`config.yaml` is in `.gitignore` so it will never be uploaded to the repo. You can edit it to your needs. Ask for a
Postmark token to your Team Leader.

Run the API

    make run-simulator-api

### Using Kind

Prepare environment:

    make kind-up
    make all
    make kind-load
    make kind-apply-bd
    make kind-db-migration
    make kind-apply-simulator-api

Please not that postgres inside kind is mapped to port 5432, and Postgres used for local tests is mapped in port 5430.

Now you can check: http://localhost:2323/api/swagger/index.html

You can check logs with:

    make kind-logs-simulator-api

## Development <a name="development"></a>

When developing you can use both kind or local environment

### Run local

    make tidy
    make test POSTGRES_URI=MYC_DEVICES_SIMULATOR
    make start-postgres-test
    make goose-up POSTGRES_URI=$MYC_DEVICES_SIMULATOR_DBPOSTGRES
    make run-myc-devices-simulator

### Run kind

    make tidy
    make all
    make kind-up
    make kind-load
    make kind-apply-db
    make kind-db-migration
    make kind-apply-myc-api

### Update Swagger

Before any commit remember to update swagger documentation

    make swagger

### Code organization

| folder                                          | content                                                      |
|-------------------------------------------------|--------------------------------------------------------------|
| /app                                            | Applications                                                 | 
| &nbsp;&nbsp; /services                          | Services                                                     |
| &nbsp;&nbsp;&nbsp;&nbsp; /myc-devices-simulator | Web API                                                      |
| /business                                       | Business Logic                                               |
| &nbsp;&nbsp; /usecase                           | Use case application                                         |
| &nbsp;&nbsp; /core                              | Core business logic                                          |
| &nbsp;&nbsp;&nbsp;&nbsp; /models                | Core Models                                                  |
| &nbsp;&nbsp; /db                                | Database and persistence logic                               |
| &nbsp;&nbsp;&nbsp;&nbsp; /schema                | Database schema                                              |
| &nbsp;&nbsp;&nbsp;&nbsp; /store                 | Database persistence logic                                   |
| &nbsp;&nbsp; /sys                               | Common libraries and helpers (sentry, logger, ...)           |
| &nbsp;&nbsp; /web                               | Web specific libraries (middlewares, responses, models, ...) |
| &nbsp;&nbsp; /deploy                            | Deployment scripts and code                                  |
| &nbsp;&nbsp;&nbsp;&nbsp; /docker                | Docker files                                                 |
| &nbsp;&nbsp;&nbsp;&nbsp; /k8s                   | K8S logic (kind, kustomization, charts, ...)                 |
| /foundation                                     | Common libraries        (mainly for testing)                 |
| /vendor                                         | Go Vendor folder                                             |
