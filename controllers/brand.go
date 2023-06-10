package controllers

import (
	"github.com/Massad/gin-boilerplate/common"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/google/uuid"

	"net/http"

	"github.com/gin-gonic/gin"
)

// BrandController ...
type BrandController struct{}

var brandModel = new(models.BrandModel)
var brandForm = new(forms.BrandForm)

// Create ...
func (ctrl BrandController) Create(c *gin.Context) {
	var form forms.CreateBrandForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := brandForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	id, err := brandModel.Create(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Brand could not be created", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "brand created successfully", "data": map[string]uuid.UUID{"id": id}, "code": common.CODE_SUCCESS})
}

// All ...
func (ctrl BrandController) All(c *gin.Context) {
	results, err := brandModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get brands", "data": nil, "code": common.CODE_FAILURE},)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "message": "Get brands successfully", "code": common.CODE_SUCCESS})
}

// One ...
func (ctrl BrandController) One(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := brandModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Brand not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Get brand successfully", "code": common.CODE_SUCCESS})
}

// Update ...
func (ctrl BrandController) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE})
		return
	}

	var form forms.CreateBrandForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := brandForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "data": nil, "code": common.CODE_FAILURE})
		return
	}

	err := brandModel.Update(id, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Brand could not be updated", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := brandModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Brand not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Update brand successfully", "code": common.CODE_SUCCESS})
}

// Delete ...
func (ctrl BrandController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE})
		return
	}

	err := brandModel.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Brand could not be deleted", "code": common.CODE_FAILURE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brand deleted", "code": common.CODE_SUCCESS})
}
