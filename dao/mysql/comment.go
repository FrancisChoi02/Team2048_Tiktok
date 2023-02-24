package mysql

import (
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/model"
	"errors"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"time"
)

// PostComment 将comment信息保存到数据库中
func PostComment(comment *model.Comment) error {
	// 使用 Create 方法向数据库中插入记录
	return DB.Create(comment).Error
}

// GetComment  查询comment信息，并将完整的信息填充到comment结构体
func GetComment(comment *model.Comment) (boolstring bool, err error) {
	boolstring = false
	if err := DB.Where("id = ?", comment.Id).First(comment).Error; err != nil { //这里曾经是&comment
		if gorm.IsRecordNotFoundError(err) {
			// 处理记录不存在错误
			zap.L().Error("Comment doesn't exist", zap.Error(err))
		} else {
			// 处理其他错误
			zap.L().Error("DB.Where(\"id = ?\", comment.Id).First(comment) failed", zap.Error(err))
		}
		return boolstring, err
	}

	boolstring = true
	return boolstring, err
}

// RemoveComment  将相关Id的comment从MySQL中删除
func RemoveComment(commentId int64) error {
	//comment结构体中，有comment.Id这个成员变量
	return DB.Delete(&model.Comment{}, commentId).Error
}

// GetCommentResponseList 根据commentId的切片，获取完整的用于返回的comment切片
func GetCommentResponseList(commentIdList []int64) (*[]model.CommentResponse, error) {
	// 1.判断切片列表是否为空
	if commentIdList == nil {
		err := errors.New("commentIdList has nothing")
		zap.L().Error("commentIdList has nothing")
		return nil, err
	}

	// 2.构建返回切片
	size := len(commentIdList)
	commentListFull := make([]model.CommentResponse, size)

	for i, commentId := range commentIdList {
		//逐个commentResponse进行填充
		comment := new(model.Comment) //获取comment基础信息
		comment.Id = commentId
		_, err := GetComment(comment)
		if err != nil {
			zap.L().Error("comment has nothing")
		}

		user := new(model.User) //获取用户信息
		user.Id = comment.UserId
		_, err = GetUser(user)
		if err != nil {
			zap.L().Error("User do not exist")
		}

		//补全用户信息
		userFulll, err := redis.GetUserDetail(*user)

		tmpTime := time.Unix(comment.CreatedAt, 0) //将unix时间转回 time.Time
		dateStr := tmpTime.Format("01-02")

		//组装返回的commentResponse并赋值给当前单元
		commentResponse := new(model.CommentResponse)
		commentResponse.Id = comment.Id
		commentResponse.User = userFulll
		commentResponse.Content = comment.Content
		commentResponse.CreateDate = dateStr

		commentListFull[i] = *commentResponse
	}

	return &commentListFull, nil

}
