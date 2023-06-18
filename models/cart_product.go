package models

import (
	"toy-store/forms"

	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
)

// Article ...

type CartProduct struct {
	CartID        uuid.UUID `db:"cart_id" json:"cart_id"`
	ProductID     uuid.UUID `db:"product_id" json:"product_id"`
	ID            uuid.UUID `db:"id, primarykey" json:"id"`
	OrderQuantity int       `db:"order_quantity" json:"order_quantity"`
	BaseModel
}

// ArticleModel ...
type CartProductModel struct{}

// Create ...
func (m CartProductModel) BatchCreate(tx *gorp.Transaction, cartID string, form []forms.CartProductItem) (ids []string, err error) {
	for _, product := range form {
		var returnID string
		err = tx.QueryRow("INSERT INTO public.cart_products(cart_id, product_id, order_quantity) VALUES($1, $2, $3) RETURNING id", cartID, product.ProductID, product.OrderQuantity).Scan(&returnID)
		if err != nil {
			return ids, err
		}
		ids = append(ids, returnID)
	}

	return ids, nil
}

// One ...
// func (m CartModel) One(id string) (cart Cart, err error) {
// 	// err = db.GetDB().Raw("SELECT * FROM public.carts as b WHERE b.id=? AND deleted_at IS NULL LIMIT 1", id).Scan(&cart).Error

// 	// if cart.ID == uuid.Nil {
// 	// 	return cart, errors.New("not found")
// 	// }

// 	return cart, err
// }

// // One ...
// func (m CartModel) GetCartByUserID(userID string) (cart Cart, err error) {
// 	err = db.GetDB().SelectOne(&cart, "SELECT * FROM public.carts WHERE user_id = $1 and deleted_at IS NULL", userID)

// 	return cart, err
// }

// // All ...
// func (m CartModel) All() (carts []Cart, err error) {
// 	// rows, err := db.GetDB().Raw("SELECT * FROM public.carts where deleted_at is null").Rows()
// 	// defer rows.Close()

// 	// for rows.Next() {
// 	// 	// ScanRows scan a row into user
// 	// 	var cart Cart
// 	// 	db.GetDB().ScanRows(rows, &cart)
// 	// 	carts = append(carts, cart)
// 	// 	// do something
// 	// }
// 	return carts, err
// }

// Update ...
func (m CartProductModel) Delete(tx *gorp.Transaction, cartID string) (err error) {
	_, err = tx.Exec("UPDATE public.cart_products SET deleted_at = current_timestamp WHERE cart_id = $1", cartID)

	return err
}

func (m CartProductModel) DeleteItem(tx *gorp.Transaction, cartID, productID string) (err error) {
	_, err = tx.Exec("UPDATE public.cart_products SET deleted_at = current_timestamp WHERE cart_id = $1 and product_id = $2", cartID, productID)

	return err
}
