package biz

import (
	"context"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/data/ent"
	"github/invokerw/gintos/demo/internal/data/ent/role"
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

func (uc *RoleUsecase) GetRole(ctx *gin.Context, code string) (*common.Role, error) {
	u, err := uc.repo.GetRole(ctx, code)
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

func (uc *RoleUsecase) DeleteRoles(ctx *gin.Context, labels []string) error {
	err := uc.repo.DeleteRoles(ctx, labels)
	if err != nil {
		return err
	}
	return nil
}

func (uc *RoleUsecase) UpdateRoles(ctx context.Context, roles []*common.Role) ([]*common.Role, error) {
	us, err := uc.repo.UpdateRoles(ctx, roles)
	if err != nil {
		return nil, err
	}
	return uc.convertToRoles(us), nil
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

func (uc *RoleUsecase) GetRoleCount(ctx context.Context) (int, error) {
	count, err := uc.repo.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (uc *RoleUsecase) convertToRoles(u []*ent.Role) []*common.Role {
	ret := make([]*common.Role, 0, len(u))
	for _, role := range u {
		ret = append(ret, uc.convertToRole(role))
	}
	return ret
}

func (uc *RoleUsecase) convertToStatus(s *role.Status) *common.RoleStatus {
	if s == nil {
		return nil
	}
	var ret common.RoleStatus
	switch *s {
	case role.StatusON:
		ret = common.RoleStatus_R_ON
	case role.StatusOFF:
		ret = common.RoleStatus_R_OFF
	default:
		return nil
	}
	return trans.Ptr(ret)
}

func (uc *RoleUsecase) convertToRole(u *ent.Role) *common.Role {
	return &common.Role{
		Id:         trans.Uint64(u.ID),
		Name:       &u.Name,
		Code:       &u.Code,
		SortId:     u.SortID,
		Remark:     u.Remark,
		Status:     uc.convertToStatus(u.Status),
		CreateTime: trans.Ptr(u.CreateTime.Unix()),
		UpdateTime: trans.Ptr(u.UpdateTime.Unix()),
	}
}
