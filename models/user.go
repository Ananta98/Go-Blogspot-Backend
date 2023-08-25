package models

import (
	"blogspot-project/utils/token"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"size:255;not null;unique" json:"name"`
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Email    string `gorm:"size:255;not null;unique" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	ImageUrl string `gorm:"size:255;" json:"image_url"`
	Role     uint   `gorm:"not null;default:2"`

	// Relationship
	Posts               []Post            `json:"-"`
	Comments            []Comment         `json:"-"`
	UserListLikePost    []UserLikePost    `json:"-"`
	UserListLikeComment []UserLikeComment `json:"-"`
}

const ADMIN_USER_ROLE = 1
const BLOGGER_USER_ROLE = 2
const READER_USER_ROLE = 3

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func PasswordHashing(passwordInput string) (string, error) {
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(passwordInput), bcrypt.DefaultCost)
	if errPassword != nil {
		return "", errPassword
	}
	return string(hashedPassword), nil
}

func (u *User) CreateUser(db *gorm.DB) (*User, error) {
	hashedPassword, err := PasswordHashing(u.Password)
	if err != nil {
		return &User{}, err
	}
	u.Password = hashedPassword
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	if err := db.Create(&u).Error; err != nil {
		return &User{}, err
	}
	return &User{}, nil
}

func LoginValid(username, password string, db *gorm.DB) (string, error) {
	foundUser := User{}
	if err := db.Where("username = ?", username).Take(&foundUser).Error; err != nil {
		return "", err
	}
	if err := VerifyPassword(password, foundUser.Password); err != nil {
		return "", err
	}
	token, err := token.GenerateToken(foundUser.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
