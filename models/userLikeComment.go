package models

type UserLikeComment struct {
	ID        uint `json:"id" gorm:"primary_key;not null;"`
	UserID    uint `json:"user_id" gorm:"not null;"`
	CommentID uint `json:"comment_id" gorm:"not null;"`
	Status    uint `json:"status"`

	// Relationship
	User    User    `json:"-"`
	Comment Comment `json:"-"`
}
