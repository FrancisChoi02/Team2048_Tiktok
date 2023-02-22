package mysql

import (
	"Team2048_Tiktok/model"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

// SaveMessage 将消息结构体保存到数据库中
func SaveMessage(message *model.Message) error {
	return DB.Create(message).Error
}

// GetMessage
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
