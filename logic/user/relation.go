package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
	"go.uber.org/zap"
)

// RelationAction 关系操作逻辑
func RelationAction(userId, toUserId int64, actionType int32) error {
	// 1.查看userId 和 toUseroId是否存在
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed, userId error", zap.Error(err))
		return err
	}
	tmpToUser := new(model.User)
	tmpToUser.Id = userId
	_, err = mysql.GetUser(tmpToUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed, toUserId error", zap.Error(err))
		return err
	}

	// 2.检查是否重复操作
	recordLiked := redis.GetFollowStatus(userId, toUserId)

	if actionType == recordLiked { //不允许重复操作
		zap.L().Error("Double follow is not allowed", zap.Error(err))
		return err
	}
	// 3.更新关注关系和关注、粉丝数
	if actionType == 1 {
		//关注，更新关注相关键值对
		if err := FollowUser(userId, toUserId); err != nil {
			zap.L().Error(" FollowUser() failed", zap.Error(err))
			return err
		}
	} else if actionType == 2 {
		//取消关注，更新关注相关键值对
		if err := UnfollowUser(userId, toUserId); err != nil {
			zap.L().Error("UnfollowUser() failed", zap.Error(err))
			return err
		}
	} else {
		zap.L().Error("Invalid actionType", zap.Error(err))
		return err
	}

	return nil
}

// FollowUser  根据Id关注用户
func FollowUser(userId, toUserId int64) (err error) {
	// 1.检验userId、toUserId是否合法
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err = mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return
	}

	tmpToUser := new(model.User)
	tmpToUser.Id = toUserId
	_, err = mysql.GetUser(tmpToUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return
	}

	// 2.更新 双方关注者、粉丝关系与数量
	if err = redis.Follow(userId, toUserId); err != nil {
		zap.L().Error(" redis.Follow() failed", zap.Error(err))
		return
	}

	return nil
}

// UnfollowUser  根据Id取关用户
func UnfollowUser(userId, toUserId int64) (err error) {
	// 1.检验userId、toUserId是否合法
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err = mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return
	}

	tmpToUser := new(model.User)
	tmpToUser.Id = toUserId
	_, err = mysql.GetUser(tmpToUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return
	}

	// 2.更新 双方关注者、粉丝关系与数量
	if err = redis.Unfollow(userId, toUserId); err != nil {
		zap.L().Error("redis.Unfollow() failed", zap.Error(err))
		return
	}
	return nil
}

// GetFollowList  获取关注列表
func GetFollowList(userId int64) (*[]model.UserResponse, error) {
	// 1.检验userId是否合法
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return nil, err
	}

	followIdList, err := redis.GetFollowIdList(userId)

	if err != nil {
		zap.L().Error("redis.GetFollowList() failed", zap.Error(err))
		return nil, err
	}

	tmpUserList, err := mysql.GetUserList(followIdList)
	if err != nil {
		zap.L().Error("mysql.GetUserList() failed", zap.Error(err))
		return nil, err
	}

	followList, err := redis.GetFollowListDetail(tmpUserList)
	if err != nil {
		zap.L().Error("redis.GetFollowList() failed", zap.Error(err))
		return nil, err
	}

	return followList, nil
}

// GetFanList  获取粉丝列表
func GetFanList(userId int64) (*[]model.UserResponse, error) {
	// 1.检验userId是否合法
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return nil, err
	}

	fanIdList, err := redis.GetFanIdList(userId)
	if err != nil {
		zap.L().Error("redis.GetFollowList() failed", zap.Error(err))
		return nil, err
	}

	tmpUserList, err := mysql.GetUserList(fanIdList)
	if err != nil {
		zap.L().Error("mysql.GetUserList() failed", zap.Error(err))
		return nil, err
	}

	fanList, err := redis.GetFanListDetail(tmpUserList, userId)
	if err != nil {
		zap.L().Error("redis.GetFollowList() failed", zap.Error(err))
		return nil, err
	}

	return fanList, nil

}

// GetFriendList  获取聊天好友列表
func GetFriendList(userId int64) (*[]model.FriendResponse, error) {
	// 1.检验userId是否合法
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return nil, err
	}

	// 2.获取好友列表Id
	friendIdList, err := redis.GetFriendIdList(userId)
	if err != nil {
		zap.L().Error("redis.GetFollowList() failed", zap.Error(err))
		return nil, err
	}

	// 3.获取好友列表基础信息
	tmpUserList, err := mysql.GetUserList(friendIdList)
	if err != nil {
		zap.L().Error("mysql.GetUserList() failed", zap.Error(err))
		return nil, err
	}

	// 4.补全好友列表
	friendList, err := redis.GetFriendListDetail(tmpUserList, userId)
	if err != nil {
		zap.L().Error("redis.GetFollowList() failed", zap.Error(err))
		return nil, err
	}

	return friendList, nil

}
