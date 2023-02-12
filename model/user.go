package model

type User struct {
	Id            int64  `json:"user_id" gorm:"primary_key,id,omitempty;index"`
	Name          string `json:"username" gorm:"primary_key,name,omitempty"`
	Password      string `json:"password" gorm:"notnull"`
	FollowCount   int64  `json:"-" gorm:"follow_count,omitempty"`
	FollowerCount int64  `json:"-" gorm:"follower_count,omitempty"`
	IsFollow      bool   `json:"-" gorm:"is_follow,omitempty"`
}
