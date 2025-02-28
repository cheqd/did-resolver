// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Cheqd Foundation Limited",
            "url": "https://cheqd.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "https://github.com/cheqd/did-resolver/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/{did}": {
            "get": {
                "description": "Fetch DID Document (\"DIDDoc\") from cheqd network",
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
                    "DID Resolution"
                ],
                "summary": "Resolve DID Document on did:cheqd",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Full DID with unique identifier",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "#Fragment",
                        "name": "fragmentId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Version",
                        "name": "versionId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Created of Updated time of DID Document",
                        "name": "versionTime",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Can transform Verification Method into another type",
                        "name": "transformKeys",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Redirects to Service Endpoint",
                        "name": "service",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Addition to Service Endpoint",
                        "name": "relativeRef",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Show only metadata of DID Document",
                        "name": "metadata",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by ResourceId",
                        "name": "resourceId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by CollectionId",
                        "name": "resourceCollectionId",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Resource Type",
                        "name": "resourceType",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Resource Name",
                        "name": "resourceName",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Filter by Resource Version",
                        "name": "resourceVersion",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Get the nearest resource by creation time",
                        "name": "resourceVersionTime",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Show only metadata of resources",
                        "name": "resourceMetadata",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Sanity check that Checksum of resource is the same as expected",
                        "name": "checksum",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "versionId, versionTime, transformKeys returns Full DID Document",
                        "schema": {
                            "$ref": "#/definitions/types.DidResolution"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    }
                }
            }
        },
        "/{did}/metadata": {
            "get": {
                "description": "Get metadata for all Resources within a DID Resource Collection",
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
                    "Resource Resolution"
                ],
                "summary": "Fetch metadata for all Resources",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Full DID with unique identifier",
                        "name": "did",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.ResourceDereferencing"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "contentStream": {
                                            "$ref": "#/definitions/types.ResolutionDidDocMetadata"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    }
                }
            }
        },
        "/{did}/resources/{resourceId}": {
            "get": {
                "description": "Get specific Resource within a DID Resource Collection",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "*/*"
                ],
                "tags": [
                    "Resource Resolution"
                ],
                "summary": "Fetch specific Resource",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Full DID with unique identifier",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Resource-specific unique-identifier",
                        "name": "resourceId",
                        "in": "path",
                        "required": true
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
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    }
                }
            }
        },
        "/{did}/resources/{resourceId}/metadata": {
            "get": {
                "description": "Get metadata for a specific Resource within a DID Resource Collection",
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
                    "Resource Resolution"
                ],
                "summary": "Fetch Resource-specific metadata",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Full DID with unique identifier",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Resource-specific unique identifier",
                        "name": "resourceId",
                        "in": "path",
                        "required": true
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
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    }
                }
            }
        },
        "/{did}/version/{versionId}": {
            "get": {
                "description": "Fetch specific all version of a DID Document (\"DIDDoc\") for a given DID and version ID",
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
                    "DID Resolution"
                ],
                "summary": "Resolve DID Document Version on did:cheqd",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Full DID with unique identifier",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version of a DID document",
                        "name": "versionId",
                        "in": "path",
                        "required": true
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
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    }
                }
            }
        },
        "/{did}/version/{versionId}/metadata": {
            "get": {
                "description": "Fetch metadata of specific a DID Document (\"DIDDoc\") version for a given DID and version ID",
                "consumes": [
                    "application/did+ld+json",
                    "application/ld+json",
                    "application/did+jsonww"
                ],
                "produces": [
                    "application/did+ld+json",
                    "application/ld+json",
                    "application/did+json"
                ],
                "tags": [
                    "DID Resolution"
                ],
                "summary": "Resolve DID Document Version Metadata on did:cheqd",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Full DID with unique identifier",
                        "name": "did",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "version of a DID document",
                        "name": "versionId",
                        "in": "path",
                        "required": true
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
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    }
                }
            }
        },
        "/{did}/versions": {
            "get": {
                "description": "Fetch specific all versions of a DID Document (\"DIDDoc\") for a given DID",
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
                    "DID Resolution"
                ],
                "summary": "Resolve DID Document Versions on did:cheqd",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Full DID with unique identifier",
                        "name": "did",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/types.ResourceDereferencing"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "contentStream": {
                                            "$ref": "#/definitions/types.DereferencedDidVersionsList"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "406": {
                        "description": "Not Acceptable",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "$ref": "#/definitions/types.IdentityError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "types.ContentType": {
            "type": "string",
            "enum": [
                "application/did+json",
                "application/did+ld+json",
                "application/ld+json",
                "application/json",
                "application/did",
                "text/plain"
            ],
            "x-enum-varnames": [
                "DIDJSON",
                "DIDJSONLD",
                "JSONLD",
                "JSON",
                "DIDRES",
                "TEXT"
            ]
        },
        "types.DereferencedDidVersionsList": {
            "type": "object",
            "properties": {
                "versions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.ResolutionDidDocMetadata"
                    }
                }
            }
        },
        "types.DereferencedResource": {
            "type": "object",
            "properties": {
                "checksum": {
                    "type": "string",
                    "example": "a95380f460e63ad939541a57aecbfd795fcd37c6d78ee86c885340e33a91b559"
                },
                "created": {
                    "type": "string",
                    "example": "2021-09-01T12:00:00Z"
                },
                "mediaType": {
                    "type": "string",
                    "example": "image/png"
                },
                "nextVersionId": {
                    "type": "string",
                    "example": "d4829ac7-4566-478c-a408-b44767eddadc"
                },
                "previousVersionId": {
                    "type": "string",
                    "example": "ad7a8442-3531-46eb-a024-53953ec6e4ff"
                },
                "resourceCollectionId": {
                    "type": "string",
                    "example": "55dbc8bf-fba3-4117-855c-1e0dc1d3bb47"
                },
                "resourceId": {
                    "type": "string",
                    "example": "398cee0a-efac-4643-9f4c-74c48c72a14b"
                },
                "resourceName": {
                    "type": "string",
                    "example": "Image Resource"
                },
                "resourceType": {
                    "type": "string",
                    "example": "Image"
                },
                "resourceURI": {
                    "type": "string",
                    "example": "did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47/resources/398cee0a-efac-4643-9f4c-74c48c72a14b"
                },
                "resourceVersion": {
                    "type": "string",
                    "example": "1"
                }
            }
        },
        "types.DereferencingMetadata": {
            "type": "object",
            "properties": {
                "contentType": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/types.ContentType"
                        }
                    ],
                    "example": "application/ld+json"
                },
                "did": {
                    "$ref": "#/definitions/types.DidProperties"
                },
                "error": {
                    "type": "string"
                },
                "retrieved": {
                    "type": "string",
                    "example": "2021-09-01T12:00:00Z"
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
                    },
                    "example": [
                        "https://www.w3.org/ns/did/v1"
                    ]
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
                    },
                    "example": [
                        "did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47#key-1"
                    ]
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
                    },
                    "example": [
                        "did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47"
                    ]
                },
                "id": {
                    "type": "string",
                    "example": "did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47"
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
        "types.IdentityError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "contentType": {
                    "$ref": "#/definitions/types.ContentType"
                },
                "did": {
                    "type": "string"
                },
                "internal": {},
                "isDereferencing": {
                    "type": "boolean"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "types.ResolutionDidDocMetadata": {
            "type": "object",
            "properties": {
                "created": {
                    "type": "string",
                    "example": "2021-09-01T12:00:00Z"
                },
                "deactivated": {
                    "type": "boolean",
                    "example": false
                },
                "linkedResourceMetadata": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.DereferencedResource"
                    }
                },
                "nextVersionId": {
                    "type": "string",
                    "example": "3f3111af-dfe6-411f-adc9-02af59716ddb"
                },
                "previousVersionId": {
                    "type": "string",
                    "example": "139445af-4281-4453-b05a-ec9a8931c1f9"
                },
                "updated": {
                    "type": "string",
                    "example": "2021-09-10T12:00:00Z"
                },
                "versionId": {
                    "type": "string",
                    "example": "284f297b-b6e3-4ffa-9172-bc3bb904e286"
                }
            }
        },
        "types.ResolutionMetadata": {
            "type": "object",
            "properties": {
                "contentType": {
                    "allOf": [
                        {
                            "$ref": "#/definitions/types.ContentType"
                        }
                    ],
                    "example": "application/ld+json"
                },
                "did": {
                    "$ref": "#/definitions/types.DidProperties"
                },
                "error": {
                    "type": "string"
                },
                "retrieved": {
                    "type": "string",
                    "example": "2021-09-01T12:00:00Z"
                }
            }
        },
        "types.ResolutionResourceMetadata": {
            "type": "object",
            "properties": {
                "linkedResourceMetadata": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.DereferencedResource"
                    }
                },
                "metadata": {
                    "$ref": "#/definitions/types.DereferencedResource"
                }
            }
        },
        "types.ResourceDereferencing": {
            "type": "object",
            "properties": {
                "@context": {
                    "type": "string",
                    "example": "https://w3id.org/did-resolution/v1"
                },
                "contentMetadata": {
                    "$ref": "#/definitions/types.ResolutionResourceMetadata"
                },
                "contentStream": {},
                "dereferencingMetadata": {
                    "$ref": "#/definitions/types.DereferencingMetadata"
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
                    "type": "string",
                    "example": "did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47#service-1"
                },
                "serviceEndpoint": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "https://example.com/endpoint/8377464"
                    ]
                },
                "type": {
                    "type": "string",
                    "example": "did-communication"
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
                "publicKeyBase58": {
                    "type": "string"
                },
                "publicKeyJwk": {},
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
	Version:          "v3.0",
	Host:             "resolver.cheqd.net",
	BasePath:         "/1.0/identifiers",
	Schemes:          []string{"https", "http"},
	Title:            "DID Resolver for cheqd DID method",
	Description:      "Universal Resolver driver for cheqd DID method",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
