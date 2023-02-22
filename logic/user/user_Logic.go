package logic

import (
	"Team2048_Tiktok/dao/mysql"
	"Team2048_Tiktok/dao/redis"
	"Team2048_Tiktok/middleware"
	"Team2048_Tiktok/model"
	"go.uber.org/zap"
)

// SignUp 将新注册的用户信息保存到数据库
func SignUp(p *model.ParamSignUp) (int64, error) {
	var user model.User
	user.Name = p.Username

	//1. 查找Username是否已经存在
	err := mysql.CheckUserExist(&user)
	if err != nil { //从数据库中查找到对应的用户
		zap.L().Error(" mysql.GetUser() failed")
		return 0, err
	}

	//2. 分发用户ID
	tmpID := model.GenID()

	//3. 构建用户实例
	user = model.User{
		Id:       tmpID,
		Name:     p.Username,
		Password: p.Password,
	}

	//4. 数据保存到MySQL,返回响应
	return tmpID, mysql.InsertUser(&user)
}

// Login 用户登录
func Login(p *model.ParamSignUp) (tmpId int64, token string, err error) {

	//1. 构建用户实例
	user := model.User{
		Id:       0,
		Name:     p.Username,
		Password: p.Password,
	}
	//2. 获取登录后的User结构体 传递的是user指针，因此可以获得user.UserID
	if err := mysql.Login(&user); err != nil { //从数据库中查找后，user的值会被覆盖一遍
		zap.L().Error(" mysql.Login() failed", zap.Error(err))
		return 0, "", err
	}

	//3. 颁发JWT
	token, err = middleware.ReleaseToken(user.Id)
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
	// 构建用户实例
	user := model.User{
		Id:       id,
		Name:     "",
		Password: "",
	}
	_, err := mysql.GetUser(&user)
	if err != nil {
		//数据库查询错误
		zap.L().Error("mysql.GetUser() failed", zap.Error(err))
		return user, err
	}

	return user, nil
}

// GetUserDetail 补充用于返回的用户结构体
func GetUserDetail(user model.User) (userResponse model.UserResponse, err error) {

	userResponse, err = redis.GetUserDetail(user)
	if err != nil {
		zap.L().Error("redis.GetUserDetail() failed", zap.Error(err))
		return userResponse, err
	}
	return userResponse, nil
}

// CheckIsFollow  查看token用户是否关注当前用户
func CheckIsFollow(user *model.UserResponse, tokenUserId int64) {
	//查看自己的信息时，标记为关注
	if user.Id == tokenUserId {
		user.IsFollow = true
		return
	}

	//查看token用户是否关注了user
	record := redis.GetFollowStatus(tokenUserId, user.Id)

	if record == 1 {
		user.IsFollow = true
	} else {
		user.IsFollow = false
	}

}
