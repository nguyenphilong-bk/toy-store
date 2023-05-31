package forms

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// ArticleForm ...
type ProductForm struct{}

// CreateArticleForm ...
type CreateProductForm struct {
	Name        string `form:"name" json:"name" binding:"required,min=1,max=100"`
	Description string `form:"description" json:"description"`
	Origin      string `form:"origin" json:"origin"`
	// BrandID     string  `form:"brand_id" json:"brand_id" binding:"required"`
	ImageURL string  `form:"image_url" json:"image_url"`
	Price    float64 `form:"price" json:"price" binding:"required,min=1000"`
	Stock    int     `form:"stock" json:"stock" binding:"required,min=0"`
	// CateID      string  `form:"cate_id" json:"cate_id" binding:"required"`
}

// Name ...
func (f ProductForm) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the product name"
		}
		return errMsg[0]
	case "min", "max":
		return "Title should be between 1 to 100 characters"
	default:
		return "Something went wrong, please try again later"
	}
}

// Stock ...
func (f ProductForm) Stock(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the stock for this product"
		}
		return errMsg[0]
	case "min":
		return "Stock must me greater than or equal 0"
	default:
		return "Something went wrong, please try again later"
	}
}

// Stock ...
func (f ProductForm) Price(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the price for this product"
		}
		return errMsg[0]
	case "min":
		return "Price must me greater than or equal 1000"
	default:
		return "Something went wrong, please try again later"
	}
}

// BrandID ...
func (f ProductForm) BrandID(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the brand id"
		}
		return errMsg[0]
	default:
		return "Something went wrong, please try again later"
	}
}

// CateID
func (f ProductForm) CateID(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the cate_id"
		}
		return errMsg[0]
	default:
		return "Something went wrong, please try again later"
	}
}

// Create ...
func (b ProductForm) Create(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Name":
				return b.Name(err.Tag())
			// case "BrandID":
			// 	return b.BrandID(err.Tag())
			// case "CateID":
			// 	return b.CateID(err.Tag())
			case "Price":
				return b.Price(err.Tag())
			case "Stock":
				return b.Stock(err.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}

// Update ...
func (b ProductForm) Update(err error) string {
	switch err.(type) {
	case validator.ValidationErrors:

		if _, ok := err.(*json.UnmarshalTypeError); ok {
			return "Something went wrong, please try again later"
		}

		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Name" {
				return b.Name(err.Tag())
			}
		}

	default:
		return "Invalid request"
	}

	return "Something went wrong, please try again later"
}
