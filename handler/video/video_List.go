package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// VideoListHandler 获取视频列表
// @Summary 获取指定用户的视频列表
// @Description 根据用户ID获取视频列表
// @Tags Video
// @Accept json
// @Produce json
// @Param user_id query int true "用户ID"
// @Success 200 {object} model.VideoListResponse
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /douyin/publish/list/ [get]
func VideoListHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	//参数位置在Query
	rawId := c.Query("user_id")
	userId, ok := strconv.Atoi(rawId)
	if ok != nil {
		zap.L().Error("user_id invalid")
		ResponseVideoListError(c, CodeUserIdError)
		return
	}

	videoList, err := logic.GetVideoListByUserId(int64(userId))
	if err != nil {
		zap.L().Error("logic.GetVideoListByUserId() failed", zap.Error(err))
		ResponseVideoListError(c, CodeVideoListError)
		return
	}

	ResponseVideoListSuccess(c, CodeSuccess, videoList)
}
