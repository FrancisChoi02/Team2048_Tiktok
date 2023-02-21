package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func VideoListHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		zap.L().Error("user_id invalid")
		ResponseVideoListError(c, CodeUserIdError)
		return
	}

	videoList, err := logic.GetVideoListByUserId(userId)
	if err != nil {
		zap.L().Error("logic.GetVideoListByUserId() failed", zap.Error(err))
		ResponseVideoListError(c, CodeVideoListError)
		return
	}

	ResponseVideoListSuccess(c, CodeSuccess, videoList)
}
