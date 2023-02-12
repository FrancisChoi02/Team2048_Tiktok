package model

// redis key 要适当使用前缀字段，方便拆分和查询
// 字段的拆分 用 ':'

const (
	KeyPrefix = "Tiktok2048" //统一该项目在Redis数据库中的前缀，方便检索
	//互动接口
	KeyVideoScoreZset      = "video:score"     //视频 以及 累积的点赞数量
	KeyVideoLikedSetPrefix = "video:liked:"    //这是 键 的前缀，搭配 视频ID 成为一个完成的键；值是对该视频按了点赞的用户ID
	KeyUserFavorSetPrefix  = "video:favorite:" //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户点过赞的视频ID

	//社交接口
	KeyUserFansSetPrefix   = "user:fans:"   //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户的粉丝的ID
	KeyUserFollowSetPrefix = "user:follow:" //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户关注的博主的ID
	//好友列表 不一定用redis
	KeyFriendshipSetPrefix = "user:friend:" //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户聊天列表中的好友

)

// getRedisKey 为 Redis key 添加前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
