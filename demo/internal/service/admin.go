package service

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/log"
)

type AdminService struct {
	uc  *biz.UserUsecase
	log *log.Helper
}

var _ admin.IAdminServer = (*AdminService)(nil)

func NewAdminService(uc *biz.UserUsecase, logger log.Logger) *AdminService {
	return &AdminService{uc: uc, log: log.NewHelper(logger)}
}

func (s *AdminService) GetUserList(ctx *gin.Context, req *admin.GetUserListRequest) (*admin.GetUserListResponse, error) {
	users, err := s.uc.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}
	return &admin.GetUserListResponse{Users: users}, nil
}
