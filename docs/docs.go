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
        "/auth": {
            "post": {
                "description": "Returns a new access and refresh token pair",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create new token pair",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/auth.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/auth.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/auth.Response"
                        }
                    }
                }
            }
        },
        "/refresh": {
            "post": {
                "description": "Returns a new access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Refresh access token",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/refresh.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/refresh.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/refresh.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/refresh.Response"
                        }
                    }
                }
            }
        },
        "/users": {
            "post": {
                "description": "Returns a new user GUID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "Request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/createUser.Request"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/createUser.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/createUser.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/createUser.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.Request": {
            "type": "object",
            "required": [
                "GUID"
            ],
            "properties": {
                "GUID": {
                    "type": "string"
                }
            }
        },
        "auth.Response": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "createUser.Request": {
            "type": "object",
            "required": [
                "email",
                "name"
            ],
            "properties": {
                "access_token_lifetime_minutes": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "max_active_token_pairs": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "refresh_token_lifetime_minutes": {
                    "type": "integer"
                }
            }
        },
        "createUser.Response": {
            "type": "object",
            "properties": {
                "GUID": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                }
            }
        },
        "refresh.Request": {
            "type": "object",
            "required": [
                "access_token",
                "refresh_token"
            ],
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "refresh.Response": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "authenticationService API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
