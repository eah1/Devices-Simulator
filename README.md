# Devices-Simulator

Device and metric simulator

## Index

* [Development](#development)

## Development <a name="development"></a>

When developing you can use both kind or local environment

### Run local

    make tidy
    make test
    make run-myc-devices-simulator

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
| &nbsp;&nbsp; /sys                               | Common libraries and helpers (sentry, logger, ...)           |
| &nbsp;&nbsp; /web                               | Web specific libraries (middlewares, responses, models, ...) |
| /foundation                                     | Common libraries        (mainly for testing)                 |
| /vendor                                         | Go Vendor folder                                             |
