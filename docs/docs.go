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
        "contact": {
            "name": "Cheqd",
            "url": "https://cheqd.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/1.0/identifiers/{did}": {
            "get": {
                "description": "Get DID Doc or its fragment",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "*/*"
                ],
                "tags": [
                    "Resolution"
                ],
                "summary": "Resolve or dereferencing DID Doc",
                "parameters": [
                    {
                        "type": "string",
                        "example": "did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY",
                        "description": "DID Doc Id",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "application/did+ld+json",
                            "application/ld+json",
                            "application/did+json"
                        ],
                        "type": "string",
                        "description": "The requested media type of the DID document representation or DID resolution result. ",
                        "name": "accept",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.DidResolution"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.DidResolution"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.DidResolution"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.DidResolution"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.DidResolution"
                        }
                    }
                }
            }
        },
        "/1.0/identifiers/{did}/resources/all": {
            "get": {
                "description": "Get a list of all collection resources metadata",
                "produces": [
                    "*/*"
                ],
                "tags": [
                    "Dereferencing"
                ],
                "summary": "Collection resources",
                "parameters": [
                    {
                        "type": "string",
                        "example": "did:cheqd:testnet:MjYxNzYKMjYxNzYK",
                        "description": "Resource collection id. DID Doc Id",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "application/did+ld+json",
                            "application/ld+json",
                            "application/did+json"
                        ],
                        "type": "string",
                        "description": "The requested media type of the DID document representation or DID resolution result. ",
                        "name": "accept",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    }
                }
            }
        },
        "/1.0/identifiers/{did}/resources/{resourceId}": {
            "get": {
                "description": "Get resource value without dereferencing wrappers",
                "produces": [
                    "*/*"
                ],
                "tags": [
                    "Dereferencing"
                ],
                "summary": "Resource value",
                "parameters": [
                    {
                        "type": "string",
                        "example": "did:cheqd:testnet:MjYxNzYKMjYxNzYK",
                        "description": "Resource collection id. DID Doc Id",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "60ad67be-b65b-40b8-b2f4-3923141ef382",
                        "description": "DID Resource identifier",
                        "name": "resourceId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "application/did+ld+json",
                            "application/ld+json",
                            "application/did+json"
                        ],
                        "type": "string",
                        "description": "The requested media type of the DID document representation or DID resolution result. ",
                        "name": "accept",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "integer"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    }
                }
            }
        },
        "/1.0/identifiers/{did}/resources/{resourceId}/metadata": {
            "get": {
                "description": "Get resource metadata without value by DID Doc",
                "produces": [
                    "*/*"
                ],
                "tags": [
                    "Dereferencing"
                ],
                "summary": "Resource metadata",
                "parameters": [
                    {
                        "type": "string",
                        "example": "did:cheqd:testnet:MjYxNzYKMjYxNzYK",
                        "description": "Resource collection id. DID Doc Id",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "60ad67be-b65b-40b8-b2f4-3923141ef382",
                        "description": "DID Resource identifier",
                        "name": "resourceId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "enum": [
                            "application/did+ld+json",
                            "application/ld+json",
                            "application/did+json"
                        ],
                        "type": "string",
                        "description": "The requested media type of the DID document representation or DID resolution result. ",
                        "name": "accept",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    }
                }
            }
        },
        "/1.0/identifiers/{did}{fragmentId}": {
            "get": {
                "description": "Get DID Doc or its fragment",
                "consumes": [
                    "application/did+ld+json",
                    "application/ld+json",
                    "application/did+json"
                ],
                "produces": [
                    "application/did+ld+json",
                    "application/ld+json",
                    "application/did+json"
                ],
                "tags": [
                    "Dereferencing"
                ],
                "summary": "Resolve or dereferencing DID Doc",
                "parameters": [
                    {
                        "type": "string",
                        "example": "did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY",
                        "description": "DID Doc Id",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "#key1",
                        "description": "` + "`" + `#` + "`" + ` + DID Doc Verification Method or Service identifier",
                        "name": "fragmentId",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "example": "\"service1\"",
                        "description": "Service id",
                        "name": "service",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "application/did+ld+json",
                            "application/ld+json",
                            "application/did+json"
                        ],
                        "type": "string",
                        "description": "The requested media type of the DID document representation or DID resolution result. ",
                        "name": "accept",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.DidDereferencing"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.DereferencedResource": {
            "type": "object",
            "properties": {
                "checksum": {
                    "type": "string"
                },
                "created": {
                    "type": "string"
                },
                "mediaType": {
                    "type": "string"
                },
                "nextVersionId": {
                    "type": "string"
                },
                "previousVersionId": {
                    "type": "string"
                },
                "resourceCollectionId": {
                    "type": "string"
                },
                "resourceId": {
                    "type": "string"
                },
                "resourceName": {
                    "type": "string"
                },
                "resourceType": {
                    "type": "string"
                },
                "resourceURI": {
                    "type": "string"
                }
            }
        },
        "types.DereferencingMetadata": {
            "type": "object",
            "properties": {
                "contentType": {
                    "type": "string"
                },
                "did": {
                    "$ref": "#/definitions/types.DidProperties"
                },
                "error": {
                    "type": "string"
                },
                "retrieved": {
                    "type": "string"
                }
            }
        },
        "types.DidDereferencing": {
            "type": "object",
            "properties": {
                "@context": {
                    "type": "string",
                    "example": "https://w3id.org/did-resolution/v1"
                },
                "contentMetadata": {
                    "$ref": "#/definitions/types.ResolutionDidDocMetadata"
                },
                "contentStream": {},
                "dereferencingMetadata": {
                    "$ref": "#/definitions/types.DereferencingMetadata"
                }
            }
        },
        "types.DidDoc": {
            "type": "object",
            "properties": {
                "@context": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "alsoKnownAs": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "assertionMethod": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "authentication": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "capabilityInvocation": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "capability_delegation": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "controller": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "keyAgreement": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "service": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.Service"
                    }
                },
                "verificationMethod": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.VerificationMethod"
                    }
                }
            }
        },
        "types.DidProperties": {
            "type": "object",
            "properties": {
                "didString": {
                    "type": "string"
                },
                "method": {
                    "type": "string"
                },
                "methodSpecificId": {
                    "type": "string"
                }
            }
        },
        "types.DidResolution": {
            "type": "object",
            "properties": {
                "@context": {
                    "type": "string"
                },
                "didDocument": {
                    "$ref": "#/definitions/types.DidDoc"
                },
                "didDocumentMetadata": {
                    "$ref": "#/definitions/types.ResolutionDidDocMetadata"
                },
                "didResolutionMetadata": {
                    "$ref": "#/definitions/types.ResolutionMetadata"
                }
            }
        },
        "types.ResolutionDidDocMetadata": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "string"
                },
                "deactivated": {
                    "type": "boolean"
                },
                "linkedResourceMetadata": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.DereferencedResource"
                    }
                },
                "updated": {
                    "type": "string"
                },
                "versionId": {
                    "type": "string"
                }
            }
        },
        "types.ResolutionMetadata": {
            "type": "object",
            "properties": {
                "contentType": {
                    "type": "string"
                },
                "did": {
                    "$ref": "#/definitions/types.DidProperties"
                },
                "error": {
                    "type": "string"
                },
                "retrieved": {
                    "type": "string"
                }
            }
        },
        "types.Service": {
            "type": "object",
            "properties": {
                "@context": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "id": {
                    "type": "string"
                },
                "serviceEndpoint": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "types.VerificationMethod": {
            "type": "object",
            "properties": {
                "@context": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "controller": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "publicKeyJwk": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "publicKeyMultibase": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Cheqd DID Resolver API",
	Description:      "Cheqd DID Resolver API for DID resolution and dereferencing.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
