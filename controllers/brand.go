package controllers

import (
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
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	id, err := brandModel.Create(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Brand could not be created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "brand created successfully", "data": map[string]uuid.UUID{"id": id}})
}

// All ...
func (ctrl BrandController) All(c *gin.Context) {
	results, err := brandModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Could not get brands"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "message": "Get brands successfully"})
}

// One ...
func (ctrl BrandController) One(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	data, err := brandModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Brand not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Get brand successfully"})
}

// Update ...
func (ctrl BrandController) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	var form forms.CreateBrandForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := brandForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message})
		return
	}

	err := brandModel.Update(id, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Brand could not be updated"})
		return
	}

	data, err := brandModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Brand not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Update brand successfully"})
}

// Delete ...
func (ctrl BrandController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"Message": "Invalid parameter"})
		return
	}

	err := brandModel.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"Message": "Brand could not be deleted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Brand deleted"})
}
