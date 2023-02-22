package router

import (
	_ "Team2048_Tiktok/docs" // 上一步swagger init 生成的docs也要导入
	"Team2048_Tiktok/handler/user"
	"Team2048_Tiktok/handler/video"
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

	//用户接口
	userAPI := r.Group("/douyin/user")
	{
		userAPI.GET("/", middleware.JWTMiddleware(), user.UserInfoHandler)
		userAPI.POST("/login/", user.UserLoginHandler)
		userAPI.POST("/register/", user.UserSignUpHandler)
	}

	//视频接口
	r.GET("/douyin/feed/", video.FeedVideoListHandler)
	videoAPI := r.Group("/douyin/publish")
	videoAPI.Use(middleware.JWTMiddleware())
	{
		videoAPI.POST("/action/", video.VideoPublishHandler)
		videoAPI.GET("/list/", video.VideoListHandler)
	}

	//互动接口
	interactAPI := r.Group("/douyin/favorite")
	interactAPI.Use(middleware.JWTMiddleware())
	{
		interactAPI.POST("/action/", video.FavorateActionHandler)
		interactAPI.GET("/list/", video.FavorateListHandler)
	}
	commentAPI := r.Group("/douyin/comment")
	commentAPI.Use(middleware.JWTMiddleware())
	{
		commentAPI.POST("/action/", video.CommentActionHandler)
		commentAPI.GET("/list/ ", video.CommentListHandler)
	}

	//社交接口
	relationAPI := r.Group("/douyin/relation")
	relationAPI.Use(middleware.JWTMiddleware())
	{
		relationAPI.POST("/action/", user.RelationshipHandler)
		relationAPI.GET("/follow/list/", user.FollowRelationHandler)
		relationAPI.GET("/follower/list/", user.FanRelationHandler)
		relationAPI.GET("/friend/list/", user.FriendRelationHandler)
	}
	//消息
	messageAPI := r.Group("/douyin/message")
	messageAPI.Use(middleware.JWTMiddleware())
	{
		messageAPI.POST("/action/", user.MessageHandler)
		messageAPI.GET("/chat/", user.ChatHistoryHandler)
	}

	//swagger接口文档所需要的路由
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	//失配路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Pages not found",
		})
	})

	return r
}
