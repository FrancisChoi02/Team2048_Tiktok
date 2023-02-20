package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VideoPublishHandler 视频投稿
func VideoPublishHandler(c *gin.Context) {

	// 1. 获取请求中的参数和视频数据
	// 获取上下文中保存的user_id
	rawId, _ := c.Get("user_id")
	userId, ok := rawId.(int64)
	if !ok {
		ResponsePublish(c, CodeUserIdError)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		ResponsePublish(c, CodeVideoError)
		return
	}

	// 2. 保存视频以及视频信息,如有错误返回响应
	err = logic.VideoPublish(c, userId, form)
	if err != nil {
		zap.L().Error("logic.VideoPublish() failed", zap.Error(err))
		ResponsePublish(c, CodePublishFailed)
	}
}
