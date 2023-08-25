package controllers

import (
	"blogspot-project/models"
	"blogspot-project/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostInput struct {
	ArticleTitle       string `binding:"required" json:"article_title"`
	ArticleDescription string `binding:"required" json:"article_description"`
	CategoryID         uint   `binding:"required" json:"category_id"`
	ArticleContent     string `binding:"required" json:"article_content"`
}

// CreateNewPost godoc
// @Summary Create Blog Post
// @Description create new blog post.
// @Tags Post
// @Param Body body PostInput true "json body to create new blog post"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /post [post]
func CreateNewPost(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var input PostInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var category models.Category
	if err := db.Where("category_id = ?", input.CategoryID).Take(&category).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCategory := models.Post{
		ArticleTitle:       input.ArticleTitle,
		ArticleDescription: input.ArticleDescription,
		CategoryID:         category.ID,
		ArticleContent:     input.ArticleContent,
	}
	var createdCategory models.Category
	if err := db.Create(&newCategory).Find(&createdCategory).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Create New Category Success", "data": createdCategory})
}

// DeleteCategory godoc
// @Summary Delete existing blog post.
// @Description Delete existing blog post by id.
// @Tags Post
// @Produce json
// @Param id path string true "Post id"
// @Success 200 {object} map[string]interface{}
// @Router /category/{id} [delete]
func DeletePost(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var post models.Post
	if err := db.Where("id = ?", ctx.Param("id")).Take(&post).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&post).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Delete Category Success"})
}

func UpdatePost(ctx *gin.Context) {

}

// GetListBlogs godoc
// @Summary Get all Blog Post list.
// @Description Get all categories.
// @Tags Category
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post [get]
func GetListBlogs(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	limit, offset, err := utils.GetPagination(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var blogs []models.Post
	if err := db.Find(&blogs).Limit(limit).Offset(offset).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Get list category success", "data": blogs})
}
