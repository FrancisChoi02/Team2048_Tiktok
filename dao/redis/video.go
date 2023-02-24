package redis

import (
	"Team2048_Tiktok/model"

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

// RecordPublishNum  记录用户的视频投稿数
func RecordPublishNum(userId int64) error {
	userStr := strconv.Itoa(int(userId))
	err := client.ZIncrBy(model.GetRedisKey(model.KeyUserPublisNumZset), 1, userStr).Err() //用户的视频投稿数
	return err
}

// GetFeedListStatus 获取视频流中返回的视频状态
func GetFeedListStatus(userId, videoId string) (int64, int64, float64) {
	commentCount := int64(client.ZScore(model.GetRedisKey(model.KeyVideoCommentNumZset), videoId).Val())
	favoriteCount := int64(client.ZScore(model.GetRedisKey(model.KeyVideoScoreZset), videoId).Val())
	ok := client.ZScore(model.GetRedisKey(model.KeyVideoLikedZSetPrefix+videoId), userId).Val()
	return commentCount, favoriteCount, ok
}
