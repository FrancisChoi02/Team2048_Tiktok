package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
	"fmt"
	"go.uber.org/zap"
)

// TakeAction 点赞逻辑
func TakeAction(userId, videoId int64, actionType int32) error {

	// 1.查看userId 和 videoId是否存在
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return err
	}

	tmpVideo := new(model.Video)
	tmpVideo.Id = videoId
	_, err = mysql.GetVideo(tmpVideo)
	if err != nil {
		zap.L().Error("mysql.GetVideo() failed", zap.Error(err))
		return err
	}

	// 2.查看 user 对该视频的点赞情况(1点赞、2取消点赞）
	recordLiked := redis.GetLikedStatus(userId, videoId)

	if actionType == recordLiked { //不允许重复投票
		zap.L().Error("Double voted is not allowed", zap.Error(err))
		return err
	}

	// 3.使用事务，更新视频点赞情况、视频点赞数、用户喜爱列表
	if actionType == 1 {
		//点赞，更新投票相关键值对
		if err := redis.FavoritePositive(userId, videoId); err != nil {
			zap.L().Error("redis.FavoritePositive() failed", zap.Error(err))
			return err
		}
	} else if actionType == 2 {
		//取消赞，更新投票相关键值对
		if err := redis.FavoriteNegative(userId, videoId); err != nil {
			zap.L().Error("redis.FavoriteNegative() failed", zap.Error(err))
			return err
		}
	} else {
		zap.L().Error("Invalid actionType", zap.Error(err))
		return err
	}

	return nil
}

func GetFavorListByUserId(userId int64) (*[]model.VideoResponse, error) {
	// 1.查看UserId是否存在
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return nil, err
	}

	// 2.根据UserId，从Redis中获取User的喜欢视频的Id列表
	favorList, err := redis.GetUserFavorList(userId)

	fmt.Println("This is user's Favorite  !!!!!!!!!", favorList)

	if err != nil {
		zap.L().Error("redis.GetUserFavorList() failed", zap.Error(err))
		return nil, err
	}

	// 3.根据视频Id列表，从MySQL中获取视频的详细信息，返回视频列表
	favorListFull, err := GetUserFavorListFull(favorList)
	if err != nil {
		zap.L().Error("mysql.GetUserFavorListFull() failed", zap.Error(err))
		return nil, err
	}

	fmt.Println("This is user's FuLL  !!!!!!!!!", favorListFull)

	return favorListFull, nil
}

// GetUserFavorListFull 获得用户喜爱视频的完整数据
func GetUserFavorListFull(favorList []int64) (*[]model.VideoResponse, error) {
	lenth := len(favorList)
	favorListFull := make([]model.VideoResponse, lenth)

	for index, videoId := range favorList {
		tmpVideo := new(model.Video)
		tmpVideo.Id = videoId
		_, err := mysql.GetVideo(tmpVideo)
		if err != nil {
			zap.L().Error("mysql.GetVideo() failed", zap.Error(err))
			return nil, err
		}

		tmpVideoFull, err := GetVideoDetail(tmpVideo)
		if err != nil {
			zap.L().Error("redis.GetVideoDetail() failed", zap.Error(err))
			return nil, err
		}

		//将完整的视频填充入结构体中
		favorListFull[index] = *tmpVideoFull
	}

	return &favorListFull, nil
}
