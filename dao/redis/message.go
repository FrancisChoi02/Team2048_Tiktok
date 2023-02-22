package redis

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/model"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
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

	// 1.将收发双方的Id，加入各自的好友列表键中
	pipeline.SAdd(model.GetRedisKey(model.KeyFriendshipSetPrefix+strconv.FormatInt(userId, 10)), toUserId)
	pipeline.SAdd(model.GetRedisKey(model.KeyFriendshipSetPrefix+strconv.FormatInt(toUserId, 10)), userId)

	// 2.记录消息的时间戳与消息Id
	chatKey := getChatKey(userId, toUserId)
	pipeline.ZAdd(chatKey, redis.Z{Score: float64(unixTime), Member: messageId}).Result()

	_, err := pipeline.Exec()

	return err
}

// GetResentMessage  补全Friend结构体的message部分
func GetResentMessage(userId, toUserId int64) (string, int64, error) {

	// 1.从键值对获得两个用户时间戳最新的messageId
	chatKey := getChatKey(userId, toUserId)
	messageIds, err := client.ZRevRange(chatKey, 0, 1).Result()
	if err != nil {
		zap.L().Error("client.ZRevRange() failed", zap.Error(err))
		return "", 0, err
	}
	if len(messageIds) == 0 {
		// 如果键不存在或者没有数据，直接返回空结果
		return "", 0, nil
	}
	latestMessageId, err := strconv.ParseInt(messageIds[0], 10, 64)
	if err != nil {
		zap.L().Error("failed to parse latest message id", zap.Error(err))
		return "", 0, err
	}

	// 2.从MySQL中获取对应的Message结构体
	message := new(model.Message)
	message.Id = latestMessageId
	if _, err := mysql.GetMessage(message); err != nil { //获取最新的message的详细信息
		zap.L().Error("GetResentMessage() failed", zap.Error(err))
		return "", 0, err
	}

	// 3.赋值返回消息内容content和消息关系msgType
	var msgType int64
	var content string
	if message.FromUserId != userId {
		msgType = 2
	} else {
		msgType = 1
	}
	content = message.Content

	return content, msgType, nil
}

// GetMessageIdList 获取消息Id列表
func GetMessageIdList(userId, toUserId int64) ([]string, error) {
	chatKey := getChatKey(userId, toUserId)
	messageIds, err := client.ZRange(chatKey, 0, -1).Result()
	return messageIds, err
}
