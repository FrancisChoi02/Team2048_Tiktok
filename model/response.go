package model

// 通用的返回消息结构体
type CommonResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type UserDetailResponse struct {
	Code    int32  // 状态码，0-成功，其他值-失败
	Msg     string // 返回状态描述
	User_id int64  // 用户id
	Token   string // 用户鉴权token
}

type UserInfoResponse struct {
	Code int32  // 状态码，0-成功，其他值-失败
	Msg  string // 返回状态描述
	User User   // 返回用户结构体
}

type VideoUploadResponse struct {
	Code int32  // 状态码，0-成功，其他值-失败
	Msg  string // 返回状态描述
}

type VideoListResponse struct {
	Code      int32  // 状态码，0-成功，其他值-失败
	Msg       string // 返回状态描述
	VideoList *[]VideoResponse
}
