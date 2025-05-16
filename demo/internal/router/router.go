package router

import (
	"github/invokerw/gintos/common/middleware"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/auth"
	"github/invokerw/gintos/demo/api/v1/base"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/api/v1/helloworld"
	"github/invokerw/gintos/demo/assets"
	"github/invokerw/gintos/demo/internal/conf"
	"github/invokerw/gintos/demo/internal/router/mw"
	"github/invokerw/gintos/demo/internal/service"
	"github/invokerw/gintos/log"
	"io/fs"
	"net/http"

	"github.com/casbin/casbin/v2"

	"github.com/google/wire"

	"github.com/gin-gonic/gin"
)

var ProviderSet = wire.NewSet(NewGinHttpServer)

func NewGinHttpServer(c *conf.Server,
	greeter *service.GreeterService,
	a *service.AuthService,
	adminS *service.AdminService,
	bs *service.BaseService,
	enforce *casbin.Enforcer,
	logger log.Logger) *gin.Engine {

	//engine := gin.Default()
	engine := gin.New()
	ginHelper := log.NewHelper(log.With(logger, "module", "router"))
	engine.Use(middleware.GinZapLogger(ginHelper), middleware.GinZapRecovery(ginHelper))
	{
		g := engine.Group("/")
		g.GET("/", func(c *gin.Context) {
			c.String(200, "Hello World")
		})
		st, _ := fs.Sub(assets.Frontend, "frontend/dist")
		g.StaticFS("/admin", http.FS(st))
	}
	{
		g := engine.Group("/")
		helloworld.RegisterGreeterServer(g, greeter)
	}
	{
		g := engine.Group("/")
		auth.RegisterAuthServer(g, a)
	}
	{
		g := engine.Group("/").Use(mw.JWTAuth(), mw.CasbinAuth(common.UserAuthority_SYS_MANAGER, enforce))
		admin.RegisterAdminServer(g, adminS)
	}
	{
		g := engine.Group("/").Use(mw.JWTAuth(), mw.CasbinAuth(common.UserAuthority_CUSTOMER_USER, nil))
		base.RegisterBaseServer(g, bs)
	}
	// swagger
	{
		g := engine.Group("/")
		registerSwaggerOpenApi(g)
	}
	return engine
}
