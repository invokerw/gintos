package mixin

import (
	"github/invokerw/gintos/demo/internal/pkg/sonyflake"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

type SnowflackId struct {
	mixin.Schema
}

func (SnowflackId) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Comment("id").
			DefaultFunc(sonyflake.GenerateSonyflake).
			Positive().
			Immutable().
			StructTag(`json:"id,omitempty"`).
			SchemaType(map[string]string{
				dialect.MySQL:    "bigint",
				dialect.Postgres: "bigint",
			}),
	}
}

// Indexes of the SnowflackId.
func (SnowflackId) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id"),
	}
}
