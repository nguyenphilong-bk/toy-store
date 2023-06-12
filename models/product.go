package models

import (
	"errors"
	"fmt"

	"toy-store/db"
	"toy-store/forms"

	"github.com/go-gorp/gorp"
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
	err = db.GetDB().QueryRow("INSERT INTO public.products(name, description, origin, image_url, price, stock) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		form.Name, form.Description, form.Origin, form.ImageURL, form.Price, form.Stock).Scan(&id)
	return id, err
}

// One ...
func (m ProductModel) One(id string) (product Product, err error) {
	err = db.GetDB().SelectOne(&product, "SELECT * FROM public.products as b WHERE b.id=$1 AND deleted_at IS NULL LIMIT 1", id)
	return product, err
}

// All ...
func (m ProductModel) All() (products []Product, err error) {
	_, err = db.GetDB().Select(&products, "SELECT * FROM public.products where deleted_at is null")
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

	operation, err := db.GetDB().Exec(`UPDATE public.products SET name=$2,
		description=$3,
		origin=$4,
	    image_url=$5,
		price=$6,
		stock=$7
		WHERE id=$1`,
		id,
		form.Name,
		form.Description,
		form.Origin,
		form.ImageURL,
		form.Price,
		form.Stock)
	// form.CateID)
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
func (m ProductModel) Delete(id string) (err error) {

	operation, err := db.GetDB().Exec("UPDATE public.products SET deleted_at = CURRENT_TIMESTAMP where id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}

func (m ProductModel) List(productIDs []string) (products []Product, err error) {
	_, err = db.GetDB().Select(&products, "SELECT * FROM public.products where deleted_at is null and id in (:IDs)", map[string]interface{}{"IDs": productIDs})
	return products, err
}

func (m ProductModel) UpdateStock(tx *gorp.Transaction, form []forms.CartProductItem) (err error) {
	for _, product := range form {
		// checking out of stock
		dbProduct, _ := m.One(product.ProductID)
		if dbProduct.Stock < product.OrderQuantity {
			return fmt.Errorf("Product %v with id %v is out of stock", dbProduct.Name, dbProduct.ID.String())
		}

		// Update stock
		_, err := tx.Exec(`UPDATE public.products SET stock = stock - $1
		WHERE id=$2`,
			product.OrderQuantity,
			product.ProductID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
