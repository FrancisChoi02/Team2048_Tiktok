package video

import (
	logic "Team2048_Tiktok/logic/video"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

// CommentActionHandler  用户对视频进行评论
func CommentActionHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据

	rawId, _ := c.Get("user_id") // 获取上下文中保存的user_id
	userId, ok := rawId.(int64)
	if !ok {
		ResponseCommentError(c, CodeUserIdError)
		return
	}
	videoId, err := strconv.Atoi(c.PostForm("video_id")) //获取需要评论操作的视频Id
	if err != nil {
		zap.L().Error("videoId invalid", zap.Error(err))
		ResponseCommentError(c, CodeCommentError)
	}

	actionType, err := strconv.Atoi(c.PostForm("action_type")) //获取评论操作的类型
	if err != nil {
		zap.L().Error("ActionType invalid", zap.Error(err))
		ResponseCommentError(c, CodeCommentError)
	}

	tmpVideo := int64(videoId)

	//根据actionType区分评论操作类型
	if actionType == 1 { //添加评论
		commentText := c.PostForm("comment_text")
		comment, err := logic.AddComment(userId, tmpVideo, commentText)
		if err != nil {
			zap.L().Error("logic.AddComment() invalid", zap.Error(err))
			ResponseCommentError(c, CodeCommentError)
		}
		ResponseCommentSuccess(c, CodeSuccess, comment)

	} else if actionType == 2 { //删除评论
		commentId, err := strconv.Atoi(c.PostForm("comment_id"))
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

// CommentListHandler  获取视频的评论列表
func CommentListHandler(c *gin.Context) {
	// 1. 获取请求中的参数和视频数据

	videoId, err := strconv.Atoi(c.PostForm("video_id")) //获取需要评论操作的视频Id
	if err != nil {
		zap.L().Error("videoId invalid", zap.Error(err))
		ResponseCommentListError(c, CodeCommentError)
	}
	tmpVideo := int64(videoId)

	//获得评论列表
	commentList, err := logic.GetCommentList(tmpVideo)

	ResponseCommentListSuccess(c, CodeSuccess, commentList)
}
