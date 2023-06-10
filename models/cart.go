package models

import (
	"errors"

	"toy-store/db"
	"toy-store/forms"
	"github.com/google/uuid"
)

// Article ...

type Cart struct {
	ID   uuid.UUID `db:"id, primarykey" json:"id"`
	Name string    `db:"name" json:"name"`
	BaseModel
}

// ArticleModel ...
type CartModel struct{}

// Create ...
func (m CartModel) Create(form forms.CreateCartForm) (articleID uuid.UUID, err error) {
	var idString string
	err = db.GetDB().Raw("INSERT INTO public.carts(name) VALUES($1) RETURNING id", form.Name).Scan(&idString).Error
	return uuid.MustParse(idString), err
}

// One ...
func (m CartModel) One(id string) (cart Cart, err error) {
	err = db.GetDB().Raw("SELECT * FROM public.carts as b WHERE b.id=? AND deleted_at IS NULL LIMIT 1", id).Scan(&cart).Error

	if cart.ID == uuid.Nil {
		return cart, errors.New("not found")
	}

	return cart, err
}

// All ...
func (m CartModel) All() (carts []Cart, err error) {
	rows, err := db.GetDB().Raw("SELECT * FROM public.carts where deleted_at is null").Rows()
	defer rows.Close()

	for rows.Next() {
		// ScanRows scan a row into user
		var cart Cart
		db.GetDB().ScanRows(rows, &cart)
		carts = append(carts, cart)
		// do something
	}
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

	err = db.GetDB().Exec("UPDATE public.carts SET name=? WHERE id=?", form.Name, id).Error

	return err
}

// Delete ...
func (m CartModel) Delete(id string) (err error) {
	err = db.GetDB().Exec("UPDATE public.carts SET deleted_at = CURRENT_TIMESTAMP where id=?", id).Error
	if err != nil {
		return err
	}

	return err
}
