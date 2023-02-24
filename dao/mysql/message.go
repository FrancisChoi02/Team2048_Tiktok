package mysql

import (
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"strconv"
)

// SaveMessage 将消息结构体保存到数据库中
func SaveMessage(message *model.Message) error {
	return DB.Create(message).Error
}

// GetMessage 获得消息结构体
func GetMessage(message *model.Message) (boolstring bool, err error) {
	boolstring = false
	if err := DB.Where("id = ?", message.Id).First(message).Error; err != nil { //这里曾经是&message
		if gorm.IsRecordNotFoundError(err) {
			// 处理记录不存在错误
			zap.L().Error("Message doesn't exist", zap.Error(err))
		} else {
			// 处理其他错误
			zap.L().Error("DB.Where(\"id = ?\", messageId).First(user) failed", zap.Error(err))
		}
		return boolstring, ErrorUserNotExist
	}

	boolstring = true
	return boolstring, err
}

// GetResentMessage  补全Friend结构体的message部分
func GetResentMessage(userId, toUserId int64) (string, int64, error) {

	// 1.从键值对获得两个用户时间戳最新的messageId
	messageIds, err := redis.GetResentMessageId(userId, toUserId)
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
	if _, err := GetMessage(message); err != nil { //获取最新的message的详细信息
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
