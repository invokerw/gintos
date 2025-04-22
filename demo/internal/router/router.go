package router

import (
	"github/invokerw/gintos/demo/api/v1/auth"
	"github/invokerw/gintos/demo/api/v1/helloworld"
	"github/invokerw/gintos/demo/internal/conf"
	"github/invokerw/gintos/demo/internal/service"
	"github/invokerw/gintos/log"

	"github.com/google/wire"

	"github.com/gin-gonic/gin"
)

var ProviderSet = wire.NewSet(NewGinHttpServer)

func NewGinHttpServer(c *conf.Server, greeter *service.GreeterService, a *service.AuthService, logger log.Logger) *gin.Engine {
	engine := gin.Default()
	{
		g := engine.Group("/")
		g.GET("/", func(c *gin.Context) {
			c.String(200, "Hello World")
		})
	}
	{
		g := engine.Group("/")
		helloworld.RegisterGreeterServer(g, greeter)
	}
	{
		g := engine.Group("/")
		auth.RegisterAuthServer(g, a)
	}
	// swagger
	{
		g := engine.Group("/")
		registerSwaggerOpenApi(g)
	}
	return engine
}
