package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
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

// GetVideo 获取视频，判断对应Id是视频是否存在
func GetVideoById(videoId int64) (model.Video, error) {
	// 构建视频实例
	video := model.Video{}
	video.Id = videoId

	_, err := mysql.GetVideo(&video)
	if err != nil {
		zap.L().Error("mysql.GetVideo() failed", zap.Error(err))
		return video, err
	}
	return video, nil
}

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

		//  6.独立返回每个视频上传成功的响应
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
	videoListFinal, err := redis.GetVideoListDetail(videoList)
	if err != nil {
		zap.L().Error("redis.GetVideoListDetial() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	return videoListFinal, err
}

func FeedWithNoToken(latestTime int64) (*[]model.VideoResponse, error) {

	// 1.获取视频流的基础信息
	feedList, err := mysql.GetFeedList(latestTime)
	if err != nil {
		zap.L().Error("mysql.GetFeedList() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	// 2. 获取视频的详细信息
	videoListFinal, err := redis.GetFeedListWithNoToken(feedList)
	if err != nil {
		zap.L().Error("redis.GetFeedListWithNoToken() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	return videoListFinal, nil
}

func FeedWithToken(latestTime, userId int64) (*[]model.VideoResponse, error) {
	// 1.获取视频流的基础信息
	feedList, err := mysql.GetFeedList(latestTime)
	if err != nil {
		zap.L().Error("mysql.GetFeedList() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	// 2. 获取用户与视频的点赞关系、视频的详细信息
	videoListFinal, err := redis.GetFeedListWithToken(userId, feedList)
	if err != nil {
		zap.L().Error("redis.GetFeedListWithNoToken() failed", zap.Error(err))
		return nil, err //返回空切片
	}

	return videoListFinal, nil
}
