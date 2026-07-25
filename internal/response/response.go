package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// OK 成功响应
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// OKWithMsg 成功响应带消息
func OKWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: msg,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code:    1,
		Message: msg,
	})
}

// FailWithCode 失败响应带状态码
func FailWithCode(c *gin.Context, httpCode int, msg string) {
	c.JSON(httpCode, Response{
		Code:    1,
		Message: msg,
	})
}
