package router

import (
	v1 "github/invokerw/gintos/demo/api/helloworld/v1"
	"github/invokerw/gintos/demo/internal/conf"
	"github/invokerw/gintos/demo/internal/service"
	"github/invokerw/gintos/log"

	"github.com/google/wire"

	"github.com/gin-gonic/gin"
)

var ProviderSet = wire.NewSet(NewGinHttpServer)

func NewGinHttpServer(c *conf.Server, greeter *service.GreeterService, logger log.Logger) *gin.Engine {
	engine := gin.Default()
	g := engine.Group("/").Use(gin.Logger())
	v1.RegisterGreeterServer(g, greeter)
	return engine
}
