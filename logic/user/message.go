package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// SendMessage  发送消息
func SendMessage(userId, toUserId int64, content string) (err error) {
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

	// 2.组装Message结构体
	now := time.Now().Truncate(time.Hour)
	hourStr := now.Format("2006-01-02 15:00:00")
	unixTime := now.Unix()

	message := new(model.Message)
	message.Id = model.GenID()
	message.FromUserId = userId
	message.ToUserId = toUserId
	message.Content = content
	message.CreatTime = hourStr

	// 3.将Message保存到MySQL
	if err := mysql.SaveMessage(message); err != nil {
		zap.L().Error("mysql.SaveMessage() failed", zap.Error(err))
		return
	}

	// 4.将Message和收发双方的关系保存到Redis
	if err := redis.SaveMessageRelation(userId, toUserId, message.Id, unixTime); err != nil {
		zap.L().Error("mysql.SaveMessage() failed", zap.Error(err))
		return
	}

	return nil
}

// GetMessageList 获取聊天记录列表
func GetMessageList(userId, toUserId int64) (*[]model.Message, error) {
	// 1.根据键值对，获取排序后的MessageId
	messageIds, err := redis.GetMessageIdList(userId, toUserId)
	if err != nil {
		zap.L().Error("redis.GetMessageIdList() failed", zap.Error(err))
		return nil, err
	}

	// 2.获得Message列表
	messageList := make([]model.Message, 0, len(messageIds))
	for _, messageId := range messageIds {

		message := new(model.Message)
		tmpMessageId, _ := strconv.Atoi(messageId)
		message.Id = int64(tmpMessageId)

		if _, err := mysql.GetMessage(message); err != nil {
			zap.L().Error("failed to get message from mysql", zap.Error(err))
			continue
		}
		messageList = append(messageList, *message)
	}

	// 3.返回Message列表
	return &messageList, nil
}


