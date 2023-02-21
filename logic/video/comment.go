package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
	"go.uber.org/zap"
	"time"
)

// AddComment 将用户的评论添加到数据库
func AddComment(userId, videoId int64, commentText string) (commentResponse model.CommentResponse, err error) {
	// 1.组装User结构体,验证videoId的合法性
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err = mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return
	}

	tmpVideo := new(model.Video)
	tmpVideo.Id = videoId
	_, err = mysql.GetVideo(tmpVideo)
	if err != nil {
		zap.L().Error("mysql.GetVideo() failed", zap.Error(err))
		return
	}

	// 2.组装comment结构体，存入MySQL中
	var comment model.Comment
	now := time.Now()
	dateStr := now.Format("01-02")
	unixTime := now.Unix()
	commentId := model.GenID()

	comment.Id = commentId
	comment.UserId = userId
	comment.VideoId = videoId
	comment.Content = commentText
	comment.CreatedAt = unixTime

	if err := mysql.PostComment(comment); err != nil {
		zap.L().Error("mysql.PostComment() failed", zap.Error(err))
		return
	}

	// 3.将comment的关系信息装入Redis
	err = redis.PostComment(videoId, commentId)
	if err != nil {
		zap.L().Error("redis.PostComment() failed", zap.Error(err))
		return
	}

	// 4.组装commentResponse结构体，并进行返回
	commentResponse.Id = commentId
	commentResponse.User = *tmpUser
	commentResponse.Content = commentText
	commentResponse.CreateDate = dateStr

	return commentResponse, nil
}

func RemoveComment(videoId, commentId int64) (commentResponse model.CommentResponse, err error) {
	// 1.验证videoId,commentId的合法性
	tmpVideo := new(model.Video)
	tmpVideo.Id = videoId
	_, err = mysql.GetVideo(tmpVideo)
	if err != nil {
		zap.L().Error("mysql.GetVideo() failed", zap.Error(err))
		return
	}

	tmpComment := new(model.Comment)
	tmpComment.Id = commentId
	_, err = mysql.GetComment(tmpComment)
	if err != nil {
		zap.L().Error("mysql.GetComment() failed", zap.Error(err))
		return
	}

	// 2.从MySQL中获取对应的comment结构体，并进行删除
	if err = mysql.RemoveComment(commentId); err != nil {
		zap.L().Error("mysql.RemoveComment() failed", zap.Error(err))
		return
	}

	// 3.删除Redis中的comment关系信息
	err = redis.RemoveComment(videoId, commentId)
	if err != nil {
		zap.L().Error("redis.RemoveComment() failed", zap.Error(err))
		return
	}

	// 4.将组装成CommentResponse并进行返回
	tmpUser := new(model.User)
	tmpUser.Id = tmpComment.UserId
	_, err = mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return
	}
	tmpTime := time.Unix(tmpComment.CreatedAt, 0) //将unix时间转回 time.Time
	dateStr := tmpTime.Format("01-02")

	commentResponse.Id = commentId
	commentResponse.User = *tmpUser
	commentResponse.CreateDate = dateStr

	return commentResponse, nil

}

func GetCommentList(videoId int64) (commentResponseList *[]model.CommentResponse, err error) {
	// 1.验证videoId的合法性
	tmpVideo := new(model.Video)
	tmpVideo.Id = videoId
	_, err = mysql.GetVideo(tmpVideo)
	if err != nil {
		zap.L().Error("mysql.GetVideo() failed", zap.Error(err))
		return
	}
	// 2.从Redis中获取该videoId对应的所有comment评论id列表，并按照发布时间倒叙
	commentIdList, err := redis.GetCommentIdList(videoId)
	if err != nil {
		zap.L().Error("redis.GetCommentIdList() failed", zap.Error(err))
		return
	}

	// 3.根据commentId列表，组装并获得commentResponse切片
	commentResponseList, err = mysql.GetCommentResponseList(commentIdList)
	if err != nil {
		zap.L().Error("mysql.GetCommentResponseList() failed", zap.Error(err))
		return
	}

	// 4.返回切片
	return commentResponseList, nil

}
