package models

type Category struct {
	ID   uint   `json:"id" gorm:"primary_key"`
	Name string `json:"name" gorm:"not null"`

	// Relationship
	Posts []Post `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
