package initialize

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/google/wire"
	"github/invokerw/gintos/demo/api"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/pkg/trans"
	"github/invokerw/gintos/demo/internal/pkg/utils"
	"github/invokerw/gintos/log"
	"math/rand"
)

var ProviderSet = wire.NewSet(DoInit)

type InitRet struct {
	Err error
}

func DoInit(user biz.UserRepo, role biz.RoleRepo, enforce *casbin.Enforcer, l log.Logger) *InitRet {
	ctx := context.Background()
	logger := log.NewHelper(log.With(l, "module", "initialize"))
	createRole := func(roleName string, sortId int32, rules [][]string) error {
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
			if rules != nil {
				if _, err = enforce.AddPolicies(rules); err != nil {
					logger.Errorf("add policies failed: %s: %v", roleName, err)
					return err
				}
			}
		}
		return nil
	}
	createUser := func(roleName string, uName, pass string, uAuth common.UserAuthority) error {
		if getUser, _ := user.GetUser(ctx, uName); getUser == nil {
			// 0-2 rand
			r := rand.Int31n(3)
			_, err := user.CreateUser(ctx, &common.User{
				RoleName:  &roleName,
				UserName:  trans.Ptr(uName),
				Password:  trans.Ptr(utils.BcryptHash(pass)),
				NickName:  trans.Ptr(uName),
				Email:     trans.Ptr(uName + "@gintos.com"),
				Gender:    trans.Ptr(common.UserGender(r)),
				Phone:     trans.Ptr(fmt.Sprintf("1380001000%d", r)),
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
		if err := createRole(roleName, 1, nil); err != nil {
			return &InitRet{
				Err: err,
			}
		}
		if err := createUser(roleName, name, pass, common.UserAuthority_SYS_ADMIN); err != nil {
			return &InitRet{
				Err: err,
			}
		}
	}

	{
		var roleName = "manager"
		var name = "invoker"
		var pass = "invoker123"
		info := api.ApiTypeMap["user"]
		rules := make([][]string, 0, len(info))
		for _, v := range info {
			rules = append(rules, []string{roleName, v.Path, v.Method})
		}
		if err := createRole(roleName, 100, rules); err != nil {
			return &InitRet{
				Err: err,
			}
		}
		if err := createUser(roleName, name, pass, common.UserAuthority_SYS_MANAGER); err != nil {
			return &InitRet{
				Err: err,
			}
		}
	}

	{
		var roleName = "normal"
		var name = "user"
		var pass = "user123"
		if err := createRole(roleName, 100, nil); err != nil {
			return &InitRet{
				Err: err,
			}
		}
		if err := createUser(roleName, name, pass, common.UserAuthority_CUSTOMER_USER); err != nil {
			return &InitRet{
				Err: err,
			}
		}
	}

	return nil
}
