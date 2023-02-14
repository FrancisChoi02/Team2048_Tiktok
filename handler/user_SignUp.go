package handler

import (
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

	//2. 查找Username是否已经存在
	if boolValue := logic.CheckUser(p.Username); boolValue {
		// 用户已存在
		zap.L().Error("SignUp with invalid param")
		//res错误码与错误信息返回
		ResponseSignUp(c, res, CodeUserN)
		return
	}

	//3. 分发用户ID
	tmpID := model.GenID()

	//4. 将新用户信息添加到数据库中
	if err := logic.SignUp(tmpID, p.Username, p.Password); err != nil {
		zap.L().Error("logic.SignUp() failed", zap.Error(err))
	}

	//5. 颁发token
	token, err := middleware.ReleaseToken(tmpID)
	if err != nil {
		zap.L().Error("middleware.ReleaseToken() failed", zap.Error(err))
		ResponseSignUp(c, res, CodeSignUpSuccess)
	}

	//6. 返回注册成功的响应
	res.User_id = tmpID
	res.Token = token
	ResponseSignUp(c, res, CodeSignUpSuccess)
}
