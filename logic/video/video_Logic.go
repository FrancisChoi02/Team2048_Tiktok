package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	logic "Team2048_Tiktok/logic/user"
	"Team2048_Tiktok/model"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var (
	videoIndexMap = map[string]struct{}{
		".mp4":  {},
		".avi":  {},
		".wmv":  {},
		".flv":  {},
		".mpeg": {},
		".mov":  {},
	}
)

// VideoPublish 用户视频投稿，将视频数据保存的在服务端本地
func VideoPublish(c *gin.Context, userId int64, form *multipart.Form) (err error) {
	videoData := form.File["data"]

	for _, file := range videoData {
		// 1.获取每个视频对应的标题
		title := c.PostForm("title")

		// 2.判断文件是否为合法的视频文件
		suffix := filepath.Ext(file.Filename)    //截取文件后缀
		if _, ok := videoIndexMap[suffix]; !ok { //判断文件是否为视频格式
			zap.L().Error("videoIndexMap[suffix] failed", zap.Error(err))
			return
		}

		// 2.为视频生成合适的文件名

		//使用snowflake为视频生成唯一ID
		videoId := model.GenID()
		videoIdStr := strconv.Itoa(int(videoId))

		// 3.生成视频路径并保存到本地
		filename := videoIdStr + "_" + suffix //以用户ID和当前时间戳 构成唯一视频名称
		savePath := filepath.Join("./static/video", filename)

		err = c.SaveUploadedFile(file, savePath)
		if err != nil {
			zap.L().Error("c.SaveUploadedFile() failed", zap.Error(err))
			return err
		}

		// 4.生成缩略图并保存到本地
		coverData, err := getSnapshot(savePath)
		if err != nil {
			zap.L().Error("getSnapshot() failed", zap.Error(err))
			return err
		}
		coverName := videoIdStr + "_" + "_cover" + ".jpeg"
		coverSavePath := filepath.Join("./static/cover", coverName)

		//使用 os.Create() 函数创建一个新文件，将该文件的句柄保存到 out 变量中
		out, err := os.Create(coverSavePath)
		if err != nil {
			zap.L().Error(" os.Create() failed", zap.Error(err))
			return err
		}

		// 将image.Image类型的图片数据编码为JPEG格式，并保存到 out变量 代表的文件中
		err = jpeg.Encode(out, coverData, nil)
		if err != nil {
			zap.L().Error(" jpeg.Encode() failed", zap.Error(err))
			return err
		}

		//关闭文件句柄，释放文件资源
		out.Close()

		//	5.将视频相关信息持久化到数据库中
		properFilePath, properCoverPath := GetVideoURL(filename, coverName)

		err = mysql.PostVideo(videoId, userId, properFilePath, properCoverPath, title)
		if err != nil {
			zap.L().Error(" mysql.PostVideo() failed", zap.Error(err))
			return err
		}

		//  6.增加用户的视频发布数
		if err = redis.RecordPublishNum(userId); err != nil {
			zap.L().Error(" redis.RecordPublishNum() failed", zap.Error(err))

		}

		//  7.独立返回每个视频上传成功的响应
		var res model.VideoUploadResponse
		res.Code = 0
		res.Msg = file.Filename + "上传成功"
		c.JSON(http.StatusOK, res)
	}
	return nil
}

// GetVideoURL 获得用于记录在数据库中的标准路径
func GetVideoURL(filename, covername string) (string, string) {
	ip := "192.168.31.169"
	port := 8080
	properFileName := fmt.Sprintf("http://%s:%d/static/video/%s", ip, port, filename)
	properCoverName := fmt.Sprintf("http://%s:%d/static/cover/%s", ip, port, covername)
	return properFileName, properCoverName
}

// GetVideoListByUserId 查询某个用户ID对应的视频投稿列表
func GetVideoListByUserId(userId int64) (*[]model.VideoResponse, error) {
	//1. 查找userId是否存在
	tmpUser := new(model.User)
	tmpUser.Id = userId
	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return nil, err
	}

	//2. 获取视频的基础信息
	videoList, err := mysql.GetVideoList(userId)
	if err != nil {
		zap.L().Error("mysql.GetVideoList() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	//3. 获取用户与视频的关系、视频的详细信息
	videoListFinal, err := GetVideoListDetail(videoList)
	if err != nil {
		zap.L().Error("redis.GetVideoListDetial() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	return videoListFinal, err
}

func FeedWithNoToken(latestTime int64) (*[]model.VideoResponse, int64, error) {

	// 1.获取视频流的基础信息
	feedList, err := mysql.GetFeedList(latestTime)
	if err != nil {
		zap.L().Error("mysql.GetFeedList() failed", zap.Error(err))
		return nil, 0, err //返回空切片
	}

	// 2. 获取next_Time
	nextTime := (*feedList)[0].CreatedAt

	// 3. 获取视频的详细信息
	videoListFinal, err := GetFeedListWithNoToken(feedList)
	if err != nil {
		zap.L().Error("redis.GetFeedListWithNoToken() failed", zap.Error(err))
		return nil, 0, err //返回空切片
	}

	return videoListFinal, nextTime, nil
}

func FeedWithToken(latestTime, userId int64) (*[]model.VideoResponse, int64, error) {
	// 1.获取视频流的基础信息
	feedList, err := mysql.GetFeedList(latestTime)
	if err != nil {
		zap.L().Error("mysql.GetFeedList() failed", zap.Error(err))
		return nil, 0, err //返回空切片
	}

	// 2. 获取next_Time
	nextTime := (*feedList)[0].CreatedAt

	// 3. 获取用户与视频的点赞关系、视频的详细信息
	videoListFinal, err := GetFeedListWithToken(userId, feedList)
	if err != nil {
		zap.L().Error("redis.GetFeedListWithNoToken() failed", zap.Error(err))
		return nil, 0, err //返回空切片
	}

	return videoListFinal, nextTime, nil
}

// GetFeedListWithNoToken  获取未登录用户的视频流
func GetFeedListWithNoToken(feedList *[]model.Video) (*[]model.VideoResponse, error) {
	// 1.判断视频列表是否为空
	if feedList == nil {
		err := errors.New("videoList has nothing")
		zap.L().Error("videoList has nothing")
		return nil, err
	}

	// 2.构建VideoResponse数组
	size := len(*feedList)
	feedListFull := make([]model.VideoResponse, size)

	//	3.每个视频单独处理
	for i, video := range *feedList {
		tmpVideoID := strconv.Itoa(int(video.Id))

		// a.获取视频发布者信息
		videoUser := &model.User{}
		videoUser.Id = video.UserId

		_, err := mysql.GetUser(videoUser)
		if err != nil {
			zap.L().Error(" GetUser() failed", zap.Error(err))
		}

		// 获取用户完整的信息
		user, err := logic.GetUserDetail(*videoUser)
		commentCount, favoriteCount, _ := redis.GetFeedListStatus("0", tmpVideoID)

		// 用户未登录，默认为false
		feedListFull[i].IsFavorite = false
		feedListFull[i].Author = user

		// b.获取视频的评论数
		feedListFull[i].CommentCount = commentCount

		// c.获取视频的点赞数
		feedListFull[i].FavoriteCount = favoriteCount

		// d.组装videoListFull单元
		assembleVideoListFull(&video, i, feedListFull)
	}
	return &feedListFull, nil

}

// GetFeedListWithToken  获取登录用户的视频流
func GetFeedListWithToken(userId int64, feedList *[]model.Video) (*[]model.VideoResponse, error) {
	// 1.判断视频列表是否为空
	if feedList == nil {
		err := errors.New("videoList has nothing")
		zap.L().Error("videoList has nothing")
		return nil, err
	}

	// 2.构建VideoResponse数组
	size := len(*feedList)
	feedListFull := make([]model.VideoResponse, size)

	//	3.每个视频单独处理
	for i, video := range *feedList {
		tmpVideoID := strconv.Itoa(int(video.Id))

		// a.获取视频发布者信息
		videoUser := &model.User{}
		videoUser.Id = video.UserId

		_, err := mysql.GetUser(videoUser)
		if err != nil {
			zap.L().Error(" GetUser() failed", zap.Error(err))
		}

		// 获取用户完整的信息
		user, err := logic.GetUserDetail(*videoUser)

		//查看当前用户 给 当前视频 的赞记录
		tmpUserID := strconv.Itoa(int(userId))
		commentCount, favoriteCount, ok := redis.GetFeedListStatus(tmpUserID, tmpVideoID)

		if ok == 1 {
			feedListFull[i].IsFavorite = true
		} else {
			feedListFull[i].IsFavorite = false
		}

		// b.获取视频的评论数
		feedListFull[i].CommentCount = commentCount

		// c.获取视频的点赞数
		feedListFull[i].FavoriteCount = favoriteCount

		// d.组装videoListFull单元
		feedListFull[i].Author = user
		assembleVideoListFull(&video, i, feedListFull)
	}
	return &feedListFull, nil
}

// GetVideoDetail 返回视频的完整数据
func GetVideoDetail(tmpVideo *model.Video) (*model.VideoResponse, error) {
	tmpVideoID := strconv.Itoa(int(tmpVideo.Id))
	videoFull := new(model.VideoResponse)

	tmpUserID := "0"
	commentCount, favoriteCount, _ := redis.GetFeedListStatus(tmpVideoID, tmpUserID)
	// a.获取视频的点赞数
	videoFull.FavoriteCount = favoriteCount

	// b.获取视频的评论数
	videoFull.CommentCount = commentCount

	// c.点赞设置为有
	videoFull.IsFavorite = true
	tmpUser := &model.User{}
	tmpUser.Id = tmpVideo.UserId

	_, err := mysql.GetUser(tmpUser)
	if err != nil {
		zap.L().Error(" GetUser() failed", zap.Error(err))
	}

	// 获取用户完整的信息
	user, err := logic.GetUserDetail(*tmpUser)

	// d.组装videoListFull单元
	videoFull.Author = user
	assembleVideoFull(tmpVideo, videoFull)

	return videoFull, nil
}

// GetVideoListDetail 返回用户投稿视频列表的完整数据
func GetVideoListDetail(videoList *[]model.Video) (*[]model.VideoResponse, error) {
	// 1.判断视频列表是否为空
	if videoList == nil {
		err := errors.New("videoList has nothing")
		zap.L().Error("videoList has nothing")
		return nil, err
	}

	// 2.构建VideoResponse数组
	size := len(*videoList)
	videoListFull := make([]model.VideoResponse, size)

	//	3.每个视频单独处理
	for i, video := range *videoList {
		tmpVideoID := strconv.Itoa(int(video.Id))

		// a.查看登录的用户有没有点赞当前视频
		tmpUser := &model.User{}
		tmpUser.Id = video.UserId

		_, err := mysql.GetUser(tmpUser)
		if err != nil {
			zap.L().Error(" GetUser() failed", zap.Error(err))
		}

		//查看当前用户 给 当前视频 的赞记录
		tmpUserID := strconv.Itoa(int(tmpUser.Id))
		commentCount, favoriteCount, ok := redis.GetFeedListStatus(tmpVideoID, tmpUserID)

		if ok == 1 {
			videoListFull[i].IsFavorite = true
		} else {
			videoListFull[i].IsFavorite = false
		}

		// b.获取视频的评论数
		videoListFull[i].CommentCount = commentCount

		// c.获取视频的点赞数
		videoListFull[i].FavoriteCount = favoriteCount

		// d.组装videoListFull单元
		// 获取用户完整的信息
		user, err := logic.GetUserDetail(*tmpUser)
		videoListFull[i].Author = user
		assembleVideoListFull(&video, i, videoListFull)
	}
	return &videoListFull, nil
}

// assembleVideoListFull  组装videoListFull单元
func assembleVideoListFull(video *model.Video, index int, videoListFull []model.VideoResponse) {
	videoListFull[index].Id = video.Id
	videoListFull[index].PlayUrl = video.PlayUrl
	videoListFull[index].CoverUrl = video.CoverUrl
	videoListFull[index].Title = video.Title
	videoListFull[index].CreatedAt = video.CreatedAt
}

// assembleVideoFull  组装视频
func assembleVideoFull(video *model.Video, videoFull *model.VideoResponse) {
	videoFull.Id = video.Id
	videoFull.PlayUrl = video.PlayUrl
	videoFull.CoverUrl = video.CoverUrl
	videoFull.Title = video.Title
	videoFull.CreatedAt = video.CreatedAt
}
