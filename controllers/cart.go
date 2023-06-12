package controllers

import (
	"fmt"
	"toy-store/common"
	"toy-store/forms"
	"toy-store/services"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CartController ...
type CartController struct{}

var cartService = new(services.CartService)
var cartForm = new(forms.CartForm)

// Create ...
func (ctrl CartController) Update(c *gin.Context) {
	userID := GetUserID(c)
	var form forms.CreateCartForm

	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := cartForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	cart, err := cartService.UpdateCart(userID, form.ProductIDs)
	if err != nil {
		fmt.Println("error when updating cart " + cart.ID.String())

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
	}


	c.JSON(http.StatusOK, gin.H{"message": "cart created successfully", "data": "fuck", "code": common.CODE_SUCCESS})
}

// All ...
// func (ctrl CartController) All(c *gin.Context) {
// 	results, err := cartModel.All()
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Could not get carts", "data": nil, "code": common.CODE_FAILURE},)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": results, "message": "Get carts successfully", "code": common.CODE_SUCCESS})
// }

// // One ...
// func (ctrl CartController) MyCart(c *gin.Context) {
// 	// id := c.Param("id")
// 	id := GetUserID(c)

// 	fmt.Println(id)
// 	// data, err := cartModel.One(id)
// 	// if err != nil {
// 	// 	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Cart not found", "code": common.CODE_FAILURE, "data": nil})
// 	// 	return
// 	// }

// 	c.JSON(http.StatusOK, gin.H{"data": "alo", "message": "Get cart successfully", "code": common.CODE_SUCCESS})
// }

// // Update ...
// func (ctrl CartController) Update(c *gin.Context) {
// 	id := c.Param("id")

// 	if id == "" {
// 		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE})
// 		return
// 	}

// 	var form forms.CreateCartForm

// 	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
// 		message := cartForm.Create(validationErr)
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "data": nil, "code": common.CODE_FAILURE})
// 		return
// 	}

// 	err := cartModel.Update(id, form)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Cart could not be updated", "code": common.CODE_FAILURE, "data": nil})
// 		return
// 	}

// 	data, err := cartModel.One(id)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Cart not found", "code": common.CODE_FAILURE, "data": nil})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": data, "message": "Update cart successfully", "code": common.CODE_SUCCESS})
// }

// // Delete ...
// func (ctrl CartController) Delete(c *gin.Context) {
// 	id := c.Param("id")

// 	if id == "" {
// 		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid parameter", "code": common.CODE_FAILURE})
// 		return
// 	}

// 	err := cartModel.Delete(id)
// 	if err != nil {
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Cart could not be deleted", "code": common.CODE_FAILURE})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted", "code": common.CODE_SUCCESS})
// }
