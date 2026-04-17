package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// History holds the schema definition for the translation history table.
type History struct {
	ent.Schema
}

// Fields of the History.
func (History) Fields() []ent.Field {
	return []ent.Field{
		field.String("source").NotEmpty(),
		field.String("destination").NotEmpty(),
		field.String("original").NotEmpty(),
		field.String("translation").NotEmpty(),
	}
}

// Edges of the History.
func (History) Edges() []ent.Edge {
	return nil
}

// Annotations of the History.
func (History) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "history"},
	}
}
