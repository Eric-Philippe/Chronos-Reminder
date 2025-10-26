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
    },
    "/auth/discord/callback": {
      "post": {
        "tags": ["Discord OAuth"],
        "summary": "Discord OAuth callback",
        "description": "Complete Discord OAuth authentication flow",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "OAuth callback request",
            "required": true,
            "schema": {
              "$ref": "#/definitions/OAuthCallbackRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Discord OAuth successful",
            "schema": {
              "$ref": "#/definitions/OAuthCallbackResponse"
            }
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
    },
    "/auth/discord/setup": {
      "post": {
        "tags": ["Discord OAuth"],
        "summary": "Complete Discord OAuth setup",
        "description": "Complete app identity setup for Discord-only users",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Setup request",
            "required": true,
            "schema": {
              "$ref": "#/definitions/CompleteDiscordSetupRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Setup completed successfully",
            "schema": {
              "$ref": "#/definitions/OAuthCallbackResponse"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/discord/guilds": {
      "post": {
        "tags": ["Discord"],
        "summary": "Get user guilds",
        "description": "Get list of Discord guilds for the authenticated user",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Request to get user guilds",
            "required": true,
            "schema": {
              "$ref": "#/definitions/GetUserGuildsRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully retrieved guilds",
            "schema": {
              "$ref": "#/definitions/GetUserGuildsResponse"
            }
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
          },
          "404": {
            "description": "Account not found",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/discord/guilds/channels": {
      "post": {
        "tags": ["Discord"],
        "summary": "Get guild channels",
        "description": "Get list of channels for a Discord guild",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Request to get guild channels",
            "required": true,
            "schema": {
              "$ref": "#/definitions/GetGuildChannelsRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully retrieved channels",
            "schema": {
              "$ref": "#/definitions/GetGuildChannelsResponse"
            }
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
          },
          "404": {
            "description": "Account not found",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/discord/guilds/roles": {
      "post": {
        "tags": ["Discord"],
        "summary": "Get guild roles",
        "description": "Get list of roles for a Discord guild",
        "consumes": ["application/json"],
        "produces": ["application/json"],
        "security": [
          {
            "BearerAuth": []
          }
        ],
        "parameters": [
          {
            "in": "body",
            "name": "body",
            "description": "Request to get guild roles",
            "required": true,
            "schema": {
              "$ref": "#/definitions/GetGuildRolesRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Successfully retrieved roles",
            "schema": {
              "$ref": "#/definitions/GetGuildRolesResponse"
            }
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
          },
          "404": {
            "description": "Account not found",
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
    "OAuthCallbackRequest": {
      "type": "object",
      "required": ["code"],
      "properties": {
        "code": {
          "type": "string",
          "description": "Authorization code from Discord",
          "example": "authorization_code_here"
        },
        "state": {
          "type": "string",
          "description": "State parameter for OAuth security",
          "example": "state_parameter"
        }
      }
    },
    "OAuthCallbackResponse": {
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
          "example": "Authentication successful"
        }
      }
    },
    "OAuthSetupRequiredResponse": {
      "type": "object",
      "properties": {
        "status": {
          "type": "string",
          "example": "setup_required"
        },
        "message": {
          "type": "string",
          "example": "Please set up your app identity to complete registration"
        },
        "account_id": {
          "type": "string",
          "example": "123e4567-e89b-12d3-a456-426614174000"
        },
        "discord_email": {
          "type": "string",
          "example": "user@discord.com"
        },
        "discord_username": {
          "type": "string",
          "example": "discord_user"
        },
        "needs_setup": {
          "type": "boolean",
          "example": true
        }
      }
    },
    "CompleteDiscordSetupRequest": {
      "type": "object",
      "required": ["account_id", "email", "password", "timezone"],
      "properties": {
        "account_id": {
          "type": "string",
          "example": "123e4567-e89b-12d3-a456-426614174000"
        },
        "email": {
          "type": "string",
          "example": "user@example.com"
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
    "GetUserGuildsRequest": {
      "type": "object",
      "required": ["account_id"],
      "properties": {
        "account_id": {
          "type": "string",
          "example": "9e6a311e-b86b-40dd-96ee-c9555210d68ca"
        }
      }
    },
    "GetUserGuildsResponse": {
      "type": "object",
      "properties": {
        "guilds": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/GuildData"
          }
        },
        "error": {
          "type": "string"
        }
      }
    },
    "GuildData": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "123456789"
        },
        "name": {
          "type": "string",
          "example": "My Discord Server"
        },
        "icon": {
          "type": "string"
        },
        "owner": {
          "type": "boolean",
          "example": true
        },
        "permissions": {
          "type": "integer",
          "format": "int64",
          "example": 8
        },
        "features": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "GetGuildChannelsRequest": {
      "type": "object",
      "required": ["account_id", "guild_id"],
      "properties": {
        "account_id": {
          "type": "string",
          "example": "9e6a311e-b86b-40dd-96ee-c9555210d68ca"
        },
        "guild_id": {
          "type": "string",
          "example": "123456789"
        }
      }
    },
    "GetGuildChannelsResponse": {
      "type": "object",
      "properties": {
        "channels": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/ChannelData"
          }
        },
        "error": {
          "type": "string"
        }
      }
    },
    "ChannelData": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "987654321"
        },
        "name": {
          "type": "string",
          "example": "general"
        },
        "type": {
          "type": "integer",
          "example": 0
        },
        "position": {
          "type": "integer",
          "example": 0
        },
        "topic": {
          "type": "string"
        }
      }
    },
    "GetGuildRolesRequest": {
      "type": "object",
      "required": ["account_id", "guild_id"],
      "properties": {
        "account_id": {
          "type": "string",
          "example": "9e6a311e-b86b-40dd-96ee-c9555210d68ca"
        },
        "guild_id": {
          "type": "string",
          "example": "123456789"
        }
      }
    },
    "GetGuildRolesResponse": {
      "type": "object",
      "properties": {
        "roles": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/RoleData"
          }
        },
        "error": {
          "type": "string"
        }
      }
    },
    "RoleData": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "555555555"
        },
        "name": {
          "type": "string",
          "example": "Moderator"
        },
        "color": {
          "type": "integer",
          "example": 0
        },
        "position": {
          "type": "integer",
          "example": 1
        },
        "permissions": {
          "type": "integer",
          "format": "int64",
          "example": 0
        },
        "managed": {
          "type": "boolean",
          "example": false
        },
        "mentionable": {
          "type": "boolean",
          "example": true
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
