package biz

import (
	"context"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/data/ent"
	"github/invokerw/gintos/demo/internal/pkg/trans"
	"github/invokerw/gintos/log"

	"github.com/gin-gonic/gin"
)

type RoleUsecase struct {
	repo RoleRepo
	log  *log.Helper
}

func NewRoleUsecase(repo RoleRepo, logger log.Logger) *RoleUsecase {
	return &RoleUsecase{repo: repo, log: log.NewHelper(log.With(logger, "usecase", "user"))}
}

func (uc *RoleUsecase) CreateRole(ctx *gin.Context, user *common.Role) (*common.Role, error) {
	u, err := uc.repo.CreateRole(ctx, user)
	if err != nil {
		return nil, err
	}
	return uc.convertToRole(u), nil
}

func (uc *RoleUsecase) GetRole(ctx *gin.Context, username string) (*common.Role, error) {
	u, err := uc.repo.GetRole(ctx, username)
	if err != nil {
		return nil, err
	}
	return uc.convertToRole(u), nil
}

func (uc *RoleUsecase) GetRoleByID(ctx *gin.Context, id uint64) (*common.Role, error) {
	u, err := uc.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.convertToRole(u), nil
}

func (uc *RoleUsecase) DeleteRole(ctx *gin.Context, username string) error {
	err := uc.repo.DeleteRole(ctx, username)
	if err != nil {
		return err
	}
	return nil
}

func (uc *RoleUsecase) UpdateRole(ctx context.Context, user *common.Role) (*common.Role, error) {
	u, err := uc.repo.UpdateRole(ctx, user)
	if err != nil {
		return nil, err
	}
	return uc.convertToRole(u), nil
}

func (uc *RoleUsecase) GetRoleList(ctx *gin.Context, req *admin.GetRoleListRequest) ([]*common.Role, error) {
	users, err := uc.repo.GetRoleList(ctx, req)
	if err != nil {
		return nil, err
	}
	var userList []*common.Role
	for _, u := range users {
		userList = append(userList, uc.convertToRole(u))
	}
	return userList, nil
}

func (uc *RoleUsecase) convertToRole(u *ent.Role) *common.Role {
	return &common.Role{
		Id:         trans.Uint64(u.ID),
		Name:       &u.Name,
		Desc:       u.Desc,
		ParentId:   u.ParentID,
		SortId:     u.SortID,
		CreateTime: trans.Ptr(u.CreateTime.Unix()),
		UpdateTime: trans.Ptr(u.UpdateTime.Unix()),
	}
}
