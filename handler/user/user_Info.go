package user

import (
	logic "Team2048_Tiktok/logic/user"
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// UserInfoHandler 查询用户信息
// @Summary 查询用户信息
// @Description 查询指定用户信息的接口
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param user_id query int true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.UserInfoResponse "查询到的用户信息"
// @Router /douyin/user/ [get]
func UserInfoHandler(c *gin.Context) {

	//1. 获取参数
	var res model.UserInfoResponse
	userInfoId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	// 获取上下文中保存的当前user_id
	rawId, _ := c.Get("user_id")
	tokenUserId, ok := rawId.(int64)
	if !ok {
		ResponseInfo(c, res, CodeUserIdError)
		return
	}

	// 2.获取用户
	tmpUser, err := logic.GetUser(userInfoId)
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
