package model

import "time"

type Video struct {
	Id           int64     `json:"id,omitempty" gorm:"primary_key;index;AUTO_INCREMENT:false"`
	Author       User      `json:"author,omitempty"`
	PlayUrl      string    `json:"play_url,omitempty"`
	CoverUrl     string    `json:"cover_url,omitempty"`
	CommentCount int64     `json:"comment_count,omitempty"`
	Title        string    `json:"title,omitempty"`
	CreatedAt    time.Time `json:"-"`
}
