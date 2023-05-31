package models

import (
	"database/sql"
	"errors"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/google/uuid"
)

// Article ...
type BaseModel struct {
	CreatedAt sql.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}
type Brand struct {
	ID   uuid.UUID `db:"id, primarykey" json:"id"`
	Name string    `db:"name" json:"name"`
	BaseModel
}

// ArticleModel ...
type BrandModel struct{}

// Create ...
func (m BrandModel) Create(form forms.CreateBrandForm) (articleID uuid.UUID, err error) {
	err = db.GetDB().QueryRow("INSERT INTO public.brands(name) VALUES($1) RETURNING id", form.Name).Scan(&articleID)
	return articleID, err
}

// One ...
func (m BrandModel) One(id string) (brand Brand, err error) {
	err = db.GetDB().SelectOne(&brand, "SELECT * FROM public.brands as b WHERE b.id=$1 AND deleted_at IS NULL LIMIT 1", id)
	return brand, err
}

// All ...
func (m BrandModel) All() (brands []Brand, err error) {
	_, err = db.GetDB().Select(&brands, "SELECT * FROM public.brands where deleted_at is null")
	return brands, err
}

// Update ...
func (m BrandModel) Update(id string, form forms.CreateBrandForm) (err error) {
	//METHOD 1
	//Check the article by ID using this way
	// _, err = m.One(userID, id)
	// if err != nil {
	// 	return err
	// }

	operation, err := db.GetDB().Exec("UPDATE public.brands SET name=$2 WHERE id=$1", id, form.Name)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()

	if success == 0 {
		return errors.New("updated 0 records")
	}

	return err
}

// Delete ...
func (m BrandModel) Delete(id string) (err error) {

	operation, err := db.GetDB().Exec("UPDATE public.brands SET deleted_at = CURRENT_TIMESTAMP where id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}
