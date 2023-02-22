package redis

import (
	"Team2048_Tiktok/model"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// FavoritePositive  点赞
func FavoritePositive(userId, videoId int64) (err error) {
	videoStr := strconv.Itoa(int(videoId))
	userStr := strconv.Itoa(int(userId))

	//用户对视频的点赞状况、视频的点赞数、用户的喜爱列表需要同步更新
	//因此需要使用事务
	pipeline := client.TxPipeline()
	//Zadd, Zset会自动去重，会覆盖原有的记录
	pipeline.ZAdd(model.GetRedisKey(model.KeyVideoLikedZSetPrefix+videoStr), redis.Z{ //用户点赞状况
		Score:  1, // 赞
		Member: userStr,
	})

	timeStamp := float64(model.GenID())
	pipeline.ZAdd(model.GetRedisKey(model.KeyUserFavorZsetPrefix+userStr), redis.Z{ //用户喜爱列表（带时间）
		Score:  timeStamp, //点赞时间
		Member: videoStr,  // 赞
	})

	pipeline.ZIncrBy(model.GetRedisKey(model.KeyVideoScoreZset), 1, videoStr) //视频的点赞数

	pipeline.ZIncrBy(model.GetRedisKey(model.KeyUserLikedNumZset), 1, userStr) //用户的被点赞数

	_, err = pipeline.Exec()
	if err != nil {
		zap.L().Error("Positive pipeline.Exec() failed", zap.Error(err))
		return err
	}

	return nil
}

// FavoriteNegative  取消点赞
func FavoriteNegative(userId, videoId int64) (err error) {
	videoStr := strconv.Itoa(int(videoId))
	userStr := strconv.Itoa(int(userId))

	//用户对视频的点赞状况、视频的点赞数、用户的喜爱列表需要同步更新
	//因此需要使用事务
	pipeline := client.TxPipeline()
	//ZRem
	pipeline.ZRem(model.GetRedisKey(model.KeyVideoLikedZSetPrefix+videoStr), userStr) //用户点赞状况

	pipeline.ZRem(model.GetRedisKey(model.KeyUserFavorZsetPrefix+userStr), videoStr) ////用户喜爱列表（带时间）

	pipeline.ZIncrBy(model.GetRedisKey(model.KeyVideoScoreZset), -1, videoStr) //视频的点赞数

	pipeline.ZIncrBy(model.GetRedisKey(model.KeyUserLikedNumZset), -1, userStr) //用户的被点赞数

	_, err = pipeline.Exec()
	if err != nil {
		zap.L().Error("pipeline.Exec() failed", zap.Error(err))
		return err
	}

	return nil
}

// GetLikedStatus  返回用户对当前视频的点赞状况
func GetLikedStatus(userId, videoId int64) int32 {
	videoStr := strconv.Itoa(int(videoId))
	userStr := strconv.Itoa(int(userId))
	record := client.ZScore(model.GetRedisKey(model.KeyVideoLikedZSetPrefix+videoStr), userStr).Val()

	res := int32(record)
	return res
}

// GetUserFavorList 返回用户喜爱列表的视频Id切片
func GetUserFavorList(userId int64) ([]int64, error) {
	userStr := strconv.Itoa(int(userId))

	// 1.获取该用户的点赞列表切片
	minScore := "0"
	maxScore := strconv.FormatInt(time.Now().Unix(), 10) //当前时间，绝对是最新的时间
	// 根据Score中表示的时间戳从大到小进行排序，返回对应的视频Id string切片
	tmpFavorId, err := client.ZRevRangeByScore(model.GetRedisKey(model.KeyUserFavorZsetPrefix+userStr), redis.ZRangeBy{
		Min: minScore,
		Max: maxScore,
	}).Result()

	if err != nil {
		zap.L().Error("ZRevRangeByScore() failed", zap.Error(err))
		return nil, err
	}

	// 2.将切片转化为int64类型
	var favorList []int64
	for _, id := range tmpFavorId {
		i, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			zap.L().Error("videoId error when parsing", zap.Error(err))
		}
		favorList = append(favorList, i)
	}

	return favorList, nil
}
