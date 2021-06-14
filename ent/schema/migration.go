package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Migration holds the schema definition for the Migration entity.
type Migration struct {
	ent.Schema
}

func (Migration) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique().NotEmpty(),
		field.Time("applied_at").Immutable().Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Migration.
func (Migration) Edges() []ent.Edge {
	return nil
}
