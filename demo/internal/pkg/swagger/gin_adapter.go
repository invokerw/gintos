package swaggerUI

import (
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/internal/pkg/swagger/internal/swagger"
	"path"
	"strings"
)

func RegisterGinSwaggerUIServerWithOption(r gin.IRoutes, handlerOpts ...HandlerOption) {
	opts := swagger.NewConfig()

	for _, o := range handlerOpts {
		o(opts)
	}

	var _openJsonFileHandler = &openJsonFileHandler{}
	var pattern string
	if opts.LocalOpenApiFile != "" {
		pattern = strings.TrimRight(opts.BasePath, "/") + "/openapi" + path.Ext(opts.LocalOpenApiFile)
		_ = _openJsonFileHandler.LoadFile(opts.LocalOpenApiFile)
	} else if len(opts.OpenApiData) != 0 {
		pattern = strings.TrimRight(opts.BasePath, "/") + "/openapi." + opts.OpenApiDataType
		_openJsonFileHandler.Content = opts.OpenApiData
	}
	opts.SwaggerJSON = pattern

	swaggerHandler := newHandlerWithConfig(opts)
	r.GET(swaggerHandler.BasePath+"/*action", func(context *gin.Context) {
		if context.Request.URL.Path == pattern {
			_openJsonFileHandler.ServeHTTP(context.Writer, context.Request)
			return
		}
		swaggerHandler.ServeHTTP(context.Writer, context.Request)
	})
}
