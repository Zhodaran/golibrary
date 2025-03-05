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
        "/api/authors": {
            "post": {
                "description": "This endpoint allows you to add a new author to the library.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authors"
                ],
                "summary": "Add a new author to the library",
                "parameters": [
                    {
                        "description": "Author name",
                        "name": "author",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.AuthorRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Author added successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/book": {
            "post": {
                "description": "This endpoint allows you to add a new book to the library.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Books"
                ],
                "summary": "Add a new book to the library",
                "parameters": [
                    {
                        "description": "Book details",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/repository.Book"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Book added successfully",
                        "schema": {
                            "$ref": "#/definitions/repository.Book"
                        }
                    },
                    "400": {
                        "description": "Invalid request",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/book/return/{index}": {
            "delete": {
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
                    "User"
                ],
                "summary": "Get Geo Coordinates by Address",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Book INDEX",
                        "name": "index",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Request body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.TakeBookRequest"
                        }
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
                    "User"
                ],
                "summary": "Get Geo Coordinates by Address",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Book INDEX",
                        "name": "index",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Request body",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.TakeBookRequest"
                        }
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
        "/api/book/{index}": {
            "put": {
                "description": "Этот эндпоинт позволяет обновить информацию о книге по индексу.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Books"
                ],
                "summary": "Обновление информации о книге",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Индекс книги",
                        "name": "index",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Bearer Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Обновленная информация о книге",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/repository.Book"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешное обновление книги",
                        "schema": {
                            "$ref": "#/definitions/repository.Book"
                        }
                    },
                    "400": {
                        "description": "Ошибка запроса",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Книга не найдена",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/books": {
            "get": {
                "description": "This description created new SQL user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Books"
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
        "/api/get-authors": {
            "get": {
                "description": "Get a list of all authors in the library",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authors"
                ],
                "summary": "Get all authors",
                "responses": {
                    "200": {
                        "description": "List of authors",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "No authors found",
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
                "description": "This endpoint returns a list of all registered users.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Get List of Registered Users",
                "responses": {
                    "200": {
                        "description": "List of registered users",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/auth.User"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/main.mErrorResponse"
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
        "main.AuthorRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "main.TakeBookRequest": {
            "type": "object",
            "properties": {
                "username": {
                    "description": "Поле для имени пользователя",
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
        "repository.Book": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "block": {
                    "type": "boolean"
                },
                "book": {
                    "type": "string"
                },
                "index": {
                    "type": "integer"
                },
                "take_count": {
                    "type": "integer"
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
