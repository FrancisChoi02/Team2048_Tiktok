package router

import (
	"Team2048_Tiktok/logger"
	"github.com/gin-gonic/gin"
	//swaggerFiles "github.com/swaggo/files"
	//gs "github.com/swaggo/gin-swagger"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(2*time.Second, 1))
	//使用自定义中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//社区分类相关路由

	//swagger接口文档所需要的路由（暂未实现）
	//r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Pages not found",
		})
	})

	return r
}
