package forms

import (
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// ArticleForm ...
type BrandForm struct{}

// CreateArticleForm ...
type CreateBrandForm struct {
	Name string `form:"name" json:"name" binding:"required,min=1,max=100"`
}

// Name ...
func (f BrandForm) Name(tag string, errMsg ...string) (message string) {
	switch tag {
	case "required":
		if len(errMsg) == 0 {
			return "Please enter the article title"
		}
		return errMsg[0]
	case "min", "max":
		return "Title should be between 1 to 100 characters"
	default:
		return "Something went wrong, please try again later"
	}
}

// Create ...
func (b BrandForm) Create(err error) string {
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

// Update ...
func (b BrandForm) Update(err error) string {
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
