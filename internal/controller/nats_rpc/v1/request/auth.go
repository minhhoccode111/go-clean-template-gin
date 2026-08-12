package request

import "github.com/goccy/go-json"

// AuthenticatedRequest is the envelope for all authenticated RPC calls.
type AuthenticatedRequest struct {
	Token string          `json:"token" validate:"required"`
	Data  json.RawMessage `json:"data"`
}

// Register -.
type Register struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// Login -.
type Login struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
