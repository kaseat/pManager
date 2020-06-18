// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2020-06-18 19:22:29.212856766 +0300 MSK m=+0.088258398

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://github.com/kaseat/pManager/blob/master/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/portfolios": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Gets all portfolios avaliable",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "portfolios"
                ],
                "summary": "Get all portfolios",
                "operationId": "portfolio-get-all",
                "responses": {
                    "200": {
                        "description": "Returns portfolio info",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Portfolio"
                            }
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Creates single portfolio",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "portfolios"
                ],
                "summary": "Add new portfolio",
                "operationId": "portfolio-add",
                "parameters": [
                    {
                        "description": "Portfolio info",
                        "name": "portfolio",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.portfolioRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns portfolio Id just created",
                        "schema": {
                            "$ref": "#/definitions/api.addPortfoliioSuccess"
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Deletes all portfolios an all associated operations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "portfolios"
                ],
                "summary": "Delete all portfolios",
                "operationId": "portfolio-del-all",
                "responses": {
                    "200": {
                        "description": "Returns true if portfolios has deleted",
                        "schema": {
                            "$ref": "#/definitions/api.delPortfoliioSuccess"
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            }
        },
        "/portfolios/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Gets portfolio info by Id",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "portfolios"
                ],
                "summary": "Get portfolio by Id",
                "operationId": "portfolio-get-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Portfolio Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns portfolio info if any",
                        "schema": {
                            "$ref": "#/definitions/models.Portfolio"
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Updates portfolio info by Id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "portfolios"
                ],
                "summary": "Update portfolio info",
                "operationId": "portfolio-put-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Portfolio Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Portfolio info",
                        "name": "portfolio",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.portfolioRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns portfolio info if any",
                        "schema": {
                            "$ref": "#/definitions/api.putPortfoliioSuccess"
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Deletes portfolio an all associated operations",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "portfolios"
                ],
                "summary": "Delete portfolio",
                "operationId": "portfolio-del-by-id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Portfolio Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns true if portfolio has deleted",
                        "schema": {
                            "$ref": "#/definitions/api.delPortfoliioSuccess"
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            }
        },
        "/portfolios/{id}/operations": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Gets all operations for specified portfolio",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "operations",
                    "portfolios"
                ],
                "summary": "Get all operations",
                "operationId": "operation-get-all",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Portfolio Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Filter by FIGI",
                        "name": "figi",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter operations from this date",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter operations till this date",
                        "name": "to",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns operations info",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/portfolio.Operation"
                            }
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Adds operation to specified portfolio",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "operations",
                    "portfolios"
                ],
                "summary": "Add new operation",
                "operationId": "operation-add",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Portfolio Id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Operation info",
                        "name": "portfolio",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.operationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns portfolio Id just created",
                        "schema": {
                            "$ref": "#/definitions/api.addOperationSuccess"
                        }
                    },
                    "400": {
                        "description": "Returns when any processing error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    },
                    "401": {
                        "description": "Returns when authentication error occurs",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            }
        },
        "/token/validate": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get string by ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "security"
                ],
                "summary": "Validate token",
                "operationId": "validate-token",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.tokenResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "description": "Checks user credentials and returns JWT if ok",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Login",
                "operationId": "login",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User name",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.tokenResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            }
        },
        "/user/signup": {
            "post": {
                "description": "Creates new user",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Create new user",
                "operationId": "sign-up",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User name",
                        "name": "username",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.tokenResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/api.errorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.addOperationSuccess": {
            "type": "object",
            "properties": {
                "createdOperationId": {
                    "type": "string",
                    "example": "5edbc0a72c857652a0542fab"
                }
            }
        },
        "api.addPortfoliioSuccess": {
            "type": "object",
            "properties": {
                "createdPortfolioId": {
                    "type": "string",
                    "example": "5edb2a0e550dfc5f16392838"
                }
            }
        },
        "api.delPortfoliioSuccess": {
            "type": "object",
            "properties": {
                "hasDeleted": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "api.errorResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "Something went wrong"
                }
            }
        },
        "api.operationRequest": {
            "type": "object",
            "properties": {
                "currency": {
                    "type": "string",
                    "example": "USD"
                },
                "date": {
                    "type": "string",
                    "example": "2020-06-06T22:54:05.000+07:00"
                },
                "figi": {
                    "type": "string",
                    "example": "BBG00MVRXDB0"
                },
                "operationType": {
                    "type": "string",
                    "example": "sell"
                },
                "price": {
                    "type": "number",
                    "example": 293.61
                },
                "vol": {
                    "type": "integer",
                    "example": 100
                }
            }
        },
        "api.portfolioRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Best portfolio ever!!!"
                },
                "name": {
                    "type": "string",
                    "example": "Best portfolio"
                }
            }
        },
        "api.putPortfoliioSuccess": {
            "type": "object",
            "properties": {
                "hasModified": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "api.tokenResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "models.Portfolio": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string",
                    "example": "Best portfolio ever!!!"
                },
                "id": {
                    "type": "string",
                    "example": "5edb2a0e550dfc5f16392838"
                },
                "name": {
                    "type": "string",
                    "example": "Best portfolio"
                }
            }
        },
        "portfolio.Operation": {
            "type": "object",
            "properties": {
                "currency": {
                    "type": "string",
                    "example": "USD"
                },
                "date": {
                    "type": "string",
                    "example": "2020-06-06T15:54:05Z"
                },
                "figi": {
                    "type": "string",
                    "example": "BBG00MVRXDB0"
                },
                "id": {
                    "type": "string",
                    "example": "5edbc0a72c857652a0542fab"
                },
                "operationType": {
                    "type": "string",
                    "example": "sell"
                },
                "pid": {
                    "type": "string"
                },
                "price": {
                    "type": "number",
                    "example": 293.61
                },
                "vol": {
                    "type": "integer",
                    "example": 100
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Portfolio manager API",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
