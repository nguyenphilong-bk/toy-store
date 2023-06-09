package controllers

import (
	"fmt"
	"net/http"
	"os"

	"toy-store/common"
	"toy-store/forms"
	"toy-store/models"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v4"
)

// AuthController ...
type AuthController struct{}

var authModel = new(models.AuthModel)

// TokenValid ...
func (ctl AuthController) TokenValid(c *gin.Context) {

	tokenAuth, err := authModel.ExtractTokenMetadata(c.Request)
	if err != nil {
		//Token either expired or not valid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please login first", "code": common.CODE_FAILURE, "data": nil})
		return
	}


	//To be called from GetUserID()
	c.Set("userID", tokenAuth.UserID)
}

// Refresh ...
func (ctl AuthController) Refresh(c *gin.Context) {
	var tokenForm forms.Token

	if c.ShouldBindJSON(&tokenForm) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid form", "form": tokenForm, "code": common.CODE_FAILURE})
		c.Abort()
		return
	}

	//verify the token
	token, err := jwt.Parse(tokenForm.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again", "code": common.CODE_FAILURE})
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again", "code": common.CODE_FAILURE})
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		_, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
			return
		}
		userID := claims["user_id"]
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		// 	return
		// }
		//Delete the previous Refresh Token
		// _, _ = authModel.DeleteAuth(refreshUUID)
		// if delErr != nil || deleted == 0 { //if any goes wrong
		// 	c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again"})
		// 	return
		// }

		//Create new pairs of refresh and access tokens
		ts, createErr := authModel.CreateToken(userID.(string))
		if createErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization, please login again", "code": common.CODE_FAILURE})
			return
		}
		//save the tokens metadata to redis
		// saveErr := authModel.CreateAuth(userID.(string), ts)
		// if saveErr != nil {
		// 	c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization, please login again"})
		// 	return
		// }
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusOK, gin.H{"message": "Refreshed token successfully", "data": tokens, "code": common.CODE_SUCCESS})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authorization, please login again", "code": common.CODE_FAILURE})
	}
}
