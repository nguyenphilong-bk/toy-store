package forms

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// ArticleForm ...
type CartForm struct{}

// CreateArticleForm ...
type CartProductItem struct {
	ProductID     string `json:"product_id"`
	OrderQuantity int    `json:"order_quantity"`
}

type CreateCartForm struct {
	Products []CartProductItem `form:"products" json:"products"`
}

// Name ...
func (f CartForm) UserID(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the article title"
		}
		return errMsg[0]
	case "uuid":
		return "User ID has to be UUID format"
	default:
		return "Something went wrong, please try again later"
	}
}

// Create ...
func (b CartForm) Create(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch {
			case err.Field() == "UserID":
				return b.UserID(err.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}

// Update ...
func (b CartForm) Update(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "UserID" {
				return b.UserID(err.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}
