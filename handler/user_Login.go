package handler

import (
	"Team2048_Tiktok/middleware"
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

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

	//2. 查找Username是否已经存在
	if boolValue := logic.CheckUser(p.Username); !boolValue {
		// 用户已存在
		zap.L().Error("logic.CheckUser() failed")
		//res错误码与错误信息返回
		ResponseLogin(c, res, CodeUserNotExisted)
		return
	}

	//3. 获取用户ID
	tmpUser, err := logic.GetUser_ByName(p.Username)
	if err != nil {
		// 用户不存在
		zap.L().Error("logic.GetUserID()")
		//res错误码与错误信息返回
		ResponseLogin(c, res, CodeUserNotExisted)
		return
	}

	//4. 颁发token
	token, err := middleware.ReleaseToken(tmpID)
	if err != nil {
		zap.L().Error("middleware.ReleaseToken() failed", zap.Error(err))
		ResponseLogin(c, res, CodeSignUpSuccess)
	}

	//5. 返回登录成功的响应
	res.User_id = tmpUser //此处要改为 tmpUser.Id
	res.Token = token
	ResponseLogin(c, res, CodeSignUpSuccess)
}
