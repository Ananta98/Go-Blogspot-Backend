package controllers

import (
	"blogspot-project/models"
	"blogspot-project/utils/token"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LikeCommentController godoc
// @Summary Like comment
// @Description like comment in existing post.
// @Tags Like
// @Param id path string true "Comment id"
// @Param status path string true "status 0/1 (dislike or like)"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /post/comment/{id}/like/{status} [post]
func LikeCommentController(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	user_id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var count int64 = 0
	db.Model(&models.UserLikeComment{}).Where("comment_id = ? AND user_id = ?", ctx.Param("id"), user_id).Count(&count)
	comment_id, err := strconv.Atoi(ctx.Param("id"))
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
		UserLikeComment := models.UserLikeComment{
			UserID:    user_id,
			CommentID: uint(comment_id),
			Status:    uint(statusLike),
		}

		if err := db.Create(&UserLikeComment).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		oldUserLikeComment := models.UserLikeComment{}
		if err := db.Where("comment_id = ? and user_id = ?", ctx.Param("id"), user_id).Take(&oldUserLikeComment).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Model(&oldUserLikeComment).Update("status", status).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	// update like counts and dislike counts
	var dislike_count int64 = 0
	db.Model(&models.UserLikeComment{}).Where("comment_id = ? AND status = ?", ctx.Param("id"), STATUS_DISLIKE).Count(&dislike_count)

	var like_count int64 = 0
	db.Model(&models.UserLikeComment{}).Where("comment_id = ? AND status = ?", ctx.Param("id"), STATUS_LIKE).Count(&like_count)

	comment := models.Comment{}
	if err := db.Table("comments").Where("id = ?", comment_id).Take(&comment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&comment).Select("comment_like_count", "comment_dislike_count").Updates(models.Comment{CommentLikeCount: uint(like_count), CommentDislikeCount: uint(dislike_count)}).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Success like or dislike comment post"})
}

// GetListUserLikeComment godoc
// @Summary Get all User likes based on comment blog post id.
// @Description Get all users who likes comment in blog post based on id.
// @Tags Like
// @Produce json
// @Param id path string true "Comment id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/comment/{id}/user-likes [get]
func GetListUserLikeComment(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	userLikeComment := []models.UserLikeComment{}
	if err := db.Table("user_like_comments").Where("comment_id = ? and status = ?", ctx.Param("id"), STATUS_LIKE).Find(&userLikeComment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	listOfUsers := []models.UserResponse{}
	for _, item := range userLikeComment {
		user := models.User{}
		if err := db.Table("users").Where("id = ?", item.UserID).Find(&user).Error; err != nil {
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

// GetListUserDislikeComment godoc
// @Summary Get all User dislike comment blog.
// @Description Get all users who dislikes comment in blog post based on id.
// @Tags Like
// @Produce json
// @Param id path string true "Comment id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/comment/{id}/user-dislikes [get]
func GetListUserDislikeComment(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	userLikeComment := []models.UserLikeComment{}
	if err := db.Table("user_like_comments").Where("comment_id = ? and status = ?", ctx.Param("id"), STATUS_DISLIKE).Find(&userLikeComment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	listOfUsers := []models.UserResponse{}
	for _, item := range userLikeComment {
		user := models.User{}
		if err := db.Table("users").Where("id = ?", item.UserID).Find(&user).Error; err != nil {
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
	ctx.JSON(http.StatusOK, gin.H{"message": "Success get all user dislike comment blog post", "data": listOfUsers})
}
