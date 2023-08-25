package controllers

import (
	"blogspot-project/models"
	"blogspot-project/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	ImageUrl string `json:"image_url"`
	Role     uint   `json:"role" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdatePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Register godoc
// @Summary Register new user or create new user.
// @Description registering a user to get access blog.
// @Tags Auth
// @Param Body body RegisterInput true "json body to register a user or create new user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/register [post]
func RegisterNewUser(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var inputRegister RegisterInput
	if err := ctx.ShouldBindJSON(&inputRegister); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newUser := models.User{
		Name:     inputRegister.Name,
		Email:    inputRegister.Email,
		Username: inputRegister.Username,
		ImageUrl: inputRegister.ImageUrl,
		Password: inputRegister.Password,
		Role:     inputRegister.Role,
	}
	_, err := newUser.CreateUser(db)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := map[string]string{
		"name":     inputRegister.Name,
		"username": inputRegister.Username,
		"email":    inputRegister.Email,
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Registration success", "user": result})
}

// Login godoc
// @Summary Login using registered user.
// @Description login into blog to get access all blog list and CRUD blogs
// @Tags Auth
// @Param Body body LoginInput true "json body to Login existing user"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /auth/login [post]
func LoginUser(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var inputLogin LoginInput
	if err := ctx.ShouldBindJSON(&inputLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := models.LoginValid(inputLogin.Username, inputLogin.Password, db)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Login success", "token": token})
}

// UpdatePassword godoc
// @Summary Update password for current user.
// @Description Ability user to change their password.
// @Tags Auth
// @Param Body body UpdatePasswordInput true "json body to update password for current existing user"
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /login/update-password [patch]
func UpdatePassword(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	oldUser := models.User{}
	if err := db.Where("ID = ?", id).Take(&oldUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatePasswordInput := UpdatePasswordInput{}
	if err := ctx.ShouldBindJSON(&updatePasswordInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = models.VerifyPassword(updatePasswordInput.OldPassword, oldUser.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := models.PasswordHashing(updatePasswordInput.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser := models.User{
		Name:     oldUser.Name,
		Username: oldUser.Username,
		Email:    oldUser.Email,
		ImageUrl: oldUser.ImageUrl,
		Password: hashedPassword,
	}
	if err := db.Model(&oldUser).Updates(&updatedUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Update Password Success"})
}
