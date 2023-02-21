package video

import (
	logic "Team2048_Tiktok/logic/video"
	"Team2048_Tiktok/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// FeedVideoListHandler 视频流
func FeedVideoListHandler(c *gin.Context) {
	// 1.从上下文中获取token 和 latest_time
	latestTime, err := strconv.ParseInt(c.PostForm("latest_time"), 10, 64)
	if err != nil {
		zap.L().Error("latest_time invalid", zap.Error(err))
		latestTime = time.Now().Unix() //如果latestTime为空，则返回当前时间戳
	}

	tokenStr := c.Query("token")
	if tokenStr == "" {
		tokenStr = c.PostForm("token")
	}

	// 2.区分有无token两种情况（有token的，加载点赞情况）
	if tokenStr == "" {
		//使用没有登录的feed
		//返回没有登录的Msg提醒
		videoListFull, err := logic.FeedWithNoToken(latestTime)
		if err != nil {
			ResponseVideoListError(c, CodeVideoListError)
			return
		}
		//返回正常响应
		ResponseVideoListSuccess(c, CodeSuccess, videoListFull)
	} else {
		//验证token合法性
		claims, err := middleware.ParseToken(tokenStr)
		if err != nil {
			//使用没有登录的feed
			//返回没有登录的Msg提醒
			videoListFull, err := logic.FeedWithNoToken(latestTime)
			if err != nil {
				ResponseVideoListError(c, CodeVideoListError)
				return
			}
			//返回正常响应
			ResponseVideoListSuccess(c, CodeSuccess, videoListFull)
			return
		}

		//检查token是否过期
		if time.Now().Unix() > claims.ExpiresAt {
			//使用没有登录的feed
			//返回没有登录的Msg题型
			videoListFull, err := logic.FeedWithNoToken(latestTime)
			if err != nil {
				zap.L().Error("logic.FeedWithNoToken() failed", zap.Error(err))
				ResponseVideoListError(c, CodeVideoListError)
				return
			}
			//返回正常响应
			ResponseVideoListSuccess(c, CodeSuccess, videoListFull)
			return
		}

		//使用登录后的feed，获取user_id
		tmpId := claims.UserId
		videoListFull, err := logic.FeedWithToken(latestTime, tmpId)
		if err != nil {
			zap.L().Error("logic.FeedWithToken() failed", zap.Error(err))
			ResponseVideoListError(c, CodeVideoListError)
			return
		}
		//返回正常响应
		ResponseVideoListSuccess(c, CodeSuccess, videoListFull)
	}

}
