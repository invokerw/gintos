package service

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/auth"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/log"
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
		return nil, errs.ErrUserNotFound
	}

	if user.Password != nil && *user.Password != req.Password {
		return nil, errs.ErrUserPasswordWrong
	}

	// TODO: 生成token

	return &auth.LoginResponse{}, nil
}

func (s *AuthService) Logout(ctx *gin.Context, req *auth.LogoutRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

// RefreshToken 刷新认证令牌
func (s *AuthService) RefreshToken(ctx *gin.Context, req *auth.RefreshTokenRequest) (*auth.LoginResponse, error) {
	panic("implement me")
}

func (s *AuthService) Register(ctx *gin.Context, req *auth.RegisterRequest) (*emptypb.Empty, error) {
	user, err := s.uc.CreateUser(ctx, &common.User{
		UserName: &req.Username,
		Password: &req.Password,
		Email:    &req.Email,
	})
	if err != nil {
		return nil, err
	}
	_ = user
	s.log.Infof("注册用户成功, 用户名: %s", req.Username)
	return nil, nil
}
