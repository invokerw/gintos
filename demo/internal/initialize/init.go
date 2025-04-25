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

func DoInit(user biz.UserRepo, l log.Logger) *InitRet {
	ctx := context.Background()
	logger := log.NewHelper(log.With(l, "module", "initialize"))
	var adminName = "admin"
	var adminPass = "admin123"
	getUser, _ := user.GetUser(ctx, adminName)
	if getUser == nil {
		_, err := user.CreateUser(ctx, &common.User{
			RoleId:    nil,
			UserName:  trans.Ptr(adminName),
			Password:  trans.Ptr(utils.BcryptHash(adminPass)),
			NickName:  trans.Ptr(adminName),
			Email:     trans.Ptr(adminName + "@gintos.com"),
			Status:    trans.Ptr(common.UserStatus_ON),
			Authority: trans.Ptr(common.UserAuthority_SYS_ADMIN),
			Roles:     nil,
		})
		if err != nil {
			logger.Errorf("create admin user failed: %v", err)
			return &InitRet{
				Err: err,
			}
		}
	}

	return nil
}
