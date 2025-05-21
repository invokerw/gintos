package service

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/base"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/demo/internal/pkg/utils"
	"github/invokerw/gintos/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type BaseService struct {
	uc  *biz.UserUsecase
	log *log.Helper
}

var _ base.IBaseServer = (*BaseService)(nil)

func NewBaseService(uc *biz.UserUsecase, logger log.Logger) *BaseService {
	return &BaseService{uc: uc, log: log.NewHelper(logger)}
}

func (s *BaseService) GetMe(ctx *gin.Context, _ *emptypb.Empty) (*base.GetMeResponse, error) {
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

func (s *BaseService) UpdateAvatar(context *gin.Context, request *base.UpdateAvatarRequest) (*base.UpdateAvatarResponse, error) {
	id, err := utils.GetUserID(context)
	if err != nil {
		return nil, errs.ErrUserNotFound.Wrap(err)
	}
	if request.AvatarData == "" {
		return nil, errs.ErrAvatarDataWrong
	}

	data, ext, err := utils.DecodeImageDataURI(request.AvatarData)
	if err != nil {
		return nil, err
	}
	u, err := s.uc.UpdateUserAvatar(context, id, bytes.NewReader(data), ext)
	if err != nil {
		return nil, err
	}
	return &base.UpdateAvatarResponse{User: u}, nil
}

func (s *BaseService) UpdateMe(context *gin.Context, request *base.UpdateMeRequest) (*base.UpdateMeResponse, error) {
	id, err := utils.GetUserID(context)
	if err != nil {
		return nil, errs.ErrUserNotFound.Wrap(err)
	}
	request.User.Id = &id
	user, err := s.uc.UpdateUsers(context, []*common.User{request.User}, true)
	if err != nil {
		return nil, err
	}
	return &base.UpdateMeResponse{User: user[0]}, nil
}
