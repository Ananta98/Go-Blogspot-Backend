package controllers

import (
	"blogspot-project/models"
	"blogspot-project/utils"
	"blogspot-project/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UpdateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	ImageUrl string `json:"image_url"`
}

// GetCurrentUserProfile godoc
// @Summary Get current user detail that have login
// @Description login into blog to get all user detail
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /user/{id} [get]
func GetCurrentUserProfile(ctx *gin.Context) {
	u := models.User{}
	db := ctx.MustGet("db").(*gorm.DB)
	id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := db.Where("id = ?", id).Take(&u).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Get user profile success", "data": u})
}

// GetListUsers godoc
// @Summary Get all User list.
// @Description Get all users that have been registered except current user.
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param   current_page      query    int        false        "current_page"
// @Param   page_size         query    int        false        "page_size"
// @Param   input_search      query    string     false        "input_search"
// @Success 200 {object} map[string]interface{}
// @Router /user [get]
func GetListUsers(ctx *gin.Context) {
	var users []models.User
	db := ctx.MustGet("db").(*gorm.DB)
	id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	limit, offset, err := utils.GetPagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Where("id != ? AND name LIKE ?", id, "%"+ctx.Query("input_search")+"%").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Get all users success", "data": users})
}

// UpdateCurrentUser godoc
// @Summary Update current user.
// @Description Update current user without update password that have logged in into blog.
// @Tags User
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /login/update-current-user [patch]
func UpdateCurrentUser(ctx *gin.Context) {
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
	var input UpdateUserInput
	if err := ctx.BindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUser := models.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		ImageUrl: input.ImageUrl,
	}
	if err := db.Model(&oldUser).Updates(&updateUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success update current user data", "data": updateUser})
}
