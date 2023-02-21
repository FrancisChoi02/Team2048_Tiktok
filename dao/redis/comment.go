package redis

import (
	"Team2048_Tiktok/model"
	"github.com/go-redis/redis"
	"strconv"
)

// PostComment  记录评论的相关信息
func PostComment(videoId, commentId int64) error {
	// 记录评论的关系，以及该视频的评论总数+1
	pipeline := client.TxPipeline()

	videoStr := strconv.Itoa(int(videoId))

	//将video的comment数+1
	pipeline.ZIncrBy(model.GetRedisKey(model.KeyVideoCommentNumZset), 1, videoStr)

	//记录评论Id和评论的时间
	createAt := float64(model.GenID())
	pipeline.ZAdd(model.GetRedisKey(model.KeyVideoCommentZsetPrefix)+videoStr, redis.Z{
		Score:  createAt,
		Member: commentId,
	})
	_, err := pipeline.Exec()

	return err
}

// RemoveComment  移除评论的相关信息
func RemoveComment(videoId, commentId int64) error {
	// 修改评论的关系，以及该视频的评论总数-1
	pipeline := client.TxPipeline()
	videoStr := strconv.Itoa(int(videoId))
	//将video的comment数-1
	pipeline.ZIncrBy(model.GetRedisKey(model.KeyVideoCommentNumZset), -1, videoStr)
	//删除评论Id和评论的时间
	pipeline.ZRem(model.GetRedisKey(model.KeyVideoCommentZsetPrefix)+videoStr, commentId)
	_, err := pipeline.Exec()

	return err
}

// GetCommentIdList 获得评论Id切片
func GetCommentIdList(videoId int64) ([]int64, error) {
	// 获取Zset所有的member，并且按照发布时间的倒序
	videoStr := strconv.Itoa(int(videoId))

	commentIds, err := client.ZRevRange(model.GetRedisKey(model.KeyVideoCommentZsetPrefix)+videoStr, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, len(commentIds))
	for _, idStr := range commentIds {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}