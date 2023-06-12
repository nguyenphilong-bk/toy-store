package controllers

import (
	"toy-store/common"
	"toy-store/forms"
	"toy-store/models"

	"github.com/google/uuid"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CategoryController ...
type CategoryController struct{}

var categoryModel = new(models.CategoryModel)
var categoryForm = new(forms.CategoryForm)

// Create ...
func (ctrl CategoryController) Create(c *gin.Context) {
	var form forms.CreateCategoryForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := categoryForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	id, err := categoryModel.Create(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Category could not be created", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category created successfully", "data": map[string]uuid.UUID{"id": id}, "code": common.CODE_SUCCESS})
}

// All ...
func (ctrl CategoryController) All(c *gin.Context) {
	results, err := categoryModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get categorys", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "message": "Get categories successfully", "code": common.CODE_SUCCESS})
}

// One ...
func (ctrl CategoryController) One(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := categoryModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Category not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Get category successfully", "code": common.CODE_SUCCESS})
}

// Update ...
func (ctrl CategoryController) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	var form forms.CreateCategoryForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := categoryForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	err := categoryModel.Update(id, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Category could not be updated", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := categoryModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Category not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Update category successfully", "code": common.CODE_SUCCESS})
}

// Delete ...
func (ctrl CategoryController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	err := categoryModel.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Category could not be deleted", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted", "code": common.CODE_SUCCESS})
}
