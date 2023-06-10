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
	var idString string
	err = db.GetDB().Raw("INSERT INTO public.brands(name) VALUES($1) RETURNING id", form.Name).Scan(&idString).Error
	return uuid.MustParse(idString), err
}

// One ...
func (m BrandModel) One(id string) (brand Brand, err error) {
	err = db.GetDB().Raw("SELECT * FROM public.brands as b WHERE b.id=? AND deleted_at IS NULL LIMIT 1", id).Scan(&brand).Error

	if brand.ID == uuid.Nil {
		return brand, errors.New("not found")
	}

	return brand, err
}

// All ...
func (m BrandModel) All() (brands []Brand, err error) {
	rows, err := db.GetDB().Raw("SELECT * FROM public.brands where deleted_at is null").Rows()
	defer rows.Close()

	for rows.Next() {
		// ScanRows scan a row into user
		var brand Brand
		db.GetDB().ScanRows(rows, &brand)
		brands = append(brands, brand)
		// do something
	}
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

	err = db.GetDB().Exec("UPDATE public.brands SET name=? WHERE id=?", form.Name, id).Error

	return err
}

// Delete ...
func (m BrandModel) Delete(id string) (err error) {
	err = db.GetDB().Exec("UPDATE public.brands SET deleted_at = CURRENT_TIMESTAMP where id=?", id).Error
	if err != nil {
		return err
	}

	return err
}
