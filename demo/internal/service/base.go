package service

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/base"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/demo/internal/pkg/utils"
	"github/invokerw/gintos/log"
)

type BaseService struct {
	uc  *biz.UserUsecase
	log *log.Helper
}

var _ base.IBaseServer = (*BaseService)(nil)

func NewBaseService(uc *biz.UserUsecase, logger log.Logger) *BaseService {
	return &BaseService{uc: uc, log: log.NewHelper(logger)}
}

func (s *BaseService) GetMe(ctx *gin.Context, req *base.GetMeRequest) (*base.GetMeResponse, error) {
	info := utils.GetUserInfo(ctx)
	if info == nil {
		return nil, errs.ErrUserNotFound
	}
	user, err := s.uc.GetUserByID(ctx, info.BaseClaims.ID, true)
	if err != nil {
		return nil, err
	}
	return &base.GetMeResponse{User: user}, nil
}
