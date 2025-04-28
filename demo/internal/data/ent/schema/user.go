package schema

import (
	"entgo.io/ent/schema/edge"
	"github/invokerw/gintos/demo/internal/pkg/mixin"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:     "users",
			Charset:   "utf8mb4",
			Collation: "utf8mb4_bin",
		},
		entsql.WithComments(true),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Comment("用户名").
			Unique().
			NotEmpty().
			Immutable().
			Optional().
			Nillable(),

		field.String("password").
			Comment("登录密码").
			MaxLen(255).
			Optional().
			Nillable().
			NotEmpty(),

		field.String("nick_name").
			Comment("昵称").
			MaxLen(255).
			Optional().
			Nillable(),

		field.String("email").
			Comment("电子邮箱").
			MaxLen(320).
			Optional().
			Nillable(),

		field.String("mobile").
			Comment("手机号码").
			Default("").
			MaxLen(255).
			Optional().
			Nillable(),

		field.String("avatar").
			Comment("头像").
			MaxLen(1023).
			Optional().
			Nillable(),

		field.Enum("gender").
			Comment("性别").
			Values(
				"UNKNOWN",
				"MALE",
				"FEMALE",
			).
			Optional().
			Nillable(),

		field.Enum("authority").
			Comment("授权").
			Optional().
			Nillable().
			Values(
				"SYS_ADMIN",
				"SYS_MANAGER",
				"CUSTOMER_USER",
				"GUEST_USER",
				"REFRESH_TOKEN",
			).
			Default("CUSTOMER_USER"),

		field.Int64("last_login_time").
			Comment("最后一次登录的时间").
			Optional().
			Nillable(),
	}
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.AutoIncrementId{},
		mixin.CreateBy{},
		mixin.UpdateBy{},
		mixin.Time{},
		mixin.Remark{},
		mixin.SwitchStatus{},
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role", Role.Type).Unique(),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id", "username"),
	}
}
