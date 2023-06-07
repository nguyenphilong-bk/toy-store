package models

import (
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
	var idString string
	err = db.GetDB().Raw("INSERT INTO public.categories(name, description) VALUES(?, ?) RETURNING id", form.Name, form.Description).Scan(&idString).Error
	return uuid.MustParse(idString), err
}

// One ...
func (m CategoryModel) One(id string) (category Category, err error) {
	err = db.GetDB().Raw("SELECT * FROM public.categories as b WHERE b.id=$1 AND deleted_at IS NULL LIMIT 1", id).Scan(&category).Error
	return category, err
}

// All ...
func (m CategoryModel) All() (categories []Category, err error) {
	rows, err := db.GetDB().Select(&categories, "SELECT * FROM public.categories where deleted_at is null").Rows()
	defer rows.Close()

	for rows.Next() {
		var category Category

		db.GetDB().ScanRows(rows, &category)
		categories = append(categories, category)
	}

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

	err = db.GetDB().Exec("UPDATE public.categories SET name=?, description=? WHERE id=?", id, form.Name, form.Description).Error
	if err != nil {
		return err
	}

	return err
}

// Delete ...
func (m CategoryModel) Delete(id string) (err error) {

	err = db.GetDB().Exec("UPDATE public.categories SET deleted_at = CURRENT_TIMESTAMP where id=?", id).Error
	if err != nil {
		return err
	}

	return err
}
