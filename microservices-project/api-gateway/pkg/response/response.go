package response

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *app.RequestContext, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *app.RequestContext, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *app.RequestContext, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 400 错误请求
func BadRequest(c *app.RequestContext, message string) {
	c.JSON(400, Response{
		Code:    400,
		Message: message,
	})
}

// Unauthorized 401 未授权
func Unauthorized(c *app.RequestContext, message string) {
	c.JSON(401, Response{
		Code:    401,
		Message: message,
	})
}

// Forbidden 403 禁止访问
func Forbidden(c *app.RequestContext, message string) {
	c.JSON(403, Response{
		Code:    403,
		Message: message,
	})
}

// NotFound 404 未找到
func NotFound(c *app.RequestContext, message string) {
	c.JSON(404, Response{
		Code:    404,
		Message: message,
	})
}

// InternalServerError 500 服务器错误
func InternalServerError(c *app.RequestContext, message string) {
	c.JSON(500, Response{
		Code:    500,
		Message: message,
	})
}

// CustomError 自定义错误
func CustomError(c *app.RequestContext, httpCode int, code int, message string) {
	c.JSON(httpCode, Response{
		Code:    code,
		Message: message,
	})
}
