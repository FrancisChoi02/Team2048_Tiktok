package user

import (
	logic "Team2048_Tiktok/logic/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// RelationshipHandler  用户关系操作
func RelationshipHandler(c *gin.Context) {

	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseRelation(c, CodeUserIdError)
		return
	}
	rawToUser := c.PostForm("to_user_id")
	tmpToUser, err := strconv.Atoi(rawToUser)
	if err != nil {
		ResponseRelation(c, CodeUserIdError)
		return
	}

	actionType, err := strconv.Atoi(c.PostForm("action_type")) //获取评论操作的类型
	if err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseRelation(c, CodeRelationTypeError)
	}

	// 2.处理关注逻辑
	toUserId := int64(tmpToUser)
	tmpAction := int32(actionType)

	if err := logic.RelationAction(userId, toUserId, tmpAction); err != nil {
		zap.L().Error("logic.RelationAction() failed", zap.Error(err))
		ResponseRelation(c, CodeRelationActionError)
		return
	}

	// 3.返回成功点赞操作的响应
	ResponseRelation(c, CodeSuccess)

}

// FollowRelationHandler  获取用户关注列表
func FollowRelationHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseRelationListError(c, CodeUserIdError)
		return
	}
	// 2. 获取关注用户列表
	followList, err := logic.GetFollowList(userId)
	if err != nil {
		zap.L().Error("logic.GetFollowList()", zap.Error(err))
		ResponseRelationListError(c, CodeRelationFollowError)
		return
	}
	// 3.返回响应
	ResponseRelationListSuccess(c, CodeSuccess, followList)

}

// FanRelationHandler  获取用户粉丝列表
func FanRelationHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseRelationListError(c, CodeUserIdError)
		return
	}
	// 2. 获取关注粉丝列表
	fanList, err := logic.GetFanList(userId)
	if err != nil {
		zap.L().Error("logic.GetFanList()", zap.Error(err))
		ResponseRelationListError(c, CodeRelationFollowError)
		return
	}
	// 3.返回响应
	ResponseRelationListSuccess(c, CodeSuccess, fanList)
}

// FriendRelationHandler  获取用户聊天好友列表
func FriendRelationHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseRelationListError(c, CodeUserIdError)
		return
	}
	// 2. 获取聊天好友列表
	friendList, err := logic.GetFriendList(userId)
	if err != nil {
		zap.L().Error("logic.GetFriendList()", zap.Error(err))
		ResponseRelationListError(c, CodeRelationFollowError)
		return
	}
	// 3.返回响应
	ResponseRelationListSuccess(c, CodeSuccess, friendList)
}
