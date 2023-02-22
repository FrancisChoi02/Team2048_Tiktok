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

// ResponseFeedError 返回视频流错误
func ResponseFeedError(c *gin.Context, code StatusCode) {
	var res model.FeedResponse
	res.Code = int32(code)
	res.Msg = code.Msg()
	res.VideoList = new([]model.VideoResponse)
	res.NextTime = 0
	c.JSON(http.StatusOK, res)
}

// ResponseFeedSuccess 返回视频流响应
func ResponseFeedSuccess(c *gin.Context, code StatusCode, videoList *[]model.VideoResponse, nextTime int64) {
	var res model.FeedResponse
	res.Code = int32(code)
	res.Msg = code.Msg()
	res.VideoList = videoList
	res.NextTime = nextTime
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

// ResponseAtcion 返回点赞操作的情况
func ResponseAtcion(c *gin.Context, code StatusCode) {
	var res model.FavorateActionResponse
	res.Code = int32(code)
	res.Msg = code.Msg()
	c.JSON(http.StatusOK, res)
}

// ResponseComment 返回成功评论的结果
func ResponseCommentSuccess(c *gin.Context, code StatusCode, comment model.CommentResponse) {
	var res model.CommentActionResponse
	res.Code = int32(code)
	res.Msg = code.Msg()
	res.Comment = comment
	c.JSON(http.StatusOK, res)
}

// ResponseCommentErr 返回评论失败的结果
func ResponseCommentError(c *gin.Context, code StatusCode) {
	var res model.CommentActionResponse
	var tmpComment model.CommentResponse

	res.Code = int32(code)
	res.Msg = code.Msg()
	res.Comment = tmpComment //返回一个空的评论
	c.JSON(http.StatusOK, res)
}

// ResponseCommentListError 返回评论失败的结果
func ResponseCommentListError(c *gin.Context, code StatusCode) {
	var res model.CommentListResponse
	var tmpComment *[]model.CommentResponse

	res.Code = int32(code)
	res.Msg = code.Msg()
	res.CommentList = tmpComment //返回一个空的评论数组
	c.JSON(http.StatusOK, res)
}

// ResponseCommentListSuccess 返回评论失败的结果
func ResponseCommentListSuccess(c *gin.Context, code StatusCode, commentList *[]model.CommentResponse) {
	var res model.CommentListResponse
	res.Code = int32(code)
	res.Msg = code.Msg()
	res.CommentList = commentList //返回一个空的评论数组
	c.JSON(http.StatusOK, res)
}
