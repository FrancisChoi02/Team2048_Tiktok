package model

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ParamUserInfo 用户信息请求参数
type ParamUserInfo struct {
	User_id int64  `form:"user_id"binding:"required"`
	Token   string `form:"token"binding:"required"`
}
