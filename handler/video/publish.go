package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary 视频投稿
// @Description 用户投稿视频
// @Tags Video
// @Accept multipart/form-data
// @Produce json
// @Param token header string true "用户token"
// @Param video formData file true "视频文件"
// @Param cover formData file true "视频封面"
// @Param title formData string true "视频标题"
// @Param desc formData string false "视频描述"
// @Param category_id formData int true "视频分类ID"
// @Success 200 {string} json "{"code":200,"msg":"OK","data":null}"
// @Failure 400 {string} json "{"code":400,"msg":"请求参数错误","data":null}"
// @Failure 401 {string} json "{"code":401,"msg":"未登录或Token过期","data":null}"
// @Failure 500 {string} json "{"code":500,"msg":"服务器内部错误","data":null}"
// @Router /douyin/publish/action [post]
func VideoPublishHandler(c *gin.Context) {

	// 1. 获取请求中的参数和视频数据

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
