package redis

import (
	"Team2048_Tiktok/model"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
)

// GetFollowStatus 判断是否已经关注当前用户
func GetFollowStatus(userId, toUserId int64) int32 {
	toUserStr := strconv.FormatInt(toUserId, 10)
	userStr := strconv.FormatInt(userId, 10)

	// 判断集合中是否已经存在关注信息
	exists, err := client.SIsMember(model.KeyUserFollowSetPrefix+userStr, toUserStr).Result()
	if err != nil {
		zap.L().Error("client.SIsMember() failed", zap.Error(err))
	}

	if exists {
		// 已经关注了
		return 1
	} else {
		// 没有关注
		return 0
	}
}

// Follow 关注用户
func Follow(userId, toUserId int64) error {
	pipeline := client.TxPipeline()

	// 1.更新被关注者的粉丝集合
	pipeline.SAdd(model.GetRedisKey(model.KeyUserFansSetPrefix+strconv.FormatInt(toUserId, 10)), userId)

	// 2.更新关注者的关注集合
	pipeline.SAdd(model.GetRedisKey(model.KeyUserFollowSetPrefix+strconv.FormatInt(userId, 10)), toUserId)

	// 3.更新关注者的关注数
	pipeline.ZIncrBy(model.KeyUserFollowNumberZset, 1, strconv.FormatInt(userId, 10))

	// 4.更新被关注者的粉丝数
	pipeline.ZIncrBy(model.KeyUserFanNumberZset, 1, strconv.FormatInt(toUserId, 10))

	// 5.执行事务
	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("pipeline.Exec() failed", zap.Error(err))
		return err
	}

	return nil
}

// Follow 取消关注
func Unfollow(userId, toUserId int64) error {
	pipeline := client.TxPipeline()

	// 1.从被关注者的粉丝集合中删除该用户
	pipeline.SRem(model.GetRedisKey(model.KeyUserFansSetPrefix+strconv.FormatInt(toUserId, 10)), userId)

	// 2.从关注者的关注集合中删除被关注者
	pipeline.SRem(model.GetRedisKey(model.KeyUserFollowSetPrefix+strconv.FormatInt(userId, 10)), toUserId)

	// 3.更新关注者的关注数
	pipeline.ZIncrBy(model.KeyUserFollowNumberZset, -1, strconv.FormatInt(userId, 10))

	// 4.更新被关注者的粉丝数
	pipeline.ZIncrBy(model.KeyUserFanNumberZset, -1, strconv.FormatInt(toUserId, 10))

	// 5.执行事务
	_, err := pipeline.Exec()
	if err != nil {
		zap.L().Error("pipeline.Exec() failed", zap.Error(err))
		return err
	}

	return nil
}

// CheckFriend 查看对方有没有互关自己，如果有则加入好友列表
func CheckFriend(userId, toUserId int64) error {

	// 1.查看自己是否在对方的关注表
	//KeyUserFollowSetPrefix  = "user:follow:"
	followKey := model.GetRedisKey(model.KeyUserFollowSetPrefix + strconv.FormatInt(toUserId, 10))
	isFriend, err := client.SIsMember(followKey, userId).Result()
	if err != nil {
		return err
	}
	if !isFriend {
		return nil
	}

	// 2.将双方的Id，加入各自的好友列表键中
	pipeline := client.TxPipeline()
	pipeline.SAdd(model.GetRedisKey(model.KeyFriendshipSetPrefix+strconv.FormatInt(userId, 10)), toUserId)
	pipeline.SAdd(model.GetRedisKey(model.KeyFriendshipSetPrefix+strconv.FormatInt(toUserId, 10)), userId)
	_, err = pipeline.Exec()
	if err != nil {
		return err
	}

	return nil
}

// CheckFriendNot 反向查看自己是否在对方的关注表中，如果不在，则将彼此移除好友列表
func CheckNotFriend(userId, toUserId int64) error {

	// 1.查看自己是否在对方的关注表
	// KeyUserFollowSetPrefix  = "user:follow:"
	followKey := model.GetRedisKey(model.KeyUserFollowSetPrefix + strconv.FormatInt(toUserId, 10))
	isFriend, err := client.SIsMember(followKey, userId).Result()
	if err != nil {
		return err
	}
	if isFriend {
		return nil
	}

	// 2.将双方的Id，从各自的好友列表键中移除
	pipeline := client.TxPipeline()
	pipeline.SRem(model.GetRedisKey(model.KeyFriendshipSetPrefix+strconv.FormatInt(userId, 10)), toUserId)
	pipeline.SRem(model.GetRedisKey(model.KeyFriendshipSetPrefix+strconv.FormatInt(toUserId, 10)), userId)
	_, err = pipeline.Exec()
	if err != nil {
		return err
	}

	return nil
}

// GetFriendDetail 补全好友列表信息
func GetFriendDetail(user model.User) (userResponse model.FriendResponse, err error) {

	// 1.获取用户关注数
	followKey := model.GetRedisKey(model.KeyUserFollowSetPrefix + strconv.FormatInt(user.Id, 10))
	followCount, err := client.SCard(followKey).Result()
	if err != nil {
		zap.L().Error("failed to get follow count from redis", zap.Error(err))
		return
	}

	// 2.获取用户粉丝数
	fansKey := model.GetRedisKey(model.KeyUserFansSetPrefix + strconv.FormatInt(user.Id, 10))
	fansCount, err := client.SCard(fansKey).Result()
	if err != nil {
		zap.L().Error("failed to get fans count from redis", zap.Error(err))
		return
	}

	// 3.获取获赞总数
	likedKey := model.GetRedisKey(model.KeyUserLikedNumZset)
	likedCount, err := client.ZScore(likedKey, strconv.FormatInt(user.Id, 10)).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("failed to get liked count from redis", zap.Error(err))
		return
	}

	// 4.获取发布的视频总数
	publishedKey := model.GetRedisKey(model.KeyUserPublisNumZset)
	publishedCount, err := client.ZScore(publishedKey, strconv.FormatInt(user.Id, 10)).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("failed to get published count from redis", zap.Error(err))
		return
	}

	// 5.点赞总数
	favorKey := model.GetRedisKey(model.KeyUserFavorZsetPrefix + strconv.FormatInt(user.Id, 10))
	favorCount, err := client.ZCard(favorKey).Result()
	if err != nil {
		zap.L().Error("failed to get favor count from redis", zap.Error(err))
		return
	}

	userResponse = model.FriendResponse{
		Id:             user.Id,
		Name:           user.Name,
		FollowCount:    followCount,
		FollowerCount:  fansCount,
		TotalFavorited: int64(likedCount),
		WorkCount:      int64(publishedCount),
		FavoriteCount:  favorCount,
	}

	return userResponse, nil
}

// GetUserDetail 获取完整的用户信息(关注关系除外）
func GetUserDetail(user model.User) (userResponse model.UserResponse, err error) {

	// 1.获取用户关注数
	followKey := model.GetRedisKey(model.KeyUserFollowSetPrefix + strconv.FormatInt(user.Id, 10))
	followCount, err := client.SCard(followKey).Result()
	if err != nil {
		zap.L().Error("failed to get follow count from redis", zap.Error(err))
		return
	}

	// 2.获取用户粉丝数
	fansKey := model.GetRedisKey(model.KeyUserFansSetPrefix + strconv.FormatInt(user.Id, 10))
	fansCount, err := client.SCard(fansKey).Result()
	if err != nil {
		zap.L().Error("failed to get fans count from redis", zap.Error(err))
		return
	}

	// 3.获取获赞总数
	likedKey := model.GetRedisKey(model.KeyUserLikedNumZset)
	likedCount, err := client.ZScore(likedKey, strconv.FormatInt(user.Id, 10)).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("failed to get liked count from redis", zap.Error(err))
		return
	}

	// 4.获取发布的视频总数
	publishedKey := model.GetRedisKey(model.KeyUserPublisNumZset)
	publishedCount, err := client.ZScore(publishedKey, strconv.FormatInt(user.Id, 10)).Result()
	if err != nil && err != redis.Nil {
		zap.L().Error("failed to get published count from redis", zap.Error(err))
		return
	}

	// 5.点赞总数
	favorKey := model.GetRedisKey(model.KeyUserFavorZsetPrefix + strconv.FormatInt(user.Id, 10))
	favorCount, err := client.ZCard(favorKey).Result()
	if err != nil {
		zap.L().Error("failed to get favor count from redis", zap.Error(err))
		return
	}

	userResponse = model.UserResponse{
		Id:             user.Id,
		Name:           user.Name,
		FollowCount:    followCount,
		FollowerCount:  fansCount,
		TotalFavorited: int64(likedCount),
		WorkCount:      int64(publishedCount),
		FavoriteCount:  favorCount,
	}

	return userResponse, nil
}

// GetFollowIdList 获取某个用户的所有关注者Id
func GetFollowIdList(userId int64) ([]int64, error) {

	// 1.获取关注者Set的所有值
	redisKey := model.GetRedisKey(model.KeyUserFollowSetPrefix + strconv.FormatInt(userId, 10))
	members, err := client.SMembers(redisKey).Result()
	if err != nil {
		zap.L().Error("client.SMembers() failed", zap.Error(err))
		return nil, err
	}

	// 2.将字符串切片调整为[]int64
	idList := make([]int64, len(members))
	for i, member := range members {
		id, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt() failed", zap.Error(err))
			return nil, err
		}
		idList[i] = id
	}

	return idList, nil
}

// GetFollowListDetail  获得完整的关注者列表
func GetFollowListDetail(userList []model.User) (*[]model.UserResponse, error) {
	var userResponseList []model.UserResponse

	for _, user := range userList {
		userResponse, err := GetUserDetail(user)
		if err != nil {
			zap.L().Error("GetUserDetail() failed", zap.Error(err))
			return nil, err
		}
		userResponse.IsFollow = true //关注列表的人，全部标记关注状态
		userResponseList = append(userResponseList, userResponse)
	}
	return &userResponseList, nil
}

// GetFanIdList 获取某个用户的所有粉丝Id
func GetFanIdList(userId int64) ([]int64, error) {

	// 1.获取粉丝Set的所有值
	redisKey := model.GetRedisKey(model.KeyUserFansSetPrefix + strconv.FormatInt(userId, 10))
	members, err := client.SMembers(redisKey).Result()
	if err != nil {
		return nil, err
	}

	// 2.将字符串切片调整为[]int64
	idList := make([]int64, len(members))
	for i, member := range members {
		id, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt() failed", zap.Error(err))
			return nil, err
		}
		idList[i] = id
	}

	return idList, nil
}

// GetFanIdListDetail 获得完整的粉丝列表
func GetFanListDetail(userList []model.User, userId int64) (*[]model.UserResponse, error) {
	var userResponseList []model.UserResponse

	for _, user := range userList {
		userResponse, err := GetUserDetail(user)
		if err != nil {
			zap.L().Error("GetUserDetail() failed", zap.Error(err))
			return nil, err
		}
		//查看用户是否有回关自己的粉丝
		record := GetFollowStatus(userId, userResponse.Id)
		if record == 1 {
			userResponse.IsFollow = true
		} else {
			userResponse.IsFollow = false
		}
		userResponseList = append(userResponseList, userResponse)
	}
	return &userResponseList, nil

}

// GetFriendIdList 获取某个用户的所有聊天好友Id
func GetFriendIdList(userId int64) ([]int64, error) {
	// 1.获取聊天好友Set的所有值
	redisKey := model.GetRedisKey(model.KeyFriendshipSetPrefix + strconv.FormatInt(userId, 10))
	members, err := client.SMembers(redisKey).Result()
	if err != nil {
		zap.L().Error("client.SMembers() failed", zap.Error(err))
		return nil, err
	}

	fmt.Println("THis is MEMBER listRRRRRRRRRRRRRRRRR", members)

	// 2.将字符串切片调整为[]int64
	idList := make([]int64, len(members))
	for i, member := range members {
		id, err := strconv.ParseInt(member, 10, 64)
		if err != nil {
			zap.L().Error("strconv.ParseInt() failed", zap.Error(err))
			return nil, err
		}
		idList[i] = id
	}

	return idList, nil
}
