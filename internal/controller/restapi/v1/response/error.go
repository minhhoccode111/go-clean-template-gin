package response

type Error struct {
	Error string `json:"error" example:"message"`
} // @name v1.Error

// Message is a generic JSON response body for success/error messages.
// The validate:"required" tag ensures generated TS types are non-optional.
type Message struct {
	Message string `json:"message" validate:"required"`
}
