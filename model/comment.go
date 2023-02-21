package model

type CommentResponse struct {
	Id         int64  `json:"comment_id" gorm:"primary_key;index;AUTO_INCREMENT:false"`
	User       User   `json:"user"`
	Content    string `json:"comment_text"`
	CreateDate string `json:"create_date" gorm:"-"`
}

type Comment struct {
	Id        int64  `json:"comment_id" gorm:"primary_key;index;AUTO_INCREMENT:false"`
	VideoId   int64  `json:"video_id"`
	UserId    int64  `json:"user_id"`
	Content   string `json:"comment_text"`
	CreatedAt int64  `json:"create_time"`
}
