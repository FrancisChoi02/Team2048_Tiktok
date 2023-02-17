package user

import (
	logic "Team2048_Tiktok/logic/user"
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserLoginHandler 用户登录
func UserLoginHandler(c *gin.Context) {
	//1. 获取参数 和 参数校验
	var p model.ParamSignUp          //用于获取 request消息 的结构体
	var res model.UserDetailResponse //用于返回 response消息 的结构体

	if err := c.ShouldBind(&p); err != nil {
		// 登录参数存在空值
		// 获取validator.ValidationErrors类型的errors
		_, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("Invalid param", zap.Error(err))
			ResponseLogin(c, res, CodeInvalidParam) //状态码为 参数错误
			return
		}
		// 参数错误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		ResponseLogin(c, res, CodeInvalidParam)
		return
	}

	//2. 登录 并获取用户Id及颁发的Token
	tmpId, token, err := logic.Login(&p) //这里应该是tmpUser.Id
	if err != nil {
		zap.L().Error("logic.Login() failed", zap.Error(err))
		ResponseLogin(c, res, CodeSuccess)
	}

	//3. 返回登录成功的响应
	res.User_id = tmpId
	res.Token = token
	ResponseLogin(c, res, CodeSuccess)
}
