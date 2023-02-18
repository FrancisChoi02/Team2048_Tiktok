package mysql

import (
	"Team2048_Tiktok/model"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// PostVideo 将视频的相关数据保存到数据库中
func PostVideo(videoId, userId int64, filePath, coverPath, title string) (err error) {
	// 1.初始化一个User结构体
	user := &model.User{}
	user.Id = userId
	fmt.Println(user.Id)
	_, err = GetUser(user)
	if err != nil {
		zap.L().Error(" GetUser() failed", zap.Error(err))
	}

	// 2.初始化一个Video结构体
	video := &model.Video{
		Id:        videoId,
		UserId:    userId,
		PlayUrl:   filePath,
		CoverUrl:  coverPath,
		Title:     title,
		CreatedAt: time.Now().Unix(), //视频的创建时间 设为 保存到数据库的时间
	}

	// 2.添加到对应的表中
	err = DB.Create(video).Error
	// 3.返回错误
	return err
}

// GetVideoList 获取用户Id对应的投稿视频列表
func GetVideoList(userId int64) (*[]model.Video, error) {
	videoList := new([]model.Video)

	//设定15条视频为上限
	err := DB.Where("user_id=?", userId).
		Select([]string{"id", "user_id", "play_url", "cover_url", "title", "create_time"}).
		Limit(15).
		Find(videoList).Error

	return videoList, err
}
