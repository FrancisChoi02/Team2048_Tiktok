package redis

import (
	"fmt"
	"go.uber.org/zap"
)

// KeyExisted 判断键是否已经存在于数据库
func KeyExisted(key string) (exist int64, err error) {
	// 判断键是否存在
	exists, err := client.Exists(key).Result()
	if err != nil {
		zap.L().Error(fmt.Sprintf("Key %s is already existed", key), zap.Error(err))
		return exists, err
	}
	return exists, nil
}
