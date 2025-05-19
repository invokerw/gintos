package router

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggest/swgui/v5emb"
	"github/invokerw/gintos/demo/assets"
)

func registerSwaggerOpenApi(r gin.IRoutes) {
	// 将 YAML 转换为 JSON
	//swaggerUI.RegisterGinSwaggerUIServerWithOption(r,
	//	swaggerUI.WithTitle("Gintos Admin"),
	//	swaggerUI.WithMemoryData(OpenApiData, "yaml"))

	r.GET("/docs/*any", func(c *gin.Context) {
		if c.Request.URL.Path == "/docs/openapi.yaml" {
			c.String(200, string(assets.OpenApiData))
			return
		}
		h := gin.WrapH(v5emb.New(
			"Gintos API",         // 标题
			"/docs/openapi.yaml", // OpenAPI JSON 文件的地址
			"/docs/",             // Swagger UI 的访问路径
		))
		h(c)
	})
}
