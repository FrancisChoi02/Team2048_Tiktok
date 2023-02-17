package logic

import (
	"bytes"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"image/jpeg"
	"os/exec"
)

// getSnapshot 获取视频截图
func getSnapshot(filePath string) ([]byte, error) {

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

	/*
		buffer := bytes.NewBuffer(nil)
		err := ffmpeg.Input(filePath).
			Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 1)}).
			Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
			WithOutput(buffer, os.Stdout).
			Run()
		if err != nil {
			zap.L().Error("生成缩略图失败", zap.Error(err))
			return nil, err
		}

		img, err := imaging.Decode(buffer)
		if err != nil {
			zap.L().Error("图片数据解码失败", zap.Error(err))
			return nil, err
		}
		bytes.NewReader(out)
	*/

	buf := new(bytes.Buffer)
	jpeg.Encode(buf, img, nil)

	return buf.Bytes(), err
}
