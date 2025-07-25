// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "http://www.swagger.io/support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "MIT",
            "url": "https://opensource.org/licenses/MIT"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/trip/preview": {
            "post": {
                "description": "Get trip preview with estimated fares for different ride packages",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "trips"
                ],
                "summary": "Preview a trip",
                "parameters": [
                    {
                        "description": "Trip preview request",
                        "name": "preview",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.previewTripRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Trip preview with estimated fares",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/trip/start": {
            "post": {
                "description": "Creates a new trip with the provided ride fare ID and user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "trips"
                ],
                "summary": "Start a trip",
                "parameters": [
                    {
                        "description": "Start trip request",
                        "name": "trip",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.startTripRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Trip started successfully",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.previewTripRequest": {
            "type": "object",
            "required": [
                "destination",
                "pickup",
                "userID"
            ],
            "properties": {
                "destination": {
                    "$ref": "#/definitions/types.Coordinate"
                },
                "pickup": {
                    "$ref": "#/definitions/types.Coordinate"
                },
                "userID": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "main.startTripRequest": {
            "type": "object",
            "required": [
                "rideFareID",
                "userID"
            ],
            "properties": {
                "rideFareID": {
                    "type": "string",
                    "example": "fare123"
                },
                "userID": {
                    "type": "string",
                    "example": "user123"
                }
            }
        },
        "types.Coordinate": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8081",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Bolt App API Gateway",
	Description:      "This is the API Gateway for the Bolt ride-sharing microservices application.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
