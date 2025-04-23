package service

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/auth"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/demo/internal/pkg/trans"
	"github/invokerw/gintos/demo/internal/pkg/utils"
	"github/invokerw/gintos/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
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
		User:             user,
		AccessToken:      token,
		RefreshToken:     rToken,
		ExpiresAt:        claims.RegisteredClaims.ExpiresAt.Unix(),
		RefreshExpiresAt: rClaims.RegisteredClaims.ExpiresAt.Unix(),
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
	}
	token, claims, err := utils.CreateAccessToken(baseC)
	if err != nil {
		return nil, err
	}
	utils.SetAccessToken(ctx, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
	return &auth.RefreshTokenResponse{
		AccessToken:      token,
		ExpiresAt:        claims.RegisteredClaims.ExpiresAt.Unix(),
		RefreshToken:     req.RefreshToken,
		RefreshExpiresAt: rClaims.RegisteredClaims.ExpiresAt.Unix(),
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
