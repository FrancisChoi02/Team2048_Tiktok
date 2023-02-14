package model

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamUserInfo 用户信息请求参数
type ParamUserInfo struct {
	User_id int64  `json:"user_id"`
	Token   string `json:"token"`
}
