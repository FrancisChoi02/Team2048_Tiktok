package model

type User struct {
	Id       int64  `json:"user_id" gorm:"id,omitempty;index;AUTO_INCREMENT:false"`
	Name     string `json:"username" gorm:"primary_key,name,omitempty"`
	Password string `json:"password" gorm:"notnull"`
}

type UserResponse struct {
	Id             int64  `json:"id" `
	Name           string `json:"name" `
	FollowCount    int64  `json:"follow_count" `
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow " `
	TotalFavorited int64  `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
}

type FriendResponse struct {
	Id             int64  `json:"user_id" `
	Name           string `json:"username" `
	FollowCount    int64  `json:"follow_count" `
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow " `
	TotalFavorited int64  `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
	Message        string `json:"message"`
	MsgType        int64  `json:"msgType"` // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}
