package redis

import (
	"Team2048_Tiktok/conf"
	"fmt"
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

// Init 初始化连接
func Init(cfg *conf.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password, // no password set
		DB:           cfg.DB,       // use default DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}

	return nil
}

//	Close 关闭Redis数据库
func Close() {
	_ = client.Close()
}
