package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//获取上下文的用户参数
//	p := NewProxyUserInfo(c)
//	//得到上层中间件根据token解析的userId

//查询该用户的个人信息

//查询数据库中的数据，如果合法则返回个人结构体

// UserInfoHandler 查询用户信息
func UserInfoHandler(c *gin.Context) {

	data, err := logic.GetUserInfo()
	if err != nil {
		zap.L().Error("logic.GetUserInfo() failed", zap.Error(err))
		return
	}

	//ResponseSuccess(c, data)
}

// UserInfoHandler 用户登录
func UserLoginHandler(c *gin.Context) {

}

// UserRegisterHandler 用户注册
func UserRegisterHandler(c *gin.Context) {

}
