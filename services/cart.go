package services

import (
	"toy-store/db"
	"toy-store/forms"
	"toy-store/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var cartModel = new(models.CartModel)
var userModel = new(models.UserModel)
var productModel = new(models.ProductModel)

// var cartForm = new(forms.CartForm)

type CartService struct{}

func (service CartService) Create(userID string, form forms.CreateCartForm) (cart models.Cart, err error) {
	tx := db.GetDB().Begin()
	if err = tx.Error; err != nil {
		return cart, err
	}
	products := []models.Product{}

	if len(form.ProductIDs) > 0 {
		products, err = productModel.List(form.ProductIDs)
		if err != nil && err != gorm.ErrRecordNotFound {
			return cart, err
		}
	}

	// Create new cart
	cart = models.Cart{
		UserID:   uuid.MustParse(userID),
		Products: products,
	}

	cart, err = cartModel.Create(cart)
	if err != nil {
		tx.Rollback()
		return cart, err
	}

	// Get user info
	user, err := userModel.One(userID)
	if err != nil {
		tx.Rollback()
		return cart, err
	}

	err = assignCartToUser(tx, cart, user)
	if err != nil {
		tx.Rollback()
		return cart, err
	}

	return cart, tx.Commit().Error
}

func (service CartService) GetCartByUserID(userID string) (cart models.Cart, err error){
	cart, err = cartModel.GetCartByUserID(userID)
	
	return cart, err
}

func assignCartToUser(tx *gorm.DB, cart models.Cart, user models.User) (err error) {
	user.Cart = cart
	tx.Save(&user)
	if err = tx.Error; err != nil {
		return err
	}

	return nil
}

func (service CartService) UpdateCart(userID string, productIDs []string) (cart models.Cart, err error){
	cart, err = service.GetCartByUserID(userID)
	if err != nil {
		return cart, err
	}

	// cartProducts = []models.CartProduct{}
	for _, productID := range productIDs {
		cartProducts = append(cartProducts, models.CartProduct{
			CartID:    cart.ID,
			ProductID: uuid.MustParse(productID),
		})
	}

	return cartProducts, nil
}
