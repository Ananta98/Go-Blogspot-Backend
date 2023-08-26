package controllers

import (
	"blogspot-project/models"
	"blogspot-project/utils"
	"blogspot-project/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InputComment struct {
	PostID         uint   `binding:"required" json:"post_id"`
	CommentContent string `binding:"required" json:"comment_content"`
}

type UpdateCommentInput struct {
	CommentContent string `binding:"required" json:"comment_content"`
}

// CreateNewComment godoc
// @Summary Create Comment Blog Post from post id.
// @Description create new comment blog post based on post id.
// @Tags Comment
// @Param Body body InputComment true "json body to create new comment post"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /post/comment [post]
func CreateNewComment(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var input InputComment
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user_id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var post models.Post
	if err := db.Where("id = ?", input.PostID).Take(&post).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newComment := models.Comment{
		UserID:         user_id,
		PostID:         input.PostID,
		CommentContent: input.CommentContent,
	}
	var createdComment models.Comment
	if err := db.Create(&newComment).Last(&createdComment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Create New Comment Success", "data": createdComment})
}

// UpdateComment godoc
// @Summary Update existing comment.
// @Description Update existing comment blog post based on comment id and post id.
// @Tags Comment
// @Produce json
// @Param id path string true "Post id"
// @Param comment_id path string true "Comment id"
// @Param Body body UpdateCommentInput true "json body to update existing comment in existing post"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id}/comment/{comment_id} [patch]
func UpdateComment(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var input UpdateCommentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var oldComment models.Comment
	if err := db.Where("post_id = ? and id = ?", ctx.Param("id"), ctx.Param("comment_id")).Take(&oldComment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedCommentInput := models.Comment{
		UserID:         oldComment.UserID,
		PostID:         oldComment.PostID,
		CommentContent: input.CommentContent,
	}
	updatedComment := models.Comment{}
	if err := db.Model(&oldComment).Updates(&updatedCommentInput).Take(&updatedComment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success update comment", "data": updatedComment})
}

// DeleteComment godoc
// @Summary Delete existing comment blog post.
// @Description Delete existing comment in blog post by post id.
// @Tags Comment
// @Produce json
// @Param id path string true "Post id"
// @Param comment_id path string true "Comment id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id}/comment/{comment_id} [delete]
func DeleteComment(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var comment models.Comment
	if err := db.Where("post_id = ? and id = ?", ctx.Param("id"), ctx.Param("comment_id")).Take(&comment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&comment).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success delete comment"})
}

// GetListComments godoc
// @Summary Get all comments list.
// @Description Get all comments based on post id.
// @Tags Comment
// @Produce json
// @Param id path string true "Post id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param   input_search      query    string     false        "input text for search comment"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id}/comment/ [get]
func GetListComments(ctx *gin.Context) {
	var comments []models.Comment
	db := ctx.MustGet("db").(*gorm.DB)
	limit, offset, err := utils.GetPagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Where("comment_content LIKE ?", "%"+ctx.Query("input_search")+"%").Limit(limit).Offset(offset).Find(&comments).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Get list category success", "data": comments})
}
