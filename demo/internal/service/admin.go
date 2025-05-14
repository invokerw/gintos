package service

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/log"
	"github/invokerw/gintos/proto/rbac"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
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
	p, err := s.casbinEnforcer.GetPermissionsForUser(request.RoleName)
	if err != nil {
		return nil, err
	}
	ret := &admin.RoleGetPolicyResponse{}
	for _, v := range p {
		k := strings.Join(v, "-")
		if apiInfo, ok := api.ApiPathMethodToApiInfo[k]; ok {
			ret.ApiInfo = append(ret.ApiInfo, s.convertApiInfo(apiInfo))
		}
	}
	return ret, nil
}

func (s *AdminService) RoleUpdatePolicy(context *gin.Context, request *admin.RoleUpdatePolicyRequest) (*emptypb.Empty, error) {
	p, err := s.casbinEnforcer.GetPermissionsForUser(request.RoleName)
	if err != nil {
		return nil, err
	}
	rules := make([][]string, 0, len(p))
	for _, v := range p {
		rules = append(rules, []string{request.RoleName, v[0], v[1]})
	}
	_, err = s.casbinEnforcer.RemovePolicies(rules)
	if err != nil {
		return nil, err
	}
	rules = nil
	for _, v := range request.ApiName {
		if apiInfo, ok := api.ApiMap[v]; ok {
			rules = append(rules, []string{request.RoleName, apiInfo.Path, apiInfo.Method})
		}

	}
	_, err = s.casbinEnforcer.AddPolicies(rules)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AdminService) convertApiInfo(i *rbac.ApiInfo) *common.ApiInfo {
	return &common.ApiInfo{
		Name:   i.Name,
		Path:   i.Path,
		Method: i.Method,
	}
}
