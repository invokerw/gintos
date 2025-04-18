package main

import (
	v1 "github/invokerw/gintos/demo/api/helloworld/v1"

	"github.com/gin-gonic/gin"
)

type GreeterServer struct {
}

// SayHello implements helloworld.v1.GreeterServer
func (s *GreeterServer) SayHello(ctx *gin.Context, req *v1.HelloRequest) (*v1.HelloReply, error) {
	reply := &v1.HelloReply{
		Message: "Hello " + req.Name,
	}
	ctx.ShouldBindUri(req)
	return reply, nil
}

func main() {
	router := gin.Default()
	g := router.Group("/").Use(gin.Logger())
	v1.RegisterGreeterServer(g, &GreeterServer{})
	// 监听并在 0.0.0.0:8080 上启动服务
	router.Run(":8080")
}
