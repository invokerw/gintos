package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

func Result(c *gin.Context, code int, data interface{}, msg string) {
	// 开始时间
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(c, SUCCESS, map[string]interface{}{}, "操作成功")
}

func OkWithMessage(c *gin.Context, message string) {
	Result(c, SUCCESS, map[string]interface{}{}, message)
}

func OkWithData(c *gin.Context, data interface{}) {
	Result(c, SUCCESS, data, "成功")
}

func OkWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, SUCCESS, data, message)
}

func Fail(c *gin.Context) {
	Result(c, ERROR, map[string]interface{}{}, "操作失败")
}

func FailWithMessage(c *gin.Context, message string) {
	Result(c, ERROR, map[string]interface{}{}, message)
}

func NoAuth(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		ERROR,
		nil,
		message,
	})
}

func FailWithDetailed(c *gin.Context, data interface{}, message string) {
	Result(c, ERROR, data, message)
}
