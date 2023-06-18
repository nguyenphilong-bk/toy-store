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

// OrderController ...
type OrderController struct{}

var orderForm = new(forms.OrderForm)
var orderModel = new(models.OrderModel)
var orderProductModel = new(models.OrderProductModel)

// Update ...
func (ctrl OrderController) Checkout(c *gin.Context) {
	userID := getUserID(c)
	var form forms.CreateOrderForm
	if validationErr := c.ShouldBindJSON(&form); validationErr != nil {
		message := orderForm.Create(validationErr)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message, "code": common.CODE_FAILURE, "data": nil})
		return
	}

	cart, err := cartModel.GetCartByUserID(userID)
	fmt.Println("error when getting cart info:", err)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	cartInfo, err := cartModel.Detail(cart.ID.String())
	if err != nil {
		fmt.Println("error when getting cart detail")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	if len(cartInfo.Products) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "can not create an empty order", "code": common.CODE_FAILURE, "data": nil})
		return
	}

	tx, err := db.GetDB().Begin()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	orderID, err := orderModel.Create(userID, cartInfo.Total, form.AddressShipping)
	if err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	// Create order products record
	_, err = orderProductModel.BatchCreate(tx, orderID, cartInfo)
	if err != nil {
		fmt.Println("error when creating order product for orderID " + orderID + err.Error())
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	err = cartModel.Delete(tx, cart.ID.String())
	if err != nil {
		fmt.Println("error when deleting the old cart: ", err.Error())
		tx.Rollback()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	tx.Commit()

	// Create a new cart for this user
	cartModel.Create(userID)

	orderInfo, err := orderModel.Detail(orderID)
	if err != nil {
		fmt.Println("error when getting order detail info", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": common.CODE_FAILURE, "data": nil})
		return
	}

	orderInfo.RedirectURL = "http://example-redirect-url.com"
	c.JSON(http.StatusOK, gin.H{"message": "order created successfully", "data": orderInfo, "code": common.CODE_SUCCESS})
}