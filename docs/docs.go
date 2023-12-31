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
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ssv/cluster/deposit": {
            "post": {
                "description": "Deposit cluster ssv",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cluster"
                ],
                "summary": "cluster Deposit ssv",
                "parameters": [
                    {
                        "description": "amount",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.clusterAmountRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Task uuid",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/cluster/nonce": {
            "get": {
                "description": "get cluster owner register nonce",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "nonce"
                ],
                "summary": "cluster owner register nonce",
                "responses": {
                    "200": {
                        "description": "cluster register nonce",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "set cluster owner register nonce",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "nonce"
                ],
                "summary": "cluster owner register nonce",
                "parameters": [
                    {
                        "description": "nonce",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.clusterNoncePutRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "massage",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/cluster/reactive": {
            "post": {
                "description": "reactive cluster",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cluster"
                ],
                "summary": "clusterReactive",
                "responses": {
                    "200": {
                        "description": "Task uuid(string)",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/cluster/withdraw": {
            "post": {
                "description": "Withdraw cluster ssv",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "cluster"
                ],
                "summary": "cluster Withdraw ssv",
                "parameters": [
                    {
                        "description": "amount",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.clusterAmountRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Task uuid",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/feeRecipient/{TaskId}/state": {
            "get": {
                "description": "check FeeRecipient ClusterValidatorTask State",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "FeeRecipient"
                ],
                "summary": "feeRecipientState",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "fee recipient task id",
                        "name": "TaskId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/registerValidator": {
            "post": {
                "description": "Provide the publicKey array to register Validator in ssv network",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Validator"
                ],
                "summary": "registerValidator",
                "parameters": [
                    {
                        "description": "publicKey array",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.validatorBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "tasks create",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/removeValidator": {
            "post": {
                "description": "remove Validator to ssv network",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Validator"
                ],
                "summary": "removeValidator",
                "parameters": [
                    {
                        "description": "publicKey array",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.validatorBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "tasks create",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/setFeeRecipientAddress": {
            "post": {
                "description": "set FeeRecipient Address on ssv network",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "FeeRecipient"
                ],
                "summary": "setFeeRecipientAddress",
                "parameters": [
                    {
                        "description": "New Fee Recipient Address",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.setFeeRecipientAddressBody"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "FeeRecipientTask Id:(number)",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/upload": {
            "post": {
                "description": "upload keystore by dvt tools",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "upload"
                ],
                "summary": "upload",
                "parameters": [
                    {
                        "description": "keystore info",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UploadKeystore"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "params error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "server error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        },
        "/ssv/{Validator}/state/": {
            "get": {
                "description": "check Validator state",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Validator"
                ],
                "summary": "validatorState",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Validator publicKey",
                        "name": "Validator",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "success",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "400": {
                        "description": "validator error",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    },
                    "500": {
                        "description": "other fail",
                        "schema": {
                            "$ref": "#/definitions/service.resultResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.KeystoreV4": {
            "type": "object",
            "properties": {
                "crypto": {
                    "type": "object",
                    "additionalProperties": true
                },
                "name": {
                    "type": "string"
                },
                "path": {
                    "type": "string"
                },
                "pubkey": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                },
                "version": {
                    "type": "integer"
                }
            }
        },
        "models.UploadKeystore": {
            "type": "object",
            "properties": {
                "keys": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.KeystoreV4"
                    }
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "service.clusterAmountRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                }
            }
        },
        "service.clusterNoncePutRequest": {
            "type": "object",
            "properties": {
                "nonce": {
                    "type": "integer"
                }
            }
        },
        "service.resultResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "response code",
                    "type": "integer"
                },
                "data": {
                    "description": "other data info"
                },
                "msg": {
                    "description": "error or success msg",
                    "type": "string"
                }
            }
        },
        "service.setFeeRecipientAddressBody": {
            "type": "object",
            "properties": {
                "fee_recipient": {
                    "description": "New Fee Recipient Address",
                    "type": "string"
                }
            }
        },
        "service.validatorBody": {
            "type": "object",
            "properties": {
                "public_keys": {
                    "description": "public key array",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "DVT Service API",
	Description:      "This is a sample server celler server.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
