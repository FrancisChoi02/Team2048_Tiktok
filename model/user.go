package model

type User struct {
	Id       int64  `json:"user_id" gorm:"id,omitempty;index;AUTO_INCREMENT:false"`
	Name     string `json:"username" gorm:"primary_key,name,omitempty"`
	Password string `json:"password" gorm:"notnull"`
}

type UserResponse struct {
	Id             int64  `json:"user_id" `
	Name           string `json:"username" `
	Password       string `json:"password"`
	FollowCount    int64  `json:"follow_count" `
	FollowerCount  int64  `json:"follower_count"`
	IsFollow       bool   `json:"is_follow " `
	TotalFavorited int64  `json:"total_favorited"`
	WorkCount      int64  `json:"work_count"`
	FavoriteCount  int64  `json:"favorite_count"`
}
