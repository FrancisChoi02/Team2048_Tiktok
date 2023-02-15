package router

import (
	_ "Team2048_Tiktok/docs" // 上一步swagger init 生成的docs也要导入
	"Team2048_Tiktok/handler/user"
	"Team2048_Tiktok/logger"
	"Team2048_Tiktok/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // gin设置成发布模式
	}
	r := gin.New()
	//使用自定义中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//基础接口
	userAPI := r.Group("/douyin/user")
	{
		userAPI.GET("/", middleware.JWTMiddleware(), user.UserInfoHandler)
		userAPI.POST("/login/", user.UserLoginHandler)
		userAPI.POST("/register/", user.UserSignUpHandler)
	}

	//videoAPI := r.Group("/douyin/publish")
	{
		//r.GET("/douyin/feed/", video.FeedVideoListHandler)
		//videoAPI.GET("/list/", middleware.NoAuthToGetUserId(), video.QueryVideoListHandler)
		//videoAPI.POST("/action/", middleware.JWTMiddleware(), video.PublishVideoHandler)
	}

	//swagger接口文档所需要的路由（暂未实现）
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	//失配路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Pages not found",
		})
	})

	return r
}
