package service

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/auth"
	"github/invokerw/gintos/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AuthService is a greeter service.
type AuthService struct {
	log *log.Helper
}

var _ auth.IAuthServer = (*AuthService)(nil)

// NewAuthService new a greeter service.
func NewAuthService(logger log.Logger) *AuthService {
	return &AuthService{log: log.NewHelper(logger)}
}

func (s *AuthService) Login(*gin.Context, *auth.LoginRequest) (*auth.LoginResponse, error) {
	panic("implement me")
}

func (s *AuthService) Logout(*gin.Context, *auth.LogoutRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

// RefreshToken 刷新认证令牌
func (s *AuthService) RefreshToken(*gin.Context, *auth.RefreshTokenRequest) (*auth.LoginResponse, error) {
	panic("implement me")
}

func (s *AuthService) Register(*gin.Context, *auth.RegisterRequest) (*emptypb.Empty, error) {
	panic("implement me")
}
