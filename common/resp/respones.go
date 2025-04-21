package resp

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IError interface {
	Error() string
	Code() int
	MetaData() map[string]string
}

type Error struct {
	code     int
	message  string
	metadata map[string]string
}

var _ IError = (*Error)(nil)

func NewErr(code int, message string) *Error {
	return &Error{
		code:     code,
		message:  message,
		metadata: make(map[string]string),
	}
}

func (e *Error) Error() string {
	return e.message
}
func (e *Error) Code() int {
	return e.code
}
func (e *Error) MetaData() map[string]string {
	return e.metadata
}

func (e *Error) CloneWithMetadata(metadata map[string]string) *Error {
	if e == nil {
		return nil
	}
	e2 := NewErr(e.code, e.message)
	e2.metadata = metadata
	return e2
}
func (e *Error) Wrap(err error) *Error {
	if e == nil {
		return nil
	}
	e2 := NewErr(e.code, e.message+": "+err.Error())
	return e2
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	SUCCESS = 0
	ERROR   = 1
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

func FailWithError(c *gin.Context, err error) {
	var e IError
	if ok := errors.As(err, &e); ok {
		Result(c, e.Code(), e.MetaData(), e.Error())
	} else {
		FailWithMessage(c, err.Error())
	}
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
