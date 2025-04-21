package service

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/helloworld"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/log"
)

// GreeterService is a greeter service.
type GreeterService struct {
	uc  *biz.GreeterUsecase
	log *log.Helper
}

var _ helloworld.IGreeterServer = (*GreeterService)(nil)

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase, logger log.Logger) *GreeterService {
	return &GreeterService{uc: uc, log: log.NewHelper(logger)}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx *gin.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &helloworld.HelloReply{Message: "Hello " + g.Hello}, nil
}
