package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// FavorateActionHandler  用户点赞操作
func FavorateActionHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	// 获取上下文中保存的user_id
	rawId, _ := c.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		ResponsePublish(c, CodeUserIdError)
		return
	}

	videoId, err := strconv.Atoi(c.PostForm("video_id"))
	if err != nil {
		zap.L().Error("videoId invalid", zap.Error(err))
		ResponseVideoListError(c, CodeActionError)
	}

	actionType, err := strconv.Atoi(c.PostForm("action_type"))
	if err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseVideoListError(c, CodeActionError)
	}

	// 2.处理点赞逻辑
	tmpVideo := int64(videoId)
	tmpAction := int32(actionType)
	if err := logic.TakeAction(userId, tmpVideo, tmpAction); err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseVideoListError(c, CodeActionError)
	}

	// 3.返回成功点赞操作的响应
	ResponseAtcion(c, CodeSuccess)
}

// FavorateListHandler 获取用户的喜欢视频列表
func FavorateListHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		zap.L().Error("user_id invalid")
		ResponseVideoListError(c, CodeUserIdError)
		return
	}

	favorList, err := logic.GetFavorListByUserId(userId)
	if err != nil {
		zap.L().Error("logic.GetFavorListByUserId()", zap.Error(err))
		ResponseVideoListError(c, CodeVideoListError)
	}

	ResponseVideoListSuccess(c, CodeSuccess, favorList)
}
