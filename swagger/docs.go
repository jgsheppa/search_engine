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
        "/api/document/article": {
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
                                "$ref": "#/definitions/models.Article"
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
                                "$ref": "#/definitions/models.Article"
                            }
                        }
                    },
                    "422": {
                        "description": ""
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
                    }
                }
            }
        },
        "/api/document/guide/{url}": {
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
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Guide"
                            }
                        }
                    },
                    "422": {
                        "description": ""
                    }
                }
            }
        },
        "/api/index/create/": {
            "post": {
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
        "/api/index/create/articles": {
            "post": {
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
        "/api/index/delete/": {
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
        "/api/search/{term}/{service}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Search"
                ],
                "summary": "Search Redisearch documents",
                "operationId": "term",
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
                        "description": "Options: guide or article",
                        "name": "service",
                        "in": "path"
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
                            "$ref": "#/definitions/controllers.SwaggerSearchResponse"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.ApiError"
                        }
                    },
                    "500": {
                        "description": "Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.ApiError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.ApiError": {
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
        "controllers.SwaggerSearchResponse": {
            "type": "object",
            "properties": {
                "response": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.swaggerResponse"
                    }
                },
                "suggestion": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/controllers.swaggerSuggestion"
                    }
                },
                "total": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "controllers.swaggerResponse": {
            "type": "object",
            "properties": {
                "author": {
                    "description": "Author of article",
                    "type": "string",
                    "example": "Alex Appleton"
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "title": {
                    "description": "Title of article",
                    "type": "string",
                    "example": "How to be awesome"
                },
                "topic": {
                    "description": "Topics of article",
                    "type": "string",
                    "example": "Awesome Stuff"
                },
                "url": {
                    "description": "URL of article",
                    "type": "string",
                    "example": "www.bestpracticer.com/awesome-article"
                }
            }
        },
        "controllers.swaggerSuggestion": {
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
        },
        "models.Article": {
            "type": "object",
            "properties": {
                "author": {
                    "description": "Author of article",
                    "type": "string"
                },
                "document": {
                    "description": "Document name, if possible a UUID",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "description": "Title of article",
                    "type": "string"
                },
                "topic": {
                    "description": "Topics of article",
                    "type": "string"
                },
                "url": {
                    "description": "URL of article",
                    "type": "string"
                }
            }
        },
        "models.Guide": {
            "type": "object",
            "properties": {
                "document": {
                    "description": "Document name, if possible a UUID",
                    "type": "string"
                },
                "title": {
                    "description": "Title of article",
                    "type": "string"
                },
                "topic": {
                    "description": "Topics of article",
                    "type": "string"
                },
                "url": {
                    "description": "URL of article",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{"http", "https"},
	Title:            "BestPracticer Search Engine",
	Description:      "This is a search engine built with Redisearch",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
