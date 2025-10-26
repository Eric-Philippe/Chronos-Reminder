package docs

// SwaggerInfo holds exported Swagger Info so clients can modify it
type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo variable
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/api",
	Schemes:     []string{"http", "https"},
	Title:       "Chronos Reminder API",
	Description: "API for Chronos Reminder application",
}

// ReadDoc returns the Swagger documentation in JSON format
func ReadDoc() string {
	return `{
  "swagger": "2.0",
  "info": {
    "description": "Chronos Reminder API Documentation",
    "version": "1.0.0",
    "title": "Chronos Reminder API",
    "contact": {
      "name": "API Support"
    }
  },
  "host": "localhost:8080",
  "basePath": "/api",
  "schemes": ["http", "https"],
  "paths": {
    "/auth/register": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Register a new user",
        "description": "Register a new user with email, username, password and timezone",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Registration request",
            "required": true,
            "schema": {
              "$ref": "#/definitions/RegisterRequest"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "User registered successfully",
            "schema": {
              "$ref": "#/definitions/RegisterResponse"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "409": {
            "description": "Email already exists",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/login": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Login user",
        "description": "Login with email and password",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Login request",
            "required": true,
            "schema": {
              "$ref": "#/definitions/LoginRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Login successful",
            "schema": {
              "$ref": "#/definitions/LoginResponse"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "401": {
            "description": "Invalid credentials",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/auth/logout": {
      "post": {
        "tags": ["Authentication"],
        "summary": "Logout user",
        "description": "Logout the current user and invalidate session",
        "produces": ["application/json"],
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "responses": {
          "200": {
            "description": "Logout successful"
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "401": {
            "description": "Unauthorized",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "RegisterRequest": {
      "type": "object",
      "required": ["email", "username", "password", "timezone"],
      "properties": {
        "email": {
          "type": "string",
          "example": "user@example.com"
        },
        "username": {
          "type": "string",
          "maxLength": 128,
          "example": "john_doe"
        },
        "password": {
          "type": "string",
          "minLength": 8,
          "example": "SecurePassword123"
        },
        "timezone": {
          "type": "string",
          "example": "America/New_York"
        }
      }
    },
    "RegisterResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "123e4567-e89b-12d3-a456-426614174000"
        },
        "email": {
          "type": "string",
          "example": "user@example.com"
        },
        "username": {
          "type": "string",
          "example": "john_doe"
        },
        "message": {
          "type": "string",
          "example": "User registered successfully"
        }
      }
    },
    "LoginRequest": {
      "type": "object",
      "required": ["email", "password"],
      "properties": {
        "email": {
          "type": "string",
          "example": "user@example.com"
        },
        "password": {
          "type": "string",
          "example": "SecurePassword123"
        },
        "remember_me": {
          "type": "boolean",
          "example": true
        }
      }
    },
    "LoginResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "123e4567-e89b-12d3-a456-426614174000"
        },
        "email": {
          "type": "string",
          "example": "user@example.com"
        },
        "username": {
          "type": "string",
          "example": "john_doe"
        },
        "token": {
          "type": "string",
          "example": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
        },
        "expires_at": {
          "type": "string",
          "example": "2024-10-27T15:04:05Z"
        },
        "message": {
          "type": "string",
          "example": "Login successful"
        }
      }
    },
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "error": {
          "type": "string",
          "example": "Invalid request"
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
}
