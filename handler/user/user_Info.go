package user

import (
	logic "Team2048_Tiktok/logic/user"
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strconv"
)

//获取上下文的用户参数
//	p := NewProxyUserInfo(c)
//	//得到上层中间件根据token解析的userId

//查询该用户的个人信息

//查询数据库中的数据，如果合法则返回个人结构体

// UserInfoHandler 查询用户信息
func UserInfoHandler(c *gin.Context) {

	//1. 获取参数
	var res model.UserInfoResponse
	var p model.ParamUserInfo

	// 获取上下文中保存的当前user_id
	rawId, _ := c.Get("user_id")
	tokenUserId, ok := rawId.(int64)
	if !ok {
		ResponseInfo(c, res, CodeUserIdError)
		return
	}

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
		zap.L().Error("Invalid param", zap.Error(err))
		ResponseInfo(c, res, CodeInvalidParam)
		return
	}
	user_id, err := strconv.ParseInt(p.User_id, 10, 64)

	// 2.获取用户
	tmpUser, err := logic.GetUser(user_id)
	if err != nil {
		// 用户不存在
		zap.L().Error("logic.GetUserID()")
		//res错误码与错误信息返回
		ResponseInfo(c, res, CodeUserNotExisted)
		return
	}

	// 3.完善用户信息结构体（点赞发布数等）
	user, err := logic.GetUserDetail(tmpUser)

	// 4.判断当前用户的关注关系,并进行标注
	logic.CheckIsFollow(&user, tokenUserId)

	// 3.返回响应
	res.User = user
	ResponseInfo(c, res, CodeSuccess)
}
