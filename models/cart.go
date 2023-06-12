package models

import (
	"toy-store/db"
	"toy-store/forms"

	"github.com/google/uuid"
)

// Article ...

type Cart struct {
	UserID   uuid.UUID `db:"user_id" json:"user_id"`
	Products []Product `json:"products"`
	ID       uuid.UUID `db:"id, primarykey" json:"id"`
	BaseModel
}

// ArticleModel ...
type CartModel struct{}

// Create ...
func (m CartModel) Create(userID string) (cartID string, err error) {
	err = db.GetDB().QueryRow("INSERT INTO public.carts(user_id) VALUES($1) RETURNING id", userID).Scan(&cartID)

	if err != nil {
		return cartID, err
	}

	return cartID, nil
}

// One ...
func (m CartModel) One(id string) (cart Cart, err error) {
	// err = db.GetDB().Raw("SELECT * FROM public.carts as b WHERE b.id=? AND deleted_at IS NULL LIMIT 1", id).Scan(&cart).Error

	// if cart.ID == uuid.Nil {
	// 	return cart, errors.New("not found")
	// }

	return cart, err
}

// One ...
func (m CartModel) GetCartByUserID(userID string) (cart Cart, err error) {
	err = db.GetDB().SelectOne(&cart, "SELECT * FROM public.carts WHERE user_id = $1 and deleted_at IS NULL", userID)

	return cart, err
}

// All ...
func (m CartModel) All() (carts []Cart, err error) {
	// rows, err := db.GetDB().Raw("SELECT * FROM public.carts where deleted_at is null").Rows()
	// defer rows.Close()

	// for rows.Next() {
	// 	// ScanRows scan a row into user
	// 	var cart Cart
	// 	db.GetDB().ScanRows(rows, &cart)
	// 	carts = append(carts, cart)
	// 	// do something
	// }
	return carts, err
}

// Update ...
func (m CartModel) Update(id string, form forms.CreateCartForm) (err error) {
	//METHOD 1
	//Check the article by ID using this way
	// _, err = m.One(userID, id)
	// if err != nil {
	// 	return err
	// }

	// err = db.GetDB().Exec("UPDATE public.carts SET name=? WHERE id=?", form.UserID, id).Error

	return err
}

// Delete ...
func (m CartModel) Delete(id string) (err error) {
	// err = db.GetDB().Exec("UPDATE public.carts SET deleted_at = CURRENT_TIMESTAMP where id=?", id).Error
	// if err != nil {
	// 	return err
	// }

	return err
}
