package mysql

import (
	"Team2048_Tiktok/model"
	"fmt"
	"github.com/jinzhu/gorm"
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

// GetVideoList 获取用户Id对应的投稿视频列表(视频基础信息)
func GetVideoList(userId int64) (*[]model.Video, error) {
	videoList := new([]model.Video)

	//设定15条视频为上限
	err := DB.Where("user_id=?", userId).
		Select([]string{"id", "user_id", "play_url", "cover_url", "title", "create_time"}).
		Limit(15).
		Find(videoList).Error

	return videoList, err
}

// GetFeedList 获取视频流列表（视频基础信息）
func GetFeedList(latestTime int64) (*[]model.Video, error) {
	videoList := new([]model.Video)

	//按照视频投稿时间的 倒序，latestTime之后的15个视频
	err := DB.Where("created_at<?", latestTime).
		Order("created_at ASC").
		Select([]string{"id", "user_id", "play_url", "cover_url", "title", "create_time"}).
		Limit(15).
		Find(videoList).Error

	return videoList, err
}

// GetVideo 获取视频信息，并查询视频是否存在
func GetVideo(video *model.Video) (boolstring bool, err error) {
	boolstring = false
	if err := DB.Where("id = ?", video.Id).First(video).Error; err != nil { //这里曾经是&user
		if gorm.IsRecordNotFoundError(err) {
			// 处理记录不存在错误
			zap.L().Error("Video doesn't exist", zap.Error(err))
		} else {
			// 处理其他错误
			zap.L().Error("DB.Where(\"id = ?\", video.Id).First(video) failed", zap.Error(err))
		}
		return boolstring, ErrorVideoNotExist
	}

	boolstring = true
	return boolstring, err
}
