package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/model"
	"bytes"
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
func VideoPublish(c *gin.Context, userId int64, title string, form *multipart.Form) (err error) {
	videoData := form.File["data"]

	for _, file := range videoData {

		// 1.判断文件是否为合法的视频文件
		suffix := filepath.Ext(file.Filename)    //截取文件后缀
		if _, ok := videoIndexMap[suffix]; !ok { //判断文件是否为视频格式
			zap.L().Error("videoIndexMap[suffix] failed", zap.Error(err))
			return
		}

		// 2.为视频生成合适的文件名
		//timeStamp := time.Now().Format("2006-01-02 15:04")
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
		fmt.Println("Now the video is successfully saved")

		// 4.生成缩略图并保存到本地
		coverData, err := getSnapshot(savePath) //输出图片类型可以优化
		if err != nil {
			zap.L().Error("getSnapshot() failed", zap.Error(err))
			return err
		}
		coverName := videoIdStr + "_" + "_cover" + ".jpeg"
		coverSavePath := filepath.Join("./static/cover", coverName)

		// 解码JPEG格式的图片数据，转换为image.Image类型
		img, err := jpeg.Decode(bytes.NewReader(coverData))
		if err != nil {
			zap.L().Error(" jpeg.Decode() failed", zap.Error(err))
			return err
		}

		//使用 os.Create() 函数创建一个新文件，将该文件的句柄保存到 out 变量中
		out, err := os.Create(coverSavePath)
		if err != nil {
			zap.L().Error(" os.Create() failed", zap.Error(err))
			return err
		}

		// 将image.Image类型的图片数据编码为JPEG格式，并保存到 out变量 代表的文件中
		err = jpeg.Encode(out, img, nil)
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

func GetVideoURL(filename, covername string) (string, string) {
	ip := "192.168.31.169"
	port := 8080

	properFileName := fmt.Sprintf("http://%s:%d/static/video/%s", ip, port, filename)
	properCoverName := fmt.Sprintf("http://%s:%d/static/cover/%s", ip, port, covername)
	return properFileName, properCoverName
}
