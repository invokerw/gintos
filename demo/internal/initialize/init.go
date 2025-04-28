package initialize

import (
	"context"
	"github.com/google/wire"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/pkg/trans"
	"github/invokerw/gintos/demo/internal/pkg/utils"
	"github/invokerw/gintos/log"
)

var ProviderSet = wire.NewSet(DoInit)

type InitRet struct {
	Err error
}

func DoInit(user biz.UserRepo, role biz.RoleRepo, l log.Logger) *InitRet {
	ctx := context.Background()
	logger := log.NewHelper(log.With(l, "module", "initialize"))

	createRoleAndUser := func(roleName string, sortId int32, uName, pass string, uAuth common.UserAuthority) error {
		if getRole, _ := role.GetRole(ctx, roleName); getRole == nil {
			_, err := role.CreateRole(ctx, &common.Role{
				Name:     &roleName,
				Desc:     &roleName,
				ParentId: nil,
				SortId:   trans.Ptr(sortId),
			})
			if err != nil {
				logger.Errorf("create %v role failed: %v", roleName, err)
				return err
			}
		}

		if getUser, _ := user.GetUser(ctx, uName); getUser == nil {
			_, err := user.CreateUser(ctx, &common.User{
				RoleName:  &roleName,
				UserName:  trans.Ptr(uName),
				Password:  trans.Ptr(utils.BcryptHash(pass)),
				NickName:  trans.Ptr(uName),
				Email:     trans.Ptr(uName + "@gintos.com"),
				Status:    trans.Ptr(common.UserStatus_ON),
				Authority: trans.Ptr(uAuth),
			})
			if err != nil {
				logger.Errorf("create %v user failed: %v", uName, err)
				return err
			}
		}
		return nil
	}

	// create admin role and user
	{
		var roleName = "admin"
		var name = "admin"
		var pass = "admin123"
		err := createRoleAndUser(roleName, 1, name, pass, common.UserAuthority_SYS_ADMIN)
		if err != nil {
			return &InitRet{
				Err: err,
			}
		}
	}

	// create common role and user
	{
		var roleName = "normal"
		var name = "invoker"
		var pass = "invoker123"
		err := createRoleAndUser(roleName, 100, name, pass, common.UserAuthority_CUSTOMER_USER)
		if err != nil {
			return &InitRet{
				Err: err,
			}
		}
	}

	return nil
}
