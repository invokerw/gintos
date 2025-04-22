package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/log"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"time"
)

// GinZapLogger Gin 日志中间件（Zap 实现）
func GinZapLogger(logger *log.Helper) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}
		// 根据状态码确定日志级别
		statusCode := c.Writer.Status()
		str := fmt.Sprintf("[GIN] | %3d | %13v | %15s | %-3s %v",
			statusCode,
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
		)
		if statusCode >= 500 {
			logger.Error(str)
		} else {
			logger.Info(str)
		}
	}
}

// GinZapRecovery Gin 异常恢复中间件（Zap 实现）
func GinZapRecovery(logger *log.Helper) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := string(httpRequest)

				// 获取网络连接信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if errors.As(ne.Err, &se) {
						if se.Err.Error() == "broken pipe" {
							brokenPipe = true
						}
					}
				}

				// 记录日志
				logger.Error("[Recovery from panic]",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", headers),
					zap.String("stack", string(debug.Stack())),
				)

				// 如果是 broken pipe 直接终止请求
				if brokenPipe {
					c.Error(err.(error))
					c.Abort()
				} else {
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}
		}()
		c.Next()
	}
}
