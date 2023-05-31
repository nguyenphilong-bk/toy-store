package controllers

import (
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
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
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	id, err := categoryModel.Create(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Category could not be created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category created successfully", "data": map[string]uuid.UUID{"id": id}})
}

// All ...
func (ctrl CategoryController) All(c *gin.Context) {
	results, err := categoryModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get categorys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "message": "Get categories successfully"})
}

// One ...
func (ctrl CategoryController) One(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	data, err := categoryModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Get category successfully"})
}

// Update ...
func (ctrl CategoryController) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	var form forms.CreateCategoryForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := categoryForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	err := categoryModel.Update(id, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Category could not be updated"})
		return
	}

	data, err := categoryModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Update category successfully"})
}

// Delete ...
func (ctrl CategoryController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	err := categoryModel.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Category could not be deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
