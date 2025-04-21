package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

var _ ent.Mixin = (*CreateBy)(nil)

type CreateBy struct{ mixin.Schema }

func (CreateBy) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("create_by").
			Comment("创建者ID").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*UpdateBy)(nil)

type UpdateBy struct{ mixin.Schema }

func (UpdateBy) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("update_by").
			Comment("更新者ID").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*DeleteBy)(nil)

type DeleteBy struct{ mixin.Schema }

func (DeleteBy) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("delete_by").
			Comment("删除者ID").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*CreatedBy)(nil)

type CreatedBy struct{ mixin.Schema }

func (CreatedBy) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("created_by").
			Comment("创建者ID").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*UpdatedBy)(nil)

type UpdatedBy struct{ mixin.Schema }

func (UpdatedBy) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("updated_by").
			Comment("更新者ID").
			Optional().
			Nillable(),
	}
}

var _ ent.Mixin = (*DeletedBy)(nil)

type DeletedBy struct{ mixin.Schema }

func (DeletedBy) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("deleted_by").
			Comment("删除者ID").
			Optional().
			Nillable(),
	}
}
