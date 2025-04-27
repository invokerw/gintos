package service

import (
	"github/invokerw/gintos/demo/api/v1/auth"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/demo/internal/pkg/trans"
	"github/invokerw/gintos/demo/internal/pkg/utils"
	"github/invokerw/gintos/log"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AuthService is a greeter service.
type AuthService struct {
	uc  *biz.UserUsecase
	log *log.Helper
}

var _ auth.IAuthServer = (*AuthService)(nil)

// NewAuthService new a greeter service.
func NewAuthService(uc *biz.UserUsecase, logger log.Logger) *AuthService {
	return &AuthService{
		uc:  uc,
		log: log.NewHelper(log.With(logger, "service", "auth")),
	}
}

func (s *AuthService) Login(ctx *gin.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := s.uc.GetUser(ctx, req.Username)
	if err != nil {
		return nil, errs.ErrUserNotFound.Wrap(err)
	}

	if user.Password != nil && !utils.BcryptCheck(req.Password, *user.Password) {
		return nil, errs.ErrUserPasswordWrong
	}

	baseC := utils.BaseClaims{
		ID:          user.GetId(),
		Username:    user.GetUserName(),
		NickName:    user.GetNickName(),
		AuthorityId: int32(user.GetAuthority()),
		Role:        "admin",
	}
	token, claims, err := utils.CreateAccessToken(baseC)
	if err != nil {
		return nil, err
	}
	utils.SetAccessToken(ctx, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))

	// 生成刷新令牌
	rToken, rClaims, err := utils.CreateRefreshToken(baseC)
	if err != nil {
		return nil, err
	}

	user.Password = nil
	return &auth.LoginResponse{
		User:           user,
		AccessToken:    token,
		RefreshToken:   rToken,
		Expires:        claims.RegisteredClaims.ExpiresAt.UnixMilli(),
		RefreshExpires: rClaims.RegisteredClaims.ExpiresAt.UnixMilli(),
	}, nil
}

func (s *AuthService) Logout(ctx *gin.Context, req *auth.LogoutRequest) (*emptypb.Empty, error) {
	utils.ClearAccessToken(ctx)
	return &emptypb.Empty{}, nil
}

// RefreshToken 刷新认证令牌
func (s *AuthService) RefreshToken(ctx *gin.Context, req *auth.RefreshTokenRequest) (*auth.RefreshTokenResponse, error) {
	rClaims, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if rClaims.RegisteredClaims.ExpiresAt.Before(time.Now()) {
		return nil, errs.ErrTokenExpired
	}
	user, err := s.uc.GetUserByID(ctx, rClaims.BaseClaims.ID)
	if err != nil {
		return nil, err
	}
	baseC := utils.BaseClaims{
		ID:          user.GetId(),
		Username:    user.GetUserName(),
		NickName:    user.GetNickName(),
		AuthorityId: int32(user.GetAuthority()),
		Role:        "admin",
	}
	token, claims, err := utils.CreateAccessToken(baseC)
	if err != nil {
		return nil, err
	}
	utils.SetAccessToken(ctx, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	return &auth.RefreshTokenResponse{
		User:           user,
		AccessToken:    token,
		Expires:        claims.RegisteredClaims.ExpiresAt.UnixMilli(),
		RefreshToken:   req.RefreshToken,
		RefreshExpires: rClaims.RegisteredClaims.ExpiresAt.UnixMilli(),
	}, nil
}

func (s *AuthService) Register(ctx *gin.Context, req *auth.RegisterRequest) (*emptypb.Empty, error) {
	user, err := s.uc.CreateUser(ctx, &common.User{
		UserName: &req.Username,
		Password: trans.Ptr(utils.BcryptHash(req.Password)),
		Email:    &req.Email,
	})
	if err != nil {
		return nil, err
	}
	_ = user
	s.log.Infof("注册用户成功, 用户名: %s", req.Username)
	return nil, nil
}

func (s *AuthService) GetAsyncRoutes(ctx *gin.Context, req *emptypb.Empty) (*auth.GetAsyncRoutesResponse, error) {
	return &auth.GetAsyncRoutesResponse{
		Routes: []*auth.RouteConfig{
			{
				Path: "/permission",
				Meta: &auth.RouteMeta{
					Title: "权限管理 >.<",
					Icon:  "ep:lollipop",
					Rank:  10,
				},
				Children: []*auth.RouteConfig{
					{
						Path: "/permission/page/index",
						Name: "PermissionPage",
						Meta: &auth.RouteMeta{Title: "页面权限", Roles: []string{"admin", "common"}},
					},
					{
						Path: "/permission/button",
						Meta: &auth.RouteMeta{Title: "按钮权限", Roles: []string{"admin", "common"}},
						Children: []*auth.RouteConfig{
							{
								Path:      "/permission/button/router",
								Component: "permission/button/index",
								Name:      "PermissionButtonRouter",
								Meta: &auth.RouteMeta{
									Title: "路由返回按钮权限",
									Auths: []string{
										"permission:btn:add",
										"permission:btn:edit",
										"permission:btn:delete",
									},
								},
							},
							{
								Path:      "/permission/button/login",
								Component: "permission/button/perms",
								Name:      "PermissionButtonLogin",
								Meta: &auth.RouteMeta{
									Title: "登录接口返回按钮权限",
								},
							},
						},
					},
				},
			},
		},
	}, nil
}
