// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/health": {
            "get": {
                "description": "Get status service to be alive",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Health"
                ],
                "summary": "Health check EndPoint",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.SuccessHealth"
                        }
                    }
                }
            }
        },
        "/api/v1/auth/login": {
            "post": {
                "description": "Login authentication user in platform",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login user EndPoint",
                "parameters": [
                    {
                        "description": "LoginUser",
                        "name": "User",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webmodels.LoginUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.SuccessLogin"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Validator"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Failed"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/responses.Failed"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Failed"
                        }
                    }
                }
            }
        },
        "/api/v1/users": {
            "post": {
                "description": "Create a new user in the system.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Create user registration EndPoint",
                "parameters": [
                    {
                        "description": "UserRegister",
                        "name": "UserRegister",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webmodels.RegisterUser"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/responses.Success"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Validator"
                        }
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {
                            "$ref": "#/definitions/responses.Failed"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Failed"
                        }
                    }
                }
            }
        },
        "/api/v1/users/activate/{activateToken}": {
            "post": {
                "description": "Activation a user in the system.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Activate user activation EndPoint",
                "parameters": [
                    {
                        "type": "string",
                        "description": "ActivateToken",
                        "name": "activateToken",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/responses.Success"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/responses.Validator"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/responses.Failed"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/responses.Failed"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "responses.Failed": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string",
                    "example": "ERROR"
                }
            }
        },
        "responses.Health": {
            "type": "object",
            "properties": {
                "buildVersion": {
                    "type": "string",
                    "example": "localhost"
                }
            }
        },
        "responses.Success": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                }
            }
        },
        "responses.SuccessHealth": {
            "type": "object",
            "properties": {
                "health": {
                    "$ref": "#/definitions/responses.Health"
                },
                "status": {
                    "type": "string",
                    "example": "OK"
                }
            }
        },
        "responses.SuccessLogin": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string",
                    "example": "OK"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "responses.Validator": {
            "type": "object",
            "properties": {
                "details": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "status": {
                    "type": "string",
                    "example": "ERROR"
                }
            }
        },
        "webmodels.LoginUser": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "webmodels.RegisterUser": {
            "type": "object",
            "required": [
                "company",
                "email",
                "firstName",
                "language",
                "lastName",
                "password"
            ],
            "properties": {
                "company": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "language": {
                    "type": "string",
                    "maxLength": 2,
                    "enum": [
                        "es",
                        "en",
                        "fr",
                        "pt"
                    ]
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 8
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Swagger MYC-DEVICE-SIMULATOR API",
	Description:      "Devices Simulator documentation API.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
