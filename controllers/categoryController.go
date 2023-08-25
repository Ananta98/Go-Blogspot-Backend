package controllers

import (
	"blogspot-project/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryInput struct {
	Name string `json:"name"`
}

// CreateNewCategory godoc
// @Summary Create Category
// @Description create new category for post.
// @Tags Category
// @Param Body body CategoryInput true "json body to create new category"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /category [post]
func CreateNewCategory(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var input CategoryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newCategory := models.Category{
		Name: input.Name,
	}
	var createdCategory models.Category
	if err := db.Create(&newCategory).Find(&createdCategory).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Create New Category Success", "data": createdCategory})
}

// UpdateCategory godoc
// @Summary Update Category.
// @Description Update existing category based on category id.
// @Tags Category
// @Produce json
// @Param id path string true "Category id"
// @Param Body body CategoryInput true "the body to update existing category"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /category/{id} [patch]
func UpdateCategory(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var oldCategory models.Category
	if err := db.Where("id = ?", ctx.Param("id")).Take(&oldCategory).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var input CategoryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var updatedCategoryResult models.Category
	updatedCategory := models.Category{
		Name: input.Name,
	}
	if err := db.Model(&oldCategory).Updates(&updatedCategory).Find(&updatedCategoryResult).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Update Category Success", "data": updatedCategoryResult})
}

// DeleteCategory godoc
// @Summary Delete existing category.
// @Description Delete existing category by id.
// @Tags Category
// @Produce json
// @Param id path string true "Category id"
// @Success 200 {object} map[string]interface{}
// @Router /category/{id} [delete]
func DeleteCategory(ctx *gin.Context) {
	db := ctx.MustGet("db").(*gorm.DB)
	var category models.Category
	if err := db.Where("id = ?", ctx.Param("id")).Take(&category).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Delete Category Success"})
}

// GetListUsers godoc
// @Summary Get all Category list.
// @Description Get all categories.
// @Tags Category
// @Produce json
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Success 200 {object} map[string]interface{}
// @Router /category [get]
func GetListCategories(ctx *gin.Context) {
	var categories []models.Category
	db := ctx.MustGet("db").(*gorm.DB)
	if err := db.Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Get list category success", "data": categories})
}
