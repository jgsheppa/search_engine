// Package swagger GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package swagger

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/document": {
            "post": {
                "tags": [
                    "Document"
                ],
                "summary": "Post documents to Redisearch",
                "parameters": [
                    {
                        "description": "The body to create a Redis document for an article",
                        "name": "Body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Document"
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Document"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        },
        "/api/document/delete/{documentName}": {
            "delete": {
                "tags": [
                    "Document"
                ],
                "summary": "Delete documents from Redisearch",
                "operationId": "documentName",
                "parameters": [
                    {
                        "type": "string",
                        "description": "search term",
                        "name": "documentName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        },
        "/api/document/{url}/{htmlTag}/{containerClass}": {
            "post": {
                "tags": [
                    "Document"
                ],
                "summary": "Post documents to Redisearch",
                "parameters": [
                    {
                        "type": "string",
                        "description": "URL path of guidebook",
                        "name": "url",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "HTML tag to scrape",
                        "name": "htmlTag",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Class of container holding HTML tag",
                        "name": "containerClass",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Document"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        },
        "/api/index/create": {
            "post": {
                "tags": [
                    "Index"
                ],
                "summary": "Create Redis index for BestPracticer guides",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/index/delete": {
            "delete": {
                "tags": [
                    "Index"
                ],
                "summary": "Delete all documents from Redisearch",
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/search/geo": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "GeoSearch"
                ],
                "summary": "Search Redisearch documents",
                "operationId": "article geo-search",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search by keyword",
                        "name": "term",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        },
        "/api/search/{term}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Search"
                ],
                "summary": "Search Redisearch documents",
                "operationId": "article search",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Search by keyword",
                        "name": "term",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Sort by field",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "boolean",
                        "description": "Ascending?",
                        "name": "ascending",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Limit number of results",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Ok",
                        "schema": {
                            "$ref": "#/definitions/swagger.SwaggerSearchResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/models.ApiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.ApiError": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Document": {
            "type": "object",
            "properties": {
                "country": {
                    "description": "URL of article",
                    "type": "string"
                },
                "lat": {
                    "description": "Latitude of city",
                    "type": "string"
                },
                "lng": {
                    "description": "Longitude of city",
                    "type": "string"
                },
                "name": {
                    "description": "Title of article",
                    "type": "string"
                }
            }
        },
        "swagger.SwaggerResponse": {
            "type": "object",
            "properties": {
                "country": {
                    "description": "Author of article",
                    "type": "string",
                    "example": "Austria"
                },
                "location": {
                    "description": "Topics of article",
                    "type": "string",
                    "example": "48.20849,16.37208"
                },
                "name": {
                    "description": "URL of doc",
                    "type": "string",
                    "example": "Wien"
                }
            }
        },
        "swagger.SwaggerSearchResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/swagger.SwaggerResponse"
                    }
                },
                "suggestion": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/swagger.SwaggerSuggestion"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "swagger.SwaggerSuggestion": {
            "type": "object",
            "properties": {
                "incr": {
                    "type": "boolean",
                    "example": false
                },
                "payload": {
                    "type": "string",
                    "example": "Pair"
                },
                "score": {
                    "type": "number",
                    "example": 70.7
                },
                "term": {
                    "type": "string",
                    "example": "Pair"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3001",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "Redisearch Search Engine",
	Description:      "This is a search engine built with Redisearch",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
