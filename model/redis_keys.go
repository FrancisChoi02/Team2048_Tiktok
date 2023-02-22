package model

// redis key 要适当使用前缀字段，方便拆分和查询
// 字段的拆分 用 ':'

const (
	KeyPrefix = "Tiktok2048:" //统一该项目在Redis数据库中的前缀，方便检索
	//互动接口
	KeyVideoScoreZset       = "video:score"      //视频 以及 累积的点赞数量
	KeyVideoCommentNumZset  = "video:commentNum" //视频 以及 累积的评论数量
	KeyVideoLikedZSetPrefix = "video:liked:"     //这是 键 的前缀，搭配 视频ID 成为一个完成的键；值是对该视频按了用户ID 与 点赞状态
	KeyUserLikedNumZset     = "user:likedNum"    //用户 以及 被点赞的数量
	KeyUserPublisNumZset    = "user:publishNum"  //用户 以及 发布视频的数量
	KeyUserFavorZsetPrefix  = "user:favorite:"   //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户点过赞的视频ID 与 操作时间戳

	//视频评论，以视频id为键，保存所有的评论id，评论本身存在Mysql中
	KeyVideoCommentZsetPrefix = "video:comment:" //这是 键 的前缀，搭配 视频ID 成为一个完成的键；值是评论的Id和评论的发布时间

	//社交接口
	KeyUserFansSetPrefix    = "user:fans:"      //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户的粉丝的ID
	KeyUserFollowSetPrefix  = "user:follow:"    //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户关注的博主的ID
	KeyUserFollowNumberZset = "relation:follow" //用户 以及 关注的用户的数量
	KeyUserFanNumberZset    = "relation:fan"    //用户 以及 用户的粉丝的数量

	KeyFriendshipSetPrefix = "relation:friend:" //这是 键 的前缀，搭配 用户ID 成为一个完整的键；值是该用户聊天列表中的好友
	//好友列表，发一条消息就将彼此加入列表
	KeyUserMessageZsetPrefix = "relation:message:" //好友消息，以聊天双方的ID为键，保存信息的id和时间戳
)

// getRedisKey 为 Redis key 添加前缀
func GetRedisKey(key string) string {
	return KeyPrefix + key
}
