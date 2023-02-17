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
