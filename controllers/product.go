package controllers

import (
	"github.com/google/uuid"
	"toy-store/common"
	"toy-store/forms"
	"toy-store/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

// ProductController ...
type ProductController struct{}

var productModel = new(models.ProductModel)
var productForm = new(forms.ProductForm)

// Create ...
func (ctrl ProductController) Create(c *gin.Context) {
	var form forms.CreateProductForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := productForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	// if form.BrandID != "" {
	// 	_, err := brandModel.One(form.BrandID)
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Brand ID doesn't exist"})
	// 		return
	// 	}
	// }

	// if form.CateID != "" {
	// 	_, err := categoryModel.One(form.CateID)
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Category ID doesn't exist"})
	// 		return
	// 	}
	// }

	id, err := productModel.Create(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Product could not be created", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created successfully", "data": map[string]uuid.UUID{"id": id}, "code": common.CODE_SUCCESS})
}

// All ...
func (ctrl ProductController) All(c *gin.Context) {
	results, err := productModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get products", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "message": "Get categories successfully", "code": common.CODE_SUCCESS})
}

// One ...
func (ctrl ProductController) One(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := productModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Product not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Get product successfully", "code": common.CODE_SUCCESS})
}

// Update ...
func (ctrl ProductController) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	var form forms.CreateProductForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := productForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	err := productModel.Update(id, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Product could not be updated", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := productModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Product not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Update product successfully", "code": common.CODE_SUCCESS})
}

// Delete ...
func (ctrl ProductController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter", "code": common.CODE_FAILURE})
		return
	}

	err := productModel.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Product could not be deleted", "code": common.CODE_FAILURE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted", "code": common.CODE_SUCCESS})
}
