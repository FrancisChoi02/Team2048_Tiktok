package middleware

import (
	"Team2048_Tiktok/handler"
	"github.com/gin-gonic/gin"
	"time"
)

//每次拉取资源的时候，检测ID是否合法、Token是否过期
//登录之后，JWTtoken会存储在手机里，然后之后可以通过shouldbind 获取对应的request数据
//后端负责处理的逻辑是response

// JWTMiddleWare 鉴权中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}

		//检验用户JWT分发状况
		if tokenStr == "" {
			handler.ResponseError(c, handler.CodeUserNotLogin)
			c.Abort() //阻止执行
			return
		}

		//验证token合法性
		token, err := ParseToken(tokenStr)
		if err != nil {
			handler.ResponseError(c, handler.CodeTokenInvalid)
			c.Abort() //阻止执行
			return
		}

		//检查token是否过期
		if time.Now().Unix() > token.ExpiresAt {
			handler.ResponseError(c, handler.CodeTokenExpired)
			c.Abort() //阻止执行
			return
		}

		c.Set("user_id", token.UserId)
		c.Next()
	}
}
