package user

import (
	logic "Team2048_Tiktok/logic/user"
	"Team2048_Tiktok/middleware"
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserRegisterHandler 用户注册
func UserSignUpHandler(c *gin.Context) {
	//1. 获取参数 和 参数校验
	var p model.ParamSignUp          //用于获取 request消息 的结构体
	var res model.UserDetailResponse //用于返回 response消息 的结构体

	if err := c.ShouldBind(&p); err != nil {
		// 注册参数存在空值
		// 获取validator.ValidationErrors类型的errors
		_, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("Invalid param", zap.Error(err))
			ResponseSignUp(c, res, CodeInvalidParam) //状态码为 参数错误
			return
		}
		// 参数错误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		ResponseSignUp(c, res, CodeInvalidParam)
		return
	}

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
	ResponseSignUp(c, res, CodeSignUpSuccess)
}
