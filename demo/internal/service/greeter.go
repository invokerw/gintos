package service

import (
	v1 "github/invokerw/gintos/demo/api/helloworld/v1"
	"github/invokerw/gintos/demo/internal/biz"

	"github.com/gin-gonic/gin"
)

// GreeterService is a greeter service.
type GreeterService struct {
	uc *biz.GreeterUsecase
}

var _ v1.IGreeterServer = (*GreeterService)(nil)

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx *gin.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}
