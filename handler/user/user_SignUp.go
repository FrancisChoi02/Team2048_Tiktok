package user

import (
	logic "Team2048_Tiktok/logic/user"
	"Team2048_Tiktok/middleware"
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// UserSignUpHandler 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 {object} model.UserDetailResponse
// @Router /douyin/user/register/ [post]
func UserSignUpHandler(c *gin.Context) {
	//1. 获取参数 和 参数校验
	var p model.ParamSignUp          //用于获取 request消息 的结构体
	var res model.UserDetailResponse //用于返回 response消息 的结构体

	//参数位置在Query
	p.Username = c.Query("username")
	p.Password = c.Query("password")

	//2. 将新用户信息添加到数据库中
	tmpId, err := logic.SignUp(&p)
	if err != nil {
		zap.L().Error("logic.SignUp() failed", zap.Error(err))
		return
	}

	//3. 颁发token
	token, err := middleware.ReleaseToken(tmpId)
	if err != nil {
		zap.L().Error("middleware.ReleaseToken() failed", zap.Error(err))
		ResponseSignUp(c, res, CodeServerBusy)
		return
	}

	//3. 返回注册成功的响应
	res.User_id = tmpId
	res.Token = token
	ResponseSignUp(c, res, CodeSuccess)
}
