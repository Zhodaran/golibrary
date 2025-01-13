// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/api/address/geocode": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint allows you to get geo coordinates by latitude and longitude.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "geo"
                ],
                "summary": "Get Geo Coordinates by Latitude and Longitude",
                "parameters": [
                    {
                        "description": "Geographic coordinates",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.GeocodeRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Bearer {token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешное выполнение",
                        "schema": {
                            "$ref": "#/definitions/service.ResponseAddress"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса",
                        "schema": {
                            "$ref": "#/definitions/service.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка подключения к серверу",
                        "schema": {
                            "$ref": "#/definitions/service.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/address/search": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint allows you to get geo coordinates by address.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "geo"
                ],
                "summary": "Get Geo Coordinates by Address",
                "parameters": [
                    {
                        "description": "Address search query",
                        "name": "address",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.RequestAddressSearch"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Bearer {token}",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешное выполнение",
                        "schema": {
                            "$ref": "#/definitions/service.ResponseAddress"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса",
                        "schema": {
                            "$ref": "#/definitions/service.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка подключения к серверу",
                        "schema": {
                            "$ref": "#/definitions/service.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/book": {
            "get": {
                "description": "This description created new SQL user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TakeBook"
                ],
                "summary": "List SQL book",
                "responses": {
                    "200": {
                        "description": "List successful",
                        "schema": {
                            "$ref": "#/definitions/control.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/book/take/{index}": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "This endpoint allows you to get geo coordinates by address.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "TakeBook"
                ],
                "summary": "Get Geo Coordinates by Address",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Book INDEX",
                        "name": "index",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешное выполнение",
                        "schema": {
                            "$ref": "#/definitions/service.ResponseAddress"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка подключения к серверу",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/login": {
            "post": {
                "description": "This endpoint allows a user to log in with their username and password.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Login a user",
                "parameters": [
                    {
                        "description": "User login details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/auth.LoginResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/register": {
            "post": {
                "description": "This endpoint allows you to register a new user with a username and password.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "User registration details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.User"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "User registered successfully",
                        "schema": {
                            "$ref": "#/definitions/auth.TokenResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "User already exists",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/auth.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users": {
            "get": {
                "description": "This description created new SQL user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Controller"
                ],
                "summary": "List SQL user",
                "responses": {
                    "200": {
                        "description": "List successful",
                        "schema": {
                            "$ref": "#/definitions/control.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "This description created new SQL user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Controller"
                ],
                "summary": "Create SQL user",
                "parameters": [
                    {
                        "description": "User login details",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/repository.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Create successful",
                        "schema": {
                            "$ref": "#/definitions/control.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users/{id}": {
            "get": {
                "description": "This description created new SQL user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Controller"
                ],
                "summary": "Get SQL user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Greate successful",
                        "schema": {
                            "$ref": "#/definitions/control.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "This description created new SQL user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Controller"
                ],
                "summary": "Delete SQL user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Delete successful",
                        "schema": {
                            "$ref": "#/definitions/control.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/users{id}": {
            "put": {
                "description": "This description created new SQL user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Controller"
                ],
                "summary": "Update SQL user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "User detail to update",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/repository.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Update successful",
                        "schema": {
                            "$ref": "#/definitions/control.CreateResponse"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Invalid credentials",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/control.rErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.ErrorResponse": {
            "type": "object",
            "properties": {
                "200": {
                    "type": "string"
                },
                "400": {
                    "type": "string"
                },
                "500": {
                    "type": "string"
                }
            }
        },
        "auth.LoginResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "auth.TokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "auth.User": {
            "type": "object",
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "control.CreateResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "control.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Код ошибки",
                    "type": "integer"
                },
                "message": {
                    "description": "Сообщение об ошибке",
                    "type": "string"
                }
            }
        },
        "control.rErrorResponse": {
            "type": "object",
            "properties": {
                "200": {
                    "type": "string"
                },
                "400": {
                    "type": "string"
                },
                "500": {
                    "type": "string"
                }
            }
        },
        "main.mErrorResponse": {
            "type": "object",
            "properties": {
                "200": {
                    "type": "string"
                },
                "400": {
                    "type": "string"
                },
                "500": {
                    "type": "string"
                }
            }
        },
        "repository.User": {
            "type": "object",
            "properties": {
                "deleted_at": {
                    "description": "Для логического удаления",
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "service.Address": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "geo_lat": {
                    "type": "string"
                },
                "geo_lon": {
                    "type": "string"
                },
                "house": {
                    "type": "string"
                },
                "street": {
                    "type": "string"
                }
            }
        },
        "service.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "Код ошибки",
                    "type": "integer"
                },
                "message": {
                    "description": "Сообщение об ошибке",
                    "type": "string"
                }
            }
        },
        "service.GeocodeRequest": {
            "type": "object",
            "properties": {
                "lat": {
                    "type": "number"
                },
                "lng": {
                    "type": "number"
                }
            }
        },
        "service.RequestAddressSearch": {
            "type": "object",
            "properties": {
                "query": {
                    "type": "string"
                }
            }
        },
        "service.ResponseAddress": {
            "type": "object",
            "properties": {
                "suggestions": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "properties": {
                            "data": {
                                "$ref": "#/definitions/service.Address"
                            }
                        }
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Address API",
	Description:      "Этот эндпоинт позволяет получить адрес по наименованию",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
