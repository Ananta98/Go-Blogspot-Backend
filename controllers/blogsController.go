package controllers

import (
	"blogspot-project/models"
	"blogspot-project/utils"
	"blogspot-project/utils/token"
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

type PostUpdate struct {
	ArticleTitle       string `json:"article_title"`
	ArticleDescription string `json:"article_description"`
	CategoryID         uint   `json:"category_id"`
	ArticleContent     string `json:"article_content"`
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
		var input PostInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var category models.Category
		if err := db.Table("categories").Where("id = ?", input.CategoryID).Take(&category).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user_id, err := token.ExtractTokenID(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		newCategory := models.Post{
			UserID:             user_id,
			ArticleTitle:       input.ArticleTitle,
			ArticleDescription: input.ArticleDescription,
			CategoryID:         category.ID,
			ArticleContent:     input.ArticleContent,
		}
		var createdPost models.Post
		if err := db.Create(&newCategory).Last(&createdPost).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Create New blog Success", "data": createdPost})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin can create new post"})
}

// DeletePost godoc
// @Summary Delete existing blog post.
// @Description Delete existing blog post by id.
// @Tags Post
// @Produce json
// @Param id path string true "Post id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id} [delete]
func DeletePost(ctx *gin.Context) {
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
		var post models.Post
		if err := db.Where("id = ?", ctx.Param("id")).Take(&post).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := db.Delete(&post).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Delete blog Success"})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin can delete post"})
}

// UpdatePost godoc
// @Summary Update existing post.
// @Description Update existing post without update password that have logged in into blog.
// @Tags Post
// @Produce json
// @Param id path string true "Post id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id} [patch]
func UpdatePost(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	user_id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	currentUser := models.User{}
	if err := db.Where("id = ?", user_id).Take(&currentUser).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if currentUser.Role == models.ADMIN_USER_ROLE {
		oldPost := models.Post{}
		if err := db.Where("id = ? AND user_id = ?", ctx.Param("id"), user_id).Take(&oldPost).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var input PostUpdate
		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var category models.Category
		if err := db.Table("categories").Where("id = ?", input.CategoryID).Take(&category).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		updatedPost := models.Post{
			ArticleTitle:       input.ArticleTitle,
			ArticleContent:     input.ArticleContent,
			ArticleDescription: input.ArticleDescription,
			CategoryID:         input.CategoryID,
		}
		if err := db.Model(&oldPost).Updates(&updatedPost).Error; err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Success update blog", "data": updatedPost})
		return
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Only Admin can update post"})
}

// GetListBlogs godoc
// @Summary Get all Blog Post list.
// @Description Get all categories.
// @Tags Post
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Param   input_search      query    string     false        "input text for search blog"
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
	if err := db.Where("article_title LIKE ?", "%"+ctx.Query("input_search")+"%").Limit(limit).Offset(offset).Find(&blogs).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Get list blog success", "data": blogs})
}

// GetDetailPost godoc
// @Summary Get detail post by id
// @Description Get post detail based on post id
// @Tags Post
// @Produce json
// @Param id path string true "Post id"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /post/{id} [get]
func GetDetailPost(ctx *gin.Context) {
	post := models.Post{}
	db := ctx.MustGet("db").(*gorm.DB)
	id, err := token.ExtractTokenID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := db.Where("id = ?", id).Take(&post).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Get blog detail success", "data": post})
}
