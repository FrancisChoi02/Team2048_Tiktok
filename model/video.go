package model

type Video struct {
	Id        int64  `json:"id,omitempty" gorm:"primary_key;index;AUTO_INCREMENT:false"`
	UserId    int64  `json:"user_id,omitempty"`
	PlayUrl   string `json:"play_url,omitempty"`
	CoverUrl  string `json:"cover_url,omitempty"`
	Title     string `json:"title,omitempty"`
	CreatedAt int64  `json:"create_time"`
}

type VideoResponse struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key;index;AUTO_INCREMENT:false"`
	Author        User   `json:"author,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	Title         string `json:"title,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	CreatedAt     int64  `json:"create_time"`
}

type List struct {
	Videos []*Video `json:"video_list,omitempty"`
}
