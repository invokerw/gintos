package biz

import (
	"context"
	"fmt"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/conf"
	"github/invokerw/gintos/demo/internal/data/ent"
	"github/invokerw/gintos/demo/internal/pkg/upload"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewOSS, NewGreeterUsecase, NewUserUsecase, NewRoleUsecase, NewCasbinEnforcer)

const (
	FILE_CATEGORY_USER_AVATAR = "user_avatar"
)

func NewOSS(data *conf.File) (upload.OSS, error) {
	switch data.Type {
	case conf.FileType_FILE_TYPE_LOCAL:
		return upload.NewLocalOSS(data.Local.Path, data.Local.Url)
	default:
		panic(fmt.Sprintf("oss type %snot support", data.Type))
	}
}

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	FindByID(context.Context, int64) (*Greeter, error)
	ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
}

// UserRepo is a User repo.
type UserRepo interface {
	CreateUsers(ctx context.Context, users []*common.User) (*ent.User, error)
	GetUser(ctx context.Context, username string) (*ent.User, error)
	GetUserByID(ctx context.Context, id uint64) (*ent.User, error)
	DeleteUsers(ctx context.Context, names []string) error
	UpdateUsers(ctx context.Context, users []*common.User) ([]*ent.User, error)
	GetUserList(ctx context.Context, req *admin.GetUserListRequest) ([]*ent.User, error)
	Count(ctx context.Context) (int, error)
}
type RoleRepo interface {
	CreateRole(ctx context.Context, in *common.Role) (*ent.Role, error)
	GetRole(ctx context.Context, code string) (*ent.Role, error)
	GetRoleByID(ctx context.Context, id uint64) (*ent.Role, error)
	DeleteRoles(ctx context.Context, labels []string) error
	UpdateRoles(ctx context.Context, roles []*common.Role) ([]*ent.Role, error)
	GetRoleList(ctx context.Context, req *admin.GetRoleListRequest) ([]*ent.Role, error)
	Count(ctx context.Context) (int, error)
}
