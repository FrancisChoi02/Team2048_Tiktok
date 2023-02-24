package redis

import (
	"Team2048_Tiktok/model"
	"github.com/go-redis/redis"
	"strconv"
)

// getChatKey 根据用户id和好友id获取Redis键
func getChatKey(userId, friendId int64) string {
	return model.KeyUserMessageZsetPrefix + strconv.FormatInt(Min(userId, friendId), 10) + "_" + strconv.FormatInt(Max(userId, friendId), 10)
}

// Min 获取两个int64类型的值中的最小值
func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Max 获取两个int64类型的值中的最大值
func Max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// SaveMessageRelation  保存消息关系
func SaveMessageRelation(userId, toUserId, messageId, unixTime int64) error {
	pipeline := client.TxPipeline()

	// 记录消息的时间戳与消息Id
	chatKey := getChatKey(userId, toUserId)
	pipeline.ZAdd(chatKey, redis.Z{Score: float64(unixTime), Member: messageId}).Result()

	_, err := pipeline.Exec()

	return err
}

// GetMessageIdList 获取最新的消息Id列表
func GetResentMessageId(userId, toUserId int64) ([]string, error) {
	chatKey := getChatKey(userId, toUserId)
	return client.ZRevRange(chatKey, 0, 1).Result()
}

// GetMessageIdList 获取消息Id列表
func GetMessageIdList(userId, toUserId int64) ([]string, error) {
	chatKey := getChatKey(userId, toUserId)
	messageIds, err := client.ZRange(chatKey, 0, -1).Result()
	return messageIds, err
}
