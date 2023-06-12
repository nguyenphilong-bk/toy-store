package controllers

import (
	"toy-store/forms"
	"toy-store/models"
)

// CartController ...
type CartController struct{}

var cartForm = new(forms.CartForm)
var cartModel = new(models.CartModel)

// Create ...
// func (ctrl CartController) Update(c *gin.Context) {
// 	userID := getUserID(c)
// 	var form forms.CreateCartForm

// 	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
// 		message := cartForm.Create(validationErr)
// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
// 		return
// 	}

// 	cart, err := cartService.UpdateCart(userID, form.ProductIDs)
// 	if err != nil {
// 		fmt.Println("error when updating cart " + cart.ID.String())

// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "cart created successfully", "data": "fuck", "code": common.CODE_SUCCESS})
// }
