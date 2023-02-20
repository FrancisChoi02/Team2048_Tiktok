package mysql

import (
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
	"go.uber.org/zap"
)

func GetUserFavorListFull(favorList []int64) (*[]model.VideoResponse, error) {
	lenth := len(favorList)
	favorListFull := make([]model.VideoResponse, lenth)

	for index, videoId := range favorList {
		tmpVideo := new(model.Video)
		tmpVideo.Id = videoId
		_, err := GetVideo(tmpVideo)
		if err != nil {
			zap.L().Error("mysql.GetVideo() failed", zap.Error(err))
			return nil, err
		}

		tmpVideoFull, err := redis.GetVideoDetail(tmpVideo)
		if err != nil {
			zap.L().Error("redis.GetVideoDetail() failed", zap.Error(err))
			return nil, err
		}

		//将完整的视频填充入结构体中
		favorListFull[index] = *tmpVideoFull
	}

	return &favorListFull, nil
}
