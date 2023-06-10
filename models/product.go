package models

import (
	"errors"

	"toy-store/db"
	"toy-store/forms"
	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `db:"id, primarykey" json:"id"`
	Name        string    `db:"name" json:"name"`
	Origin      string    `db:"origin" json:"origin"`
	Description string    `db:"description" json:"description"`
	ImageURL    string    `db:"image_url" json:"image_url"`
	Price       float64   `db:"price" json:"price"`
	Stock       int       `db:"stock" json:"stock"`
	BaseModel
}

// ArticleModel ...
type ProductModel struct{}

// Create ...
func (m ProductModel) Create(form forms.CreateProductForm) (id uuid.UUID, err error) {
	var idString string
	err = db.GetDB().Raw("INSERT INTO public.products(name, description, origin, image_url, price, stock) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		form.Name, form.Description, form.Origin, form.ImageURL, form.Price, form.Stock).Scan(&idString).Error

	return uuid.MustParse(idString), err
}

// One ...
func (m ProductModel) One(id string) (product Product, err error) {
	err = db.GetDB().Raw("SELECT * FROM public.products as b WHERE b.id=? AND deleted_at IS NULL LIMIT 1", id).Scan(&product).Error

	if product.ID == uuid.Nil {
		return product, errors.New("not found")
	}

	return product, err
}

// All ...
func (m ProductModel) All() (products []Product, err error) {
	rows, err := db.GetDB().Raw("SELECT * FROM public.products where deleted_at is null").Rows()

	defer rows.Close()
	for rows.Next() {
		var product Product
		db.GetDB().ScanRows(rows, &product)
		products = append(products, product)
	}

	return products, err
}

// Update ...
func (m ProductModel) Update(id string, form forms.CreateProductForm) (err error) {
	//METHOD 1
	//Check the article by ID using this way
	// _, err = m.One(userID, id)
	// if err != nil {
	// 	return err
	// }

	err = db.GetDB().Exec(`UPDATE public.products SET name=?,
		description=?,
		origin=?,
	    image_url=?,
		price=?,
		stock=?
		WHERE id=?`,
		form.Name,
		form.Description,
		form.Origin,
		form.ImageURL,
		form.Price,
		form.Stock,
		id).Error
	// form.CateID)
	if err != nil {
		return err
	}

	return err
}

// Delete ...
func (m ProductModel) Delete(id string) (err error) {
	err = db.GetDB().Exec("UPDATE public.products SET deleted_at = CURRENT_TIMESTAMP where id=?", id).Error

	return err
}
