package models

type UserLikePost struct {
	ID     uint `json:"id" gorm:"primary_key;not null;"`
	UserID uint `json:"user_id" gorm:"not null;"`
	PostID uint `json:"post_id" gorm:"not null;"`
	Status uint `json:"status"`

	// Relationship
	User User `json:"-"`
	Post Post `json:"-"`
}
