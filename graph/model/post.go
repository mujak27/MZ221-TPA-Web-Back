package model

import "time"

type PostLike struct {
	ID     string `json:"id" gorm:"type:varchar(191)"`
	PostId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LikeId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Post struct {
	ID             string `gorm:"type:varchar(191)"`
	Text           string `json:"Text"`
	SenderId       string
	Sender         *User      `json:"Sender" gorm:"reference:User"`
	Comments       []*Comment `json:"Comments" gorm:"foreignKey:PostId;references:ID"`
	Likes          []*User    `json:"Likes" gorm:"many2many:post_likes"`
	AttachmentLink string     `json:"AttachmentLink"`
	CreatedAt      time.Time  `json:"CreatedAt"`
}

type Comment struct {
	ID          string     `json:"ID" gorm:"type:varchar(191)"`
	Text        string     `json:"Text"`
	SenderId    string     `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Sender      *User      `json:"Sender" gorm:"reference:User"`
	RepliedToId *string    `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Replies     []*Comment `json:"Replies" gorm:"foreignKey:RepliedToId"`
	PostId      string     `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Post        *Post      `json:"Sender" gorm:"reference:Post"`
	Likes       []*User    `json:"Likes" gorm:"many2many:comment_likes"`
	CreatedAt   time.Time  `json:"CreatedAt"`
}

type CommentLike struct {
	ID        string `json:"id" gorm:"type:varchar(191)"`
	CommentId string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LikeId    string `gorm:"type:varchar(191);constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type InputComment struct {
	PostID      string `json:"PostId"`
	RepliedToID string `json:"RepliedToId"`
	Text        string `json:"Text"`
}
