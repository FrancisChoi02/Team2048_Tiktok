package user

import (
	logic "Team2048_Tiktok/logic/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// RelationshipHandler  用户关系操作
// @Summary 用户关系操作
// @Description 用户关系操作接口
// @Tags 关注相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param to_user_id query int true "接收方用户 ID"
// @Param action_type query int true "操作类型，1表示关注，2表示取消关注"
// @Security ApiKeyAuth
// @Success 200 {string} string "操作成功"
// @Router /douyin/relation/action/ [post]
func RelationshipHandler(c *gin.Context) {

	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseRelation(c, CodeUserIdError)
		return
	}
	rawToUser := c.Query("to_user_id")
	tmpToUser, err := strconv.Atoi(rawToUser)
	if err != nil {
		ResponseRelation(c, CodeUserIdError)
		return
	}

	actionType, err := strconv.Atoi(c.Query("action_type")) //获取评论操作的类型
	if err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseRelation(c, CodeRelationTypeError)
		return
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
// @Summary 获取用户关注列表
// @Description 获取用户关注列表接口
// @Tags 关注相关接口
// @Accept application/json
// @Produce application/json
// @Param user_id query int true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.FollowRelationResponse"关注列表"
// @Router /douyin/relation/follower/list/ [get]
func FollowRelationHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据

	rawUser := c.Query("user_id")
	tmpUser, err := strconv.Atoi(rawUser)
	if err != nil {
		ResponseRelation(c, CodeUserIdError)
		return
	}
	userId := int64(tmpUser)

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
// @Summary 获取用户粉丝列表
// @Description 获取用户粉丝列表接口
// @Tags 关注相关接口
// @Accept application/json
// @Produce application/json
// @Param user_id query int true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.FollowRelationResponse "关注粉丝列表"
// @Router /relation/fans [get]
func FanRelationHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据

	rawUser := c.Query("user_id")
	tmpUser, err := strconv.Atoi(rawUser)
	if err != nil {
		ResponseRelation(c, CodeUserIdError)
		return
	}
	userId := int64(tmpUser)

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

// @Description 获取用户聊天好友列表接口
// @Tags 好友相关接口
// @Accept application/json
// @Produce application/json
// @Param user_id query int true "用户 ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.FriendListResponse "好友列表"
// @Router /douyin/relation/action/friend/list/ [get]
func FriendRelationHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据

	rawUser := c.Query("user_id")
	tmpUser, err := strconv.Atoi(rawUser)
	if err != nil {
		ResponseRelation(c, CodeUserIdError)
		return
	}
	userId := int64(tmpUser)

	// 2. 获取聊天好友列表
	friendList, err := logic.GetFriendList(userId)
	if err != nil {
		zap.L().Error("logic.GetFriendList()", zap.Error(err))
		ResponseFriendListError(c, CodeRelationFollowError)
		return
	}

	// 3.返回响应
	ResponseFriendListSuccess(c, CodeSuccess, friendList)
}
