package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AdminService struct {
	userUc         *biz.UserUsecase
	roleUc         *biz.RoleUsecase
	casbinEnforcer *casbin.Enforcer
	log            *log.Helper
}

var _ admin.IAdminServer = (*AdminService)(nil)

func NewAdminService(
	userUc *biz.UserUsecase,
	roleUc *biz.RoleUsecase,
	casbinEnforcer *casbin.Enforcer,
	logger log.Logger,
) *AdminService {
	return &AdminService{
		userUc:         userUc,
		roleUc:         roleUc,
		casbinEnforcer: casbinEnforcer,
		log:            log.NewHelper(logger),
	}
}

func (s *AdminService) GetUserList(ctx *gin.Context, req *admin.GetUserListRequest) (*admin.GetUserListResponse, error) {
	users, err := s.userUc.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}
	return &admin.GetUserListResponse{Users: users}, nil
}

func (s *AdminService) DeleteRoles(context *gin.Context, request *admin.DeleteRolesRequest) (*emptypb.Empty, error) {
	err := s.roleUc.DeleteRoles(context, request.Names)
	if err != nil {
		return nil, err
	}
	for _, n := range request.Names {
		_, err = s.casbinEnforcer.DeleteUser(n)
		if err != nil {
			s.log.Errorf("DeleteRoles casbinEnforcer %s error: %v", n, err)
		}
	}
	return &emptypb.Empty{}, nil
}

func (s *AdminService) DeleteUsers(context *gin.Context, request *admin.DeleteUsersRequest) (*emptypb.Empty, error) {
	err := s.userUc.DeleteUsers(context, request.Names)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AdminService) GetRoleList(context *gin.Context, request *admin.GetRoleListRequest) (*admin.GetRoleListResponse, error) {
	roles, err := s.roleUc.GetRoleList(context, request)
	if err != nil {
		return nil, err
	}
	return &admin.GetRoleListResponse{Roles: roles}, nil
}

func (s *AdminService) UpdateRoles(context *gin.Context, request *admin.UpdateRolesRequest) (*admin.UpdateRolesResponse, error) {
	rs, err := s.roleUc.UpdateRoles(context, request.Roles)
	if err != nil {
		return nil, err
	}
	return &admin.UpdateRolesResponse{Roles: rs}, nil
}

func (s *AdminService) UpdateUsers(context *gin.Context, request *admin.UpdateUsersRequest) (*admin.UpdateUsersResponse, error) {
	rs, err := s.userUc.UpdateUsers(context, request.Users)
	if err != nil {
		return nil, err
	}
	return &admin.UpdateUsersResponse{Users: rs}, nil
}
