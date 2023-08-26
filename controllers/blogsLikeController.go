package controllers

import (
	"blogspot-project/models"
	"blogspot-project/utils/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const STATUS_DISLIKE = 0
const STATUS_LIKE = 1

// LikePostController godoc
// @Summary Create Category
// @Description create new category for post.
// @Tags Like
// @Param id path string true "Post id"
// @Param status path string true "status 0/1 (dislike or like)"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /post/{id}/like/{status} [post]
func LikePostController(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	user_id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int64 = 0
	db.Model(&models.UserLikePost{}).Where("post_id = ? AND user_id = ?", ctx.Param("id"), user_id).Count(&count)
	post_id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	status, err := strconv.Atoi(ctx.Param("status"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var statusLike int = 0
	if status == STATUS_LIKE {
		statusLike = STATUS_LIKE
	} else {
		statusLike = STATUS_DISLIKE
	}

	// if empty just add new database user like list else update status to like
	if count == 0 {
		userLikePost := models.UserLikePost{
			UserID: user_id,
			PostID: uint(post_id),
			Status: uint(statusLike),
		}

		if err := db.Create(&userLikePost).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		oldUserLikePost := models.UserLikePost{}
		if err := db.Where("post_id = ? and user_id = ?", ctx.Param("id"), user_id).Take(&oldUserLikePost).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&oldUserLikePost).Update("status", status).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// update like counts and dislike counts
	var dislike_count int64 = 0
	db.Model(&models.UserLikePost{}).Where("post_id = ? AND status = ?", ctx.Param("id"), STATUS_DISLIKE).Count(&dislike_count)

	var like_count int64 = 0
	db.Model(&models.UserLikePost{}).Where("post_id = ? AND status = ?", ctx.Param("id"), STATUS_LIKE).Count(&like_count)

	post := models.Post{}
	if err := db.Table("posts").Where("id = ?", post_id).Take(&post).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&post).Select("post_like_count", "post_dislike_count").Updates(models.Post{PostLikeCount: uint(like_count), PostDislikeCount: uint(dislike_count)}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success like blog post"})
}

// GetListUserLikePost godoc
// @Summary Get all User like based on blog post id.
// @Description Get all users who likes comment in blog post based on id.
// @Tags Post
// @Produce json
// @Param id path string true "Post id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id}/user-likes [get]
func GetListUserLikePost(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	userLikePost := []models.UserLikePost{}
	if err := db.Where("post_id = ? and status = ?", ctx.Param("id"), STATUS_LIKE).Find(&userLikePost).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	listOfUsers := []models.UserResponse{}
	for _, item := range userLikePost {
		user := models.User{}
		if err := db.Table("users").Where("id = ?", ctx.Param("id"), item.UserID).Find(&user).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		listOfUsers = append(listOfUsers, models.UserResponse{
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			ImageUrl: user.ImageUrl,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success get all user like blog post", "data": listOfUsers})
}

// GetListUserDislikePost godoc
// @Summary Get all User dislike blog.
// @Description Get all users who dislikes in blog post based on id.
// @Tags Like
// @Produce json
// @Param id path string true "Post id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id}/user-dislikes [get]
func GetListUserDislikePost(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	userLikePost := []models.UserLikePost{}
	if err := db.Where("post_id = ? and status = ?", ctx.Param("id"), STATUS_DISLIKE).Find(&userLikePost).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	listOfUsers := []models.UserResponse{}
	for _, item := range userLikePost {
		user := models.User{}
		if err := db.Table("users").Where("id = ?", ctx.Param("id"), item.UserID).Find(&user).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		listOfUsers = append(listOfUsers, models.UserResponse{
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			ImageUrl: user.ImageUrl,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success get all user dislike blog post", "data": listOfUsers})
}
