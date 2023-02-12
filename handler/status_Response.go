package handler

import (
	"Team2048_Tiktok/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResponseError 返回错误消息
func ResponseError(c *gin.Context, code StatusCode) {
	res := &model.CommonResponse{
		StatusCode: int32(code),
		StatusMsg:  code.Msg(),
	}

	c.JSON(http.StatusOK, res)
}
