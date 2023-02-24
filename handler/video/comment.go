package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CommentActionHandler 用户对视频进行评论
// @Summary 用户对视频进行评论
// @Description 用户对视频进行评论的接口
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param video_id query int true "需要评论操作的视频 Id"
// @Param action_type query int true "操作类型，1 表示添加评论，2 表示删除评论"
// @Param comment_text query string false "评论内容，当 action_type 为 1 时必填"
// @Security ApiKeyAuth
// @Success 200 {object} CommentResponse "操作成功"
// @Failure 400 {string} string "用户 ID 错误"
// @Router /douyin/comment/action/ [post]
func CommentActionHandler(c *gin.Context) {

	// 1. 获取请求中的参数和视频数据
	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseCommentError(c, CodeUserIdError)
		return
	}
	videoId, err := strconv.Atoi(c.Query("video_id")) //获取需要评论操作的视频Id
	if err != nil {
		zap.L().Error("videoId invalid", zap.Error(err))
		ResponseCommentError(c, CodeCommentError)
	}

	actionType, err := strconv.Atoi(c.Query("action_type")) //获取评论操作的类型
	if err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseCommentError(c, CodeCommentError)
	}

	tmpVideo := int64(videoId)

	//根据actionType区分评论操作类型
	if actionType == 1 { //添加评论
		commentText := c.Query("comment_text")
		comment, err := logic.AddComment(userId, tmpVideo, commentText)
		if err != nil {
			zap.L().Error("logic.AddComment() invalid", zap.Error(err))
			ResponseCommentError(c, CodeCommentError)
		}
		ResponseCommentSuccess(c, CodeSuccess, comment)

	} else if actionType == 2 { //删除评论
		commentId, err := strconv.Atoi(c.Query("comment_id"))
		if err != nil {
			zap.L().Error("commentId invalid", zap.Error(err))
			ResponseCommentError(c, CodeCommentError)
		}
		tmpCommentId := int64(commentId)

		comment, err := logic.RemoveComment(tmpVideo, tmpCommentId)
		if err != nil {
			zap.L().Error("logic.RemoveComment() invalid", zap.Error(err))
			ResponseCommentError(c, CodeCommentError)
		}
		ResponseCommentSuccess(c, CodeSuccess, comment)

	} else { //actionType不合法
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseCommentError(c, CodeCommentError)
	}
}

// CommentListHandler 获取视频的评论列表
// @Summary 获取视频的评论列表
// @Description 获取视频的评论列表的接口
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param video_id query int true "视频 ID"
// @Success 200 {object} CommentListResponse "操作成功"

// @Router /douyin/comment/list/ [get]
func CommentListHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据

	videoId, err := strconv.Atoi(c.Query("video_id")) //获取需要评论操作的视频Id
	if err != nil {
		zap.L().Error("videoId invalid", zap.Error(err))
		ResponseCommentListError(c, CodeCommentError)
	}
	tmpVideo := int64(videoId)

	//获得评论列表
	commentList, err := logic.GetCommentList(tmpVideo)

	ResponseCommentListSuccess(c, CodeSuccess, commentList)
}
