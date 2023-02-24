package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// FavorateActionHandler 用户点赞操作
// @Summary 用户点赞操作
// @Description 用户点赞操作接口
// @Tags 视频相关接口
// @Accept application/json
// @Produce application/json
// @Param video_id query int true "视频 ID"
// @Param action_type query int true "操作类型，1表示点赞，2表示取消点赞"
// @Security ApiKeyAuth
// @Success 200 {string} string "操作成功"
// @Router /douyin/favorite/action/ [post]
func FavorateActionHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	// 获取上下文中保存的user_id

	rawId, _ := c.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		ResponsePublish(c, CodeUserIdError)
		return
	}

	videoId, err := strconv.Atoi(c.Query("video_id"))
	if err != nil {
		zap.L().Error("videoId invalid", zap.Error(err))
		ResponseVideoListError(c, CodeActionError)
		return
	}

	actionType, err := strconv.Atoi(c.Query("action_type"))
	if err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseVideoListError(c, CodeActionError)
		return
	}

	// 2.处理点赞逻辑
	tmpVideo := int64(videoId)
	tmpAction := int32(actionType)
	if err := logic.TakeAction(int64(userId), tmpVideo, tmpAction); err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseVideoListError(c, CodeActionError)
		return
	}

	// 3.返回成功点赞操作的响应
	ResponseAtcion(c, CodeSuccess)
}

// FavorateListHandler 获取用户的喜欢视频列表
// @Summary 获取用户的喜欢视频列表
// @Description 获取用户的喜欢视频列表接口
// @Tags 视频相关接口
// @Accept application/json
// @Produce application/json
// @Param user_id query int true "用户 ID"
// @Success 200 {string} string "操作成功"
// @Failure 400 {string} string "用户 ID 错误"
// @Failure 500 {string} string "服务器忙，请稍后重试"
// @Router /douyin/favorite/list/ [get]
func FavorateListHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据

	rawId := c.Query("user_id")
	userId, ok := strconv.Atoi(rawId)
	if ok != nil {
		zap.L().Error("user_id invalid")
		ResponseVideoListError(c, CodeUserIdError)
		return
	}

	favorList, err := logic.GetFavorListByUserId(int64(userId))
	if err != nil {
		zap.L().Error("logic.GetFavorListByUserId()", zap.Error(err))
		ResponseVideoListError(c, CodeVideoListError)
		return
	}

	ResponseVideoListSuccess(c, CodeSuccess, favorList)
}
