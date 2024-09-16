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
        "/admin/admins": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "管理者ユーザーを一覧で取得します",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "管理者ユーザー一覧取得",
                "responses": {}
            }
        },
        "/admin/admins/i": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "ログイン中の管理者情報取得します",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "ログイン中の管理者情報取得",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Admin ID",
                        "name": "adminId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/admin/customers": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "全顧客を一覧で取得します",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "顧客一覧取得",
                "responses": {}
            }
        },
        "/admin/customers/{customerId}": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "顧客を一件取得します",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "顧客取得",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Customer ID",
                        "name": "customerId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/admin/customers/{customerId}/posts": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "顧客ごとの投稿データを一覧で取得します",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "投稿取得",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Customer ID",
                        "name": "customerId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/admin/login": {
            "post": {
                "description": "管理者としてログインします",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "ログイン",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email",
                        "name": "email",
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
                "responses": {}
            }
        },
        "/admin/register/admin": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "管理者ユーザーを作成します",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "管理者ユーザーの作成",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Email",
                        "name": "email",
                        "in": "formData"
                    }
                ],
                "responses": {}
            }
        },
        "/admin/register/customer": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "顧客を作成します",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "summary": "顧客の作成",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Email",
                        "name": "email",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "WordPress URL",
                        "name": "wordpress_url",
                        "in": "formData"
                    }
                ],
                "responses": {}
            }
        },
        "/badge/execute": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "全顧客に対して、InstagramからWordpressへのデータ連携を行います。",
                "tags": [
                    "Badge"
                ],
                "summary": "バッチ実行",
                "responses": {}
            }
        },
        "/customer/i": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "自分の顧客情報を取得します",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customer"
                ],
                "summary": "顧客情報の取得",
                "responses": {}
            }
        },
        "/customer/i/fetch/post": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "インスタグラムとWordpressの連携",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customer"
                ],
                "summary": "インスタグラムとWordpressの連携",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Customer ID",
                        "name": "customer_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/customer/i/posts": {
            "get": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "顧客ごとの投稿を一覧で取得します",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Customer"
                ],
                "summary": "顧客の投稿一覧の取得",
                "responses": {}
            }
        },
        "/customer/login": {
            "post": {
                "security": [
                    {
                        "Token": []
                    }
                ],
                "description": "顧客としてログインします",
                "consumes": [
                    "application/x-www-form-urlencoded"
                ],
                "tags": [
                    "Customer"
                ],
                "summary": "ログイン",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Email",
                        "name": "email",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Password",
                        "name": "password",
                        "in": "formData"
                    }
                ],
                "responses": {}
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
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
