package models

import (
	"errors"

	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `db:"id, primarykey" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	BaseModel
}

// ArticleModel ...
type CategoryModel struct{}

// Create ...
func (m CategoryModel) Create(form forms.CreateCategoryForm) (id uuid.UUID, err error) {
	err = db.GetDB().QueryRow("INSERT INTO public.categories(name, description) VALUES($1, $2) RETURNING id", form.Name, form.Description).Scan(&id)
	return id, err
}

// One ...
func (m CategoryModel) One(id string) (category Category, err error) {
	err = db.GetDB().SelectOne(&category, "SELECT * FROM public.categories as b WHERE b.id=$1 AND deleted_at IS NULL LIMIT 1", id)
	return category, err
}

// All ...
func (m CategoryModel) All() (categories []Category, err error) {
	_, err = db.GetDB().Select(&categories, "SELECT * FROM public.categories where deleted_at is null")
	return categories, err
}

// Update ...
func (m CategoryModel) Update(id string, form forms.CreateCategoryForm) (err error) {
	//METHOD 1
	//Check the article by ID using this way
	// _, err = m.One(userID, id)
	// if err != nil {
	// 	return err
	// }

	operation, err := db.GetDB().Exec("UPDATE public.categories SET name=$2, description=$3 WHERE id=$1", id, form.Name, form.Description)
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
func (m CategoryModel) Delete(id string) (err error) {

	operation, err := db.GetDB().Exec("UPDATE public.categories SET deleted_at = CURRENT_TIMESTAMP where id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}
