package request

import "strings"

type Translate struct {
	Source      string `json:"source"      validate:"required" example:"auto"`
	Destination string `json:"destination" validate:"required" example:"en"`
	Original    string `json:"original"    validate:"required" example:"текст для перевода"`
} // @name v1.Translate

// Clean trims whitespace from all string fields.
func (r *Translate) Clean() {
	r.Source = strings.TrimSpace(r.Source)
	r.Destination = strings.TrimSpace(r.Destination)
	r.Original = strings.TrimSpace(r.Original)
}
