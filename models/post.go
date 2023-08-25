package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID             uint   `gorm:"not null;" json:"user_id"`
	ArticleTitle       string `gorm:"size:255;not null;" json:"article_title"`
	ArticleDescription string `gorm:"text;" json:"article_description"`
	CategoryID         uint   `json:"category_id"`
	PostLikeCount      uint   `json:"post_like_count" gorm:"not null;default:0"`
	PostDislikeCount   uint   `json:"post_dislike_count" gorm:"not null;default:0"`

	// Relationship
	User         User           `json:"-"`
	Comments     []Comment      `json:"-"`
	Category     Category       `json:"-"`
	UserLikePost []UserLikePost `json:"-"`
}
