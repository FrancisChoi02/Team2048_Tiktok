package handler

import (
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

//获取上下文的用户参数
//	p := NewProxyUserInfo(c)
//	//得到上层中间件根据token解析的userId

//查询该用户的个人信息

//查询数据库中的数据，如果合法则返回个人结构体

// UserInfoHandler 查询用户信息
func UserInfoHandler(c *gin.Context) {

	//1. 获取参数
	var p model.ParamUserInfo
	var res model.UserInfoResponse

	if err := c.ShouldBind(&p); err != nil {
		// 参数存在空值
		// 获取validator.ValidationErrors类型的errors
		_, ok := err.(validator.ValidationErrors)
		if !ok {
			zap.L().Error("Invalid param", zap.Error(err))
			ResponseInfo(c, res, CodeInvalidParam) //状态码为 参数错误
			return
		}
		// 参数错误
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		ResponseInfo(c, res, CodeInvalidParam)
		return
	}

	// 2.获取用户
	tmpUser, err := logic.GetUser_ByID(p.User_id)
	if err != nil {
		// 用户不存在
		zap.L().Error("logic.GetUserID()")
		//res错误码与错误信息返回
		ResponseInfo(c, res, CodeUserNotExisted)
		return
	}

	// 3.填充用户结构体
	resUser := new(model.User)

	// 4.返回响应
	resUser.Id = tmpUser.Id
	resUser.Name = tmpUser.Name
	resUser.FollowCount = tmpUser.FollowCount
	resUser.FollowerCount = tmpUser.FollowerCount
	resUser.IsFollow = tmpUser.IsFollow
	ResponseInfo(c, res, CodeSignUpSuccess)

}
