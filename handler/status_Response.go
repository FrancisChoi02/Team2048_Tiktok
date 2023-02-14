package handler

import (
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseError 返回错误消息
func ResponseError(c *gin.Context, code StatusCode) {
	res := &model.CommonResponse{
		StatusCode: int32(code),
		StatusMsg:  code.Msg(),
	}

	c.JSON(http.StatusOK, res)
}

// ResponseErrorWithMsg 返回带指定字符串的错误消息
func ResponseErrorWithMsg(c *gin.Context, code StatusCode, tmpString string) {
	res := &model.CommonResponse{
		StatusCode: int32(code),
		StatusMsg:  tmpString,
	}

	c.JSON(http.StatusOK, res)
}

// ResponseSignUp 返回注册结果
func ResponseSignUp(c *gin.Context, res model.UserDetailResponse, code StatusCode) {
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}

// ResponseLogin 返回登录结果
func ResponseLogin(c *gin.Context, res model.UserDetailResponse, code StatusCode) {
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}

func ResponseInfo(c *gin.Context, res model.UserInfoResponse, code StatusCode) {
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}
