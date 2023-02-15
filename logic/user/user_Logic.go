package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/middleware"
	"Team2048_Tiktok/model"
	"go.uber.org/zap"
)

// SignUp 将新注册的用户信息保存到数据库
func SignUp(p *model.ParamSignUp) (int64, error) {
	var user model.User
	user.Name = p.Username

	//1. 查找Username是否已经存在
	boolstring, err := mysql.GetUser(&user)
	if err != nil { //从数据库中查找到对应的用户
		zap.L().Error(" mysql.GetUser() failed")
		return 0, err
	}

	if boolstring {
		zap.L().Error("User already exist")
		return 0, err
	}

	//2. 分发用户ID
	tmpID := model.GenID()

	//3. 构建用户实例
	user = model.User{
		Id:            tmpID,
		Name:          p.Username,
		Password:      p.Password,
		FollowCount:   0,
		FollowerCount: 0,
	}

	//4. 数据保存到MySQL,返回响应
	return tmpID, mysql.InsertUser(&user)
}

// Login 用户登录
func Login(p *model.ParamSignUp) (tmpId int64, token string, err error) {

	//1. 构建用户实例
	user := model.User{
		Id:            0,
		Name:          p.Username,
		Password:      p.Password,
		FollowCount:   0,
		FollowerCount: 0,
	}
	//2. 获取登录后的User结构体 传递的是user指针，因此可以获得user.UserID
	if err := mysql.Login(&user); err != nil { //从数据库中查找后，user的值会被覆盖一遍
		zap.L().Error(" mysql.Login() failed", zap.Error(err))
		return 0, "", err
	}

	//3. 颁发JWT
	token, err = middleware.ReleaseToken(tmpId)
	if err != nil {
		zap.L().Error(" middleware.ReleaseToken() failed", zap.Error(err))
		return
	}
	//4. 获得用于返回的user_id
	tmpId = user.Id

	return tmpId, token, err
}

// GetUser 通过Id获得对应的User结构体
func GetUser(id int64) (model.User, error) {
	//2. 构建用户实例
	user := model.User{
		Id:            id,
		Name:          "",
		Password:      "",
		FollowCount:   0,
		FollowerCount: 0,
	}
	_, err := mysql.GetUser(&user)
	if err != nil {
		//数据库查询错误
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return user, err
	}

	return user, nil
}
