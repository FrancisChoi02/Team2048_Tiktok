package redis

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/model"
	"errors"
	"go.uber.org/zap"
	"strconv"
)

//GetVideoListDetail 返回用户投稿视频列表的完整数据
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
		videoListFull[i].Id = video.Id
		videoListFull[i].Author = *tmpUser
		videoListFull[i].PlayUrl = video.PlayUrl
		videoListFull[i].CoverUrl = video.CoverUrl
		videoListFull[i].Title = video.Title
		videoListFull[i].CreatedAt = video.CreatedAt
	}
	return &videoListFull, nil
}
