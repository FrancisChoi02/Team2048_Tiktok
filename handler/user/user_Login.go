package user

import (
	logic "Team2048_Tiktok/logic/user"
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserLoginHandler 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {object} model.UserDetailResponse
// @Router /douyin/user/login/ [get]
func UserLoginHandler(c *gin.Context) {
	//1. 获取参数 和 参数校验
	var p model.ParamSignUp          //用于获取 request消息 的结构体
	var res model.UserDetailResponse //用于返回 response消息 的结构体

	//参数位置在Query
	p.Username = c.Query("username")
	p.Password = c.Query("password")

	//2. 登录 并获取用户Id及颁发的Token
	tmpId, token, err := logic.Login(&p) //这里应该是tmpUser.Id
	if err != nil {
		zap.L().Error("logic.Login() failed", zap.Error(err))
		ResponseLogin(c, res, CodeUserNotLogin)
		return
	}

	//3. 返回登录成功的响应
	res.User_id = tmpId
	res.Token = token
	ResponseLogin(c, res, CodeSuccess)
}
