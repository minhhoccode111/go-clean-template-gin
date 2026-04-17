package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// History holds the schema definition for the History entity.
type History struct {
	ent.Schema
}

// Fields of the History.
func (History) Fields() []ent.Field {
	return []ent.Field{
		field.String("source"),
		field.String("destination"),
		field.String("original"),
		field.String("translation"),
	}
}
