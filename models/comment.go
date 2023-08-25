package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID              uint   `json:"user_id" gorm:"not null"`
	PostID              uint   `json:"post_id" gorm:"not null"`
	CommentContent      string `json:"comment_content" gorm:"text;not null"`
	CommentLikeCount    uint   `json:"comment_like_count" gorm:"not null;default:0"`
	CommentDislikeCount uint   `json:"comment_dislike_count" gorm:"not null;default:0"`

	// relationship
	User            User              `json:"-"`
	Post            Post              `json:"-"`
	UserLikeComment []UserLikeComment `json:"-"`
}
