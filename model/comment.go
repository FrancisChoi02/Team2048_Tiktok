package model

import "time"

type Comment struct {
	Id         int64     `json:"comment_id" gorm:"primary_key;index"`
	VideoId    int64     `json:"video_id"   gorm:"primary_key;index"`
	User       User      `json:"user"`
	Content    string    `json:"comment_text"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date" gorm:"-"`
}
