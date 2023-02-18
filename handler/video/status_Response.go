package video

import (
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponsePublish 返回视频上传结果
func ResponsePublish(c *gin.Context, code StatusCode) {
	var res model.VideoUploadResponse
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}

// ResponseVideoListError 返回用户视频列表获取错误
func ResponseVideoListError(c *gin.Context, code StatusCode) {
	var res model.VideoListResponse

	res.Code = int32(code)
	res.Msg = code.Msg()
	res.VideoList = new([]model.VideoResponse)
	c.JSON(http.StatusOK, res)
}

// ResponseVideoListSuccess 返回用户视频列表获取成功
func ResponseVideoListSuccess(c *gin.Context, code StatusCode, videoList *[]model.VideoResponse) {
	var res model.VideoListResponse

	res.Code = int32(code)
	res.Msg = code.Msg()
	res.VideoList = videoList
	c.JSON(http.StatusOK, res)
}
