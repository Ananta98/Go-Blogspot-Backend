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
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email" binding:"email"`
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
// @Param   current_page      query    int        false        "current page for pagination"
// @Param   page_size         query    int        false        "page size for pagination"
// @Param   input_search      query    string     false        "input text for search category"
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
	currentUser := models.User{}
	if err := db.Where("id = ?", id).Take(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == models.ADMIN_USER_ROLE {
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
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin can look list users"})
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
	updatedUser := models.User{
		Name:     input.Name,
		Username: input.Username,
		Email:    input.Email,
		ImageUrl: input.ImageUrl,
	}
	if err := db.Model(&oldUser).Updates(&updatedUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success update current user data", "data": updatedUser})
}

// DeleteUser godoc
// @Summary Delete existing user.
// @Description Delete existing user by id.
// @Tags User
// @Produce json
// @Param id path string true "User id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /user/{id} [delete]
func DeleteUser(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	currentUser := models.User{}
	if err := db.Where("id = ?", id).Take(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == models.ADMIN_USER_ROLE {
		user := models.User{}
		if err := db.Where("id = ?", ctx.Param("id")).Take(&user).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Delete(&user).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Delete User Success"})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin can delete user"})
}
