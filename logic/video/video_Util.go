package logic

import (
	"bytes"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"image"
	"os/exec"
)

// getSnapshot 获取视频截图
func getSnapshot(filePath string) (image.Image, error) {

	cmd := exec.Command("ffmpeg", "-i", filePath, "-vframes", "1", "-f", "image2", "pipe:1")
	buffer, err := cmd.Output()
	if err != nil {
		zap.L().Error("生成缩略图失败", zap.Error(err))
		return nil, err
	}

	img, err := imaging.Decode(bytes.NewReader(buffer))
	if err != nil {
		zap.L().Error("图片数据解码失败", zap.Error(err))
		return nil, err
	}

	return img, err
}
