package schema

import (
	"github/invokerw/gintos/demo/internal/pkg/mixin"

	"entgo.io/ent/schema/index"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

func (Role) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "roles",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
	}
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Comment("角色名称").
			NotEmpty().
			MaxLen(128),

		field.String("code").
			Comment("角色标识").
			NotEmpty().
			Unique().
			Immutable().
			MaxLen(128),

		field.Int32("sort_id").
			Comment("排序ID").
			Optional().
			Nillable().
			Default(0),
	}
}

// Mixin of the Role.
func (Role) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.Time{},
		mixin.SwitchStatus{},
		mixin.CreateBy{},
		mixin.UpdateBy{},
		mixin.Remark{},
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Indexes of the User.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "code"),
	}
}
