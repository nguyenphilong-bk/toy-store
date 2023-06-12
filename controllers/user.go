package controllers

import (
	"fmt"
	"toy-store/common"
	"toy-store/forms"
	"toy-store/models"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserController ...
type UserController struct{}

var userModel = new(models.UserModel)
var userForm = new(forms.UserForm)

// getUserID ...
func getUserID(c *gin.Context) (userID string) {
	//MustGet returns the value for the given key if it exists, otherwise it panics.
	return c.MustGet("userID").(string)
}

// Login ...
func (ctrl UserController) Login(c *gin.Context) {
	var loginForm forms.LoginForm

	if validationErr := c.ShouldBindJSON(&loginForm); validationErr != nil {
		message := userForm.Login(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message, "code": common.CODE_FAILURE})
		return
	}

	_, token, err := userModel.Login(loginForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Invalid login details", "code": common.CODE_FAILURE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged in", "data": token, "code": common.CODE_SUCCESS})
}

// Register ...
func (ctrl UserController) Register(c *gin.Context) {
	var registerForm forms.RegisterForm

	if validationErr := c.ShouldBindJSON(&registerForm); validationErr != nil {
		message := userForm.Register(validationErr)
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": message, "code": common.CODE_FAILURE})
		return
	}

	user, err := userModel.Register(registerForm)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": err.Error(), "code": common.CODE_FAILURE})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully registered", "data": user, "code": common.CODE_SUCCESS})
}

// Logout ...
func (ctrl UserController) Logout(c *gin.Context) {

	_, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in", "code": common.CODE_FAILURE})
		return
	}

	// deleted, delErr := authModel.DeleteAuth(au.AccessUUID)
	// if delErr != nil || deleted == 0 { //if any goes wrong
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid request"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out", "code": common.CODE_SUCCESS})
}

func (ctrl UserController) Me(c *gin.Context) {
	_, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not logged in", "code": common.CODE_FAILURE})
		return
	}

	userID := getUserID(c)

	// cartService.Create(userID, )
	user, err := userModel.One(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error when getting user information", "code": common.CODE_FAILURE})
		return
	}

	cart, _ := cartModel.GetCartByUserID(userID)
	if cart.ID == uuid.Nil {
		cartID, err := cartModel.Create(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error when creating cart for user", "code": common.CODE_FAILURE})
			return
		}

		cart.ID = uuid.MustParse(cartID)
	}

	user.Cart, err = cartModel.Detail(cart.ID.String())
	if err != nil {
		fmt.Println("error when getting cart detail")
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get user information successfully", "data": user, "code": common.CODE_SUCCESS})
}
