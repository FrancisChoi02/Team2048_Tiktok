package mysql

import (
	"Team2048_Tiktok/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

const key = "Tiktok2048"

// encryptPassword MD5加密用户密码
func encryptPassword(oPassword string) string {
	//用key进行加密
	h := md5.New()
	h.Write([]byte(key))
	tmp := h.Sum([]byte(oPassword))

	return hex.EncodeToString(tmp)
}

// InsertUser 往数据库中增添新的用户
func InsertUser(user *model.User) (err error) {
	// 1. 将用户传来的密码加密
	user.Password = encryptPassword(user.Password)

	// 1使用 Create 方法向数据库中插入记录
	if err := DB.Create(user).Error; err != nil {
		// 数据插入错误
		return ErrorInserFaied
	}
	return
}

// CheckUserExist 查找用户是否存在
func CheckUserExist(user *model.User) (err error) {
	if err := DB.Where("name = ?", user.Name).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil
		}
	}
	return ErrorUserExist
}

// GetUser 获取完整的用户信息,查询用户是否存在
func GetUser(user *model.User) (boolstring bool, err error) {
	boolstring = false
	fmt.Println("this is GetUser, useridIsDDDDDDDD", user.Id)

	if err := DB.Where("id = ?", user.Id).First(user).Error; err != nil { //这里曾经是&user
		if gorm.IsRecordNotFoundError(err) {
			// 处理记录不存在错误
			zap.L().Error("User doesn't exist", zap.Error(err))
		} else {
			// 处理其他错误
			zap.L().Error("DB.Where(\"id = ?\", user.Id).First(user) failed", zap.Error(err))
		}
		return boolstring, ErrorUserNotExist
	}

	boolstring = true
	return boolstring, err
}

// GetUserList  根据Id切片，获取用户结构体切片
func GetUserList(idList []int64) ([]model.User, error) {
	userList := make([]model.User, len(idList))

	//逐个User进行获取
	for i, id := range idList {
		tmpUser := new(model.User)
		tmpUser.Id = id

		_, err := GetUser(tmpUser)
		if err != nil {
			zap.L().Error("GetUser() failed", zap.Error(err))
			return nil, err
		}

		userList[i] = *tmpUser
	}

	return userList, nil
}

// Login 检查登录用户合法性，获取完整的用户信息
func Login(user *model.User) (err error) {
	// 1.临时结构体保存数据库查询结果
	var tmpUser model.User
	fmt.Println("this is mysql username", user.Name)

	if err := DB.Where("name = ?", user.Name).First(&tmpUser).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			// 处理记录不存在错误
			zap.L().Error("User doesn't exist", zap.Error(err))
		} else {
			// 处理其他错误
			return ErrorInserFaied
		}
		return ErrorUserNotExist
	}

	// 2.将数据库中的密码密文进行比对
	user.Password = encryptPassword(user.Password)

	if tmpUser.Password != user.Password {
		return ErrorInvalidPassword
	}
	//返回登录成功的用户ID
	user.Id = tmpUser.Id

	return err
}
