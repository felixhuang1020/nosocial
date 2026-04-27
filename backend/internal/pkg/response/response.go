package response

import (
	"net/http"
	"nosocial/config"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess = 0
	CodeError   = 1
)

// 生产环境通用错误消息，避免泄露内部细节
const (
	GenericErrorMsg      = "服务器内部错误"
	GenericDBErrorMsg    = "数据操作失败"
	GenericAuthErrorMsg  = "认证失败"
	GenericParamErrorMsg = "请求参数错误"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// Error 通用错误响应，生产环境自动脱敏
func Error(c *gin.Context, code int, msg string) {
	if config.C != nil && config.C.App.Mode == "release" {
		msg = GenericErrorMsg
	}
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// SafeError 安全错误包装，生产环境隐藏原始错误详情
func SafeError(c *gin.Context, err error) {
	msg := err.Error()
	if config.C != nil && config.C.App.Mode == "release" {
		msg = GenericErrorMsg
	}
	c.JSON(http.StatusOK, Response{
		Code: CodeError,
		Msg:  msg,
		Data: nil,
	})
}

func ErrorWithStatus(c *gin.Context, httpCode int, msg string) {
	if config.C != nil && config.C.App.Mode == "release" {
		msg = GenericErrorMsg
	}
	c.JSON(httpCode, Response{
		Code: CodeError,
		Msg:  msg,
		Data: nil,
	})
}

func BadRequest(c *gin.Context, msg string) {
	if config.C != nil && config.C.App.Mode == "release" {
		msg = GenericParamErrorMsg
	}
	ErrorWithStatus(c, http.StatusBadRequest, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusUnauthorized, msg)
}

func Forbidden(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusForbidden, msg)
}

func NotFound(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusNotFound, msg)
}

func ServerError(c *gin.Context, msg string) {
	ErrorWithStatus(c, http.StatusInternalServerError, GenericErrorMsg)
}
