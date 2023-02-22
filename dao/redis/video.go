package redis

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/model"
	"errors"
	"go.uber.org/zap"
	"strconv"
)

// assembleVideoListFull  组装videoListFull单元
func assembleVideoListFull(video *model.Video, index int, videoListFull []model.VideoResponse) {
	videoListFull[index].Id = video.Id
	videoListFull[index].PlayUrl = video.PlayUrl
	videoListFull[index].CoverUrl = video.CoverUrl
	videoListFull[index].Title = video.Title
	videoListFull[index].CreatedAt = video.CreatedAt
}

// assembleVideoFull  组装视频
func assembleVideoFull(video *model.Video, videoFull *model.VideoResponse) {
	videoFull.Id = video.Id
	videoFull.PlayUrl = video.PlayUrl
	videoFull.CoverUrl = video.CoverUrl
	videoFull.Title = video.Title
	videoFull.CreatedAt = video.CreatedAt
}

// RecordPublishNum  记录用户的视频投稿数
func RecordPublishNum(userId int64) error {
	userStr := strconv.Itoa(int(userId))
	err := client.ZIncrBy(model.GetRedisKey(model.KeyUserPublisNumZset), 1, userStr).Err() //用户的视频投稿数
	return err
}

// GetVideoDetail 返回视频的完整数据
func GetVideoDetail(tmpVideo *model.Video) (*model.VideoResponse, error) {
	tmpVideoID := strconv.Itoa(int(tmpVideo.Id))
	videoFull := new(model.VideoResponse)
	// a.获取视频的点赞数
	videoFull.FavoriteCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoScoreZset), tmpVideoID).Val())

	// b.获取视频的评论数
	videoFull.CommentCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoCommentNumZset), tmpVideoID).Val())

	// c.点赞设置为有
	videoFull.IsFavorite = true
	tmpUser := &model.User{}
	tmpUser.Id = tmpVideo.UserId

	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error(" GetUser() failed", zap.Error(err))
	}

	// d.组装videoListFull单元
	videoFull.Author = *tmpUser
	assembleVideoFull(tmpVideo, videoFull)

	return videoFull, nil
}

// GetVideoListDetail 返回用户投稿视频列表的完整数据
func GetVideoListDetail(videoList *[]model.Video) (*[]model.VideoResponse, error) {
	// 1.判断视频列表是否为空
	if videoList == nil {
		err := errors.New("videoList has nothing")
		zap.L().Error("videoList has nothing")
		return nil, err
	}

	// 2.构建VideoResponse数组
	size := len(*videoList)
	videoListFull := make([]model.VideoResponse, size)

	//	3.每个视频单独处理
	for i, video := range *videoList {
		tmpVideoID := strconv.Itoa(int(video.Id))

		// a.获取视频的点赞数
		videoListFull[i].FavoriteCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoScoreZset), tmpVideoID).Val())

		// b.获取视频的评论数
		videoListFull[i].CommentCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoCommentNumZset), tmpVideoID).Val())

		// c.查看用户有没有点赞当前自己的视频
		tmpUser := &model.User{}
		tmpUser.Id = video.UserId

		_, err := mysql.GetUser(tmpUser)
		if err != nil {
			zap.L().Error(" GetUser() failed", zap.Error(err))
		}

		//查看当前用户 给 当前视频 的赞记录
		tmpUserID := strconv.Itoa(int(tmpUser.Id))
		ok := client.ZScore(model.GetRedisKey(model.KeyVideoLikedZSetPrefix+tmpVideoID), tmpUserID).Val()
		if ok == 1 {
			videoListFull[i].IsFavorite = true
		} else {
			videoListFull[i].IsFavorite = false
		}

		// d.组装videoListFull单元
		videoListFull[i].Author = *tmpUser
		assembleVideoListFull(&video, i, videoListFull)
	}
	return &videoListFull, nil
}

// GetFeedListWithNoToken  获取未登录用户的视频流
func GetFeedListWithNoToken(feedList *[]model.Video) (*[]model.VideoResponse, error) {
	// 1.判断视频列表是否为空
	if feedList == nil {
		err := errors.New("videoList has nothing")
		zap.L().Error("videoList has nothing")
		return nil, err
	}

	// 2.构建VideoResponse数组
	size := len(*feedList)
	feedListFull := make([]model.VideoResponse, size)

	//	3.每个视频单独处理
	for i, video := range *feedList {
		tmpVideoID := strconv.Itoa(int(video.Id))

		// a.获取视频的点赞数
		feedListFull[i].FavoriteCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoScoreZset), tmpVideoID).Val())

		// b.获取视频的评论数
		feedListFull[i].CommentCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoCommentNumZset), tmpVideoID).Val())

		// c.用户没登录，默认视频没有点赞
		feedListFull[i].IsFavorite = false

		// d.组装videoListFull单元
		assembleVideoListFull(&video, i, feedListFull)
	}
	return &feedListFull, nil

}

// GetFeedListWithToken  获取登录用户的视频流
func GetFeedListWithToken(userId int64, feedList *[]model.Video) (*[]model.VideoResponse, error) {
	// 1.判断视频列表是否为空
	if feedList == nil {
		err := errors.New("videoList has nothing")
		zap.L().Error("videoList has nothing")
		return nil, err
	}

	// 2.构建VideoResponse数组
	size := len(*feedList)
	feedListFull := make([]model.VideoResponse, size)

	//	3.每个视频单独处理
	for i, video := range *feedList {
		tmpVideoID := strconv.Itoa(int(video.Id))

		// a.获取视频的点赞数
		feedListFull[i].FavoriteCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoScoreZset), tmpVideoID).Val())

		// b.获取视频的评论数
		feedListFull[i].CommentCount = int64(client.ZScore(model.GetRedisKey(model.KeyVideoCommentNumZset), tmpVideoID).Val())

		// c.查看登录的用户有没有点赞当前视频
		tmpUser := &model.User{}
		tmpUser.Id = userId

		_, err := mysql.GetUser(tmpUser)
		if err != nil {
			zap.L().Error(" GetUser() failed", zap.Error(err))
		}

		//查看当前用户 给 当前视频 的赞记录
		tmpUserID := strconv.Itoa(int(tmpUser.Id))
		ok := client.ZScore(model.GetRedisKey(model.KeyVideoLikedZSetPrefix+tmpVideoID), tmpUserID).Val()
		if ok == 1 {
			feedListFull[i].IsFavorite = true
		} else {
			feedListFull[i].IsFavorite = false
		}

		// d.组装videoListFull单元
		feedListFull[i].Author = *tmpUser
		assembleVideoListFull(&video, i, feedListFull)
	}
	return &feedListFull, nil
}
