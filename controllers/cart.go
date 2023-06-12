package controllers

import (
	"fmt"
	"net/http"
	"toy-store/common"
	"toy-store/db"
	"toy-store/forms"
	"toy-store/models"

	"github.com/gin-gonic/gin"
)

// CartController ...
type CartController struct{}

var cartForm = new(forms.CartForm)
var cartModel = new(models.CartModel)
var cartProductModel = new(models.CartProductModel)

// Update ...
func (ctrl CartController) Update(c *gin.Context) {
	delete := c.Query("delete")
	userID := getUserID(c)
	var form forms.CreateCartForm
	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := cartForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	cart, err := cartModel.GetCartByUserID(userID)
	fmt.Println("error when getting cart info:", err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
	}

	// TODO: Return stock
	// err = productModel.ReturnStock(tx, cart.ID.String())
	// if err != nil {
	// 	fmt.Println("error when deleting old cart products of id: " + cart.ID.String() + err.Error())
	// 	tx.Rollback()
	// 	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
	// 	return
	// }

	// Create cart products record
	if delete == "true" {
		err = cartProductModel.Delete(tx, cart.ID.String())
		if err != nil {
			fmt.Println("error when deleting old cart products of id: " + cart.ID.String() + err.Error())
			tx.Rollback()
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
			return
		}
	}

	err = productModel.UpdateStock(tx, form.Products)
	if err != nil {
		fmt.Println("error when updating product stock", err.Error())
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	// Create cart products record
	_, err = cartProductModel.BatchCreate(tx, cart.ID.String(), form.Products)
	if err != nil {
		fmt.Println("error when creating cart product for cartID " + cart.ID.String() + err.Error())
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	tx.Commit()

	cartInfo, err := cartModel.Detail(cart.ID.String())
	if err != nil {
		fmt.Println("error when getting cart detail info", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart created successfully", "data": cartInfo, "code": common.CODE_SUCCESS})
}
