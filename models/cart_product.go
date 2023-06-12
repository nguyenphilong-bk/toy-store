package models

import (
	"toy-store/db"

	"github.com/google/uuid"
)

// Article ...

type CartProduct struct {
	CartID    uuid.UUID
	ProductID uuid.UUID
	BaseModel
}

// ArticleModel ...
type CartProductModel struct{}

// Create ...
func (m CartProductModel) Create(cartProducts []CartProduct) ([]CartProduct, error) {
	tmp := db.GetDB().Table("cart_products").Create(&cartProducts)

	if tmp.Error != nil {
		return cartProducts, tmp.Error
	}

	return cartProducts, nil
}
