package forms

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// ArticleForm ...
type OrderForm struct{}


type CreateOrderForm struct {
	AddressShipping string `form:"address_shipping" json:"address_shipping" binding:"required,min=3,max=1000"`
}

// Name ...
func (f OrderForm) AddressShipping(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the address_shipping field"
		}
		return errMsg[0]
	case "min", "max":
		return "address_shipping should be between 3 to 1000 characters"
	default:
		return "Something went wrong, please try again later"
	}
}

// Create ...
func (b OrderForm) Create(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:
		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch {
			case err.Field() == "AddressShipping":
				return b.AddressShipping(err.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}

// Update ...
// func (b OrderForm) Update(err error) string {
// 	switch err.(type) {
// 	case validator.ValidationErrors:

// 		if _, ok := err.(*json.UnmarshalTypeError); ok {
// 			return "Something went wrong, please try again later"
// 		}

// 		for _, err := range err.(validator.ValidationErrors) {
// 			if err.Field() == "UserID" {
// 				return b.UserID(err.Tag())
// 			}
// 		}

// 	default:
// 		return "Invalid request"
// 	}

// 	return "Something went wrong, please try again later"
// }
