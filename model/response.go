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
	Code int32        // 状态码，0-成功，其他值-失败
	Msg  string       // 返回状态描述
	User UserResponse // 返回用户信息反馈结构体
}

type VideoUploadResponse struct {
	Code int32  // 状态码，0-成功，其他值-失败
	Msg  string // 返回状态描述
}

type FeedResponse struct {
	Code      int32  // 状态码，0-成功，其他值-失败
	Msg       string // 返回状态描述
	NextTime  int64  //下一次视频刷新时间
	VideoList *[]VideoResponse
}

type VideoListResponse struct {
	Code      int32  // 状态码，0-成功，其他值-失败
	Msg       string // 返回状态描述
	VideoList *[]VideoResponse
}

type FavorateActionResponse struct {
	Code int32  // 状态码，0-成功，其他值-失败
	Msg  string // 返回状态描述
}

type CommentActionResponse struct {
	Code    int32           // 状态码，0-成功，其他值-失败
	Msg     string          // 返回状态描述
	Comment CommentResponse //返回评论结果
}

type CommentListResponse struct {
	Code        int32              // 状态码，0-成功，其他值-失败
	Msg         string             // 返回状态描述
	CommentList *[]CommentResponse //返回评论列表
}

type FollowRelationResponse struct {
	Code     int32           // 状态码，0-成功，其他值-失败
	Msg      string          // 返回状态描述
	UserList *[]UserResponse //返回用户列表
}

type FriendListResponse struct {
	Code       int32             // 状态码，0-成功，其他值-失败
	Msg        string            // 返回状态描述
	FriendList *[]FriendResponse //返回聊天好友列表
}

type ChatHistoryResponse struct {
	Code        int32      // 状态码，0-成功，其他值-失败
	Msg         string     // 返回状态描述
	MessageList *[]Message //返回消息列表
}
