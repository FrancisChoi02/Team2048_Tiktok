package model

type Message struct {
	Id         int64  `json:"comment_id" gorm:"primary_key;index;AUTO_INCREMENT:false"`
	ToUserId   int64  `json:"to_user_id"`
	FromUserId int64  `json:"from_user_id"`
	Content    string `json:"comment_text"`
	CreatTime  string `json:"create_time"`
}
