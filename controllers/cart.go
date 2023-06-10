package controllers

import (
	"toy-store/common"
	"toy-store/forms"
	"toy-store/models"
	"github.com/google/uuid"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CartController ...
type CartController struct{}

var cartModel = new(models.CartModel)
var cartForm = new(forms.CartForm)

// Create ...
func (ctrl CartController) Create(c *gin.Context) {
	var form forms.CreateCartForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := cartForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	id, err := cartModel.Create(form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Cart could not be created", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart created successfully", "data": map[string]uuid.UUID{"id": id}, "code": common.CODE_SUCCESS})
}

// All ...
func (ctrl CartController) All(c *gin.Context) {
	results, err := cartModel.All()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get carts", "data": nil, "code": common.CODE_FAILURE},)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results, "message": "Get carts successfully", "code": common.CODE_SUCCESS})
}

// One ...
func (ctrl CartController) One(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := cartModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Cart not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Get cart successfully", "code": common.CODE_SUCCESS})
}

// Update ...
func (ctrl CartController) Update(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE})
		return
	}

	var form forms.CreateCartForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := cartForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "data": nil, "code": common.CODE_FAILURE})
		return
	}

	err := cartModel.Update(id, form)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Cart could not be updated", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	data, err := cartModel.One(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Cart not found", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Update cart successfully", "code": common.CODE_SUCCESS})
}

// Delete ...
func (ctrl CartController) Delete(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE})
		return
	}

	err := cartModel.Delete(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Cart could not be deleted", "code": common.CODE_FAILURE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted", "code": common.CODE_SUCCESS})
}
