package mysql

import (
	"Team2048_Tiktok/model"
	"fmt"
	"go.uber.org/zap"
	"time"
)

// PostVideo 将视频的相关数据保存到数据库中
func PostVideo(videoId, userId int64, filePath, coverPath, title string) (err error) {
	// 1.初始化一个Video结构体
	user := &model.User{}
	user.Id = userId
	fmt.Println(user.Id)
	_, err = GetUser(user)
	if err != nil {
		zap.L().Error(" GetUser() failed", zap.Error(err))
	}

	video := &model.Video{
		Id:           videoId,
		UserId:       userId,
		PlayUrl:      filePath,
		CoverUrl:     coverPath,
		CommentCount: 0,
		Title:        title,
		CreatedAt:    time.Now(), //视频的创建时间 设为 保存到数据库的时间
	}

	// 2.添加到对应的表中
	err = DB.Create(video).Error
	// 3.返回错误
	return err
}
