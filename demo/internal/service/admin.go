package service

import (
	"bytes"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/demo/internal/pkg/utils"
	"github/invokerw/gintos/log"
	"github/invokerw/gintos/proto/rbac"
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

func (s *AdminService) CreateUser(context *gin.Context, request *admin.CreateUserRequest) (*admin.CreateUserResponse, error) {
	user, err := s.userUc.CreateUser(context, request.User, true)
	if err != nil {
		return nil, err
	}
	return &admin.CreateUserResponse{User: user}, nil
}

func (s *AdminService) GetUserList(ctx *gin.Context, req *admin.GetUserListRequest) (*admin.GetUserListResponse, error) {
	users, err := s.userUc.GetUserList(ctx, req, true)
	if err != nil {
		return nil, err
	}
	return &admin.GetUserListResponse{Users: users}, nil
}

func (s *AdminService) DeleteRoles(context *gin.Context, request *admin.DeleteRolesRequest) (*emptypb.Empty, error) {
	err := s.roleUc.DeleteRoles(context, request.Codes)
	if err != nil {
		return nil, err
	}
	for _, n := range request.Codes {
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

func (s *AdminService) CreateRole(context *gin.Context, request *admin.CreateRoleRequest) (*admin.CreateRoleResponse, error) {
	role, err := s.roleUc.CreateRole(context, request.Role)
	if err != nil {
		return nil, err
	}
	return &admin.CreateRoleResponse{Role: role}, nil
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
	rs, err := s.userUc.UpdateUsers(context, request.Users, true)
	if err != nil {
		return nil, err
	}
	return &admin.UpdateUsersResponse{Users: rs}, nil
}

func (s *AdminService) GetApiInfoList(context *gin.Context, empty *emptypb.Empty) (*admin.GetApiInfoListResponse, error) {
	ret := &admin.GetApiInfoListResponse{
		ApiTypeMap: make(map[string]*common.ApiTypeInfo),
	}
	for _, v := range api.ApiInfo {
		if ret.ApiTypeMap[v.Type] == nil {
			ret.ApiTypeMap[v.Type] = &common.ApiTypeInfo{
				Type: v.Type,
			}
		}
		ret.ApiTypeMap[v.Type].ApiInfo = append(ret.ApiTypeMap[v.Type].ApiInfo, s.convertApiInfo(v))
	}
	return ret, nil
}

func (s *AdminService) RoleGetPolicy(context *gin.Context, request *admin.RoleGetPolicyRequest) (*admin.RoleGetPolicyResponse, error) {
	p, err := s.casbinEnforcer.GetPermissionsForUser(request.RoleCode)
	if err != nil {
		return nil, err
	}
	ret := &admin.RoleGetPolicyResponse{}
	for _, v := range p {
		if apiInfo, ok := api.GetApiInfo(v[1], v[2]); ok {
			ret.ApiInfo = append(ret.ApiInfo, s.convertApiInfo(apiInfo))
		}
	}
	return ret, nil
}

func (s *AdminService) RoleUpdatePolicy(context *gin.Context, request *admin.RoleUpdatePolicyRequest) (*emptypb.Empty, error) {
	rules, err := s.casbinEnforcer.GetPermissionsForUser(request.RoleCode)
	if err != nil {
		return nil, err
	}
	_, err = s.casbinEnforcer.RemovePolicies(rules)
	if err != nil {
		return nil, err
	}
	rules = nil
	for _, v := range request.ApiName {
		if apiInfo, ok := api.ApiMap[v]; ok {
			rules = append(rules, []string{request.RoleCode, apiInfo.Path, apiInfo.Method})
		}
	}
	if rules != nil {
		_, err = s.casbinEnforcer.AddPolicies(rules)
		if err != nil {
			return nil, err
		}
	}
	err = s.casbinEnforcer.SavePolicy()
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AdminService) GetRoleCount(context *gin.Context, empty *emptypb.Empty) (*common.IntValue, error) {
	count, err := s.roleUc.GetRoleCount(context)
	if err != nil {
		return nil, err
	}
	return &common.IntValue{Data: int32(count)}, nil
}

func (s *AdminService) GetUserCount(context *gin.Context, empty *emptypb.Empty) (*common.IntValue, error) {
	count, err := s.userUc.GetUserCount(context)
	if err != nil {
		return nil, err
	}
	return &common.IntValue{Data: int32(count)}, nil
}

func (s *AdminService) convertApiInfo(i *rbac.ApiInfo) *common.ApiInfo {
	return &common.ApiInfo{
		Name:   i.Name,
		Path:   i.Path,
		Method: i.Method,
	}
}

func (s *AdminService) UpdateUserAvatar(context *gin.Context, request *admin.UpdateUserAvatarRequest) (*admin.UpdateUserAvatarResponse, error) {
	if request.AvatarData == "" {
		return nil, errs.ErrAvatarDataWrong
	}

	data, ext, err := utils.DecodeImageDataURI(request.AvatarData)
	if err != nil {
		return nil, err
	}
	u, err := s.userUc.UpdateUserAvatar(context, request.Id, bytes.NewReader(data), ext)
	if err != nil {
		return nil, err
	}
	return &admin.UpdateUserAvatarResponse{User: u}, nil
}
