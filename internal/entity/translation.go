// Package entity defines main entities for business logic (services).
// These entities are the "source of truth" and are used throughout the entire application.
// They are independent of database schema or transport protocols (HTTP/gRPC).
package entity

// Translation -.
type Translation struct {
	Source      string `json:"source"       example:"auto"`
	Destination string `json:"destination"  example:"en"`
	Original    string `json:"original"     example:"текст для перевода"`
	Translation string `json:"translation"  example:"text for translation"`
}
