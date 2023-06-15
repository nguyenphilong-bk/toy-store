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
	Material    string    `db:"material" json:"material"`
	Size        string    `db:"size" json:"size"`
	Barcode     string    `db:"barcode" json:"barcode"`
	BaseModel
}

type ProductDetail struct {
	Product
	CategoryName string `db:"category_name" json:"category_name"`
	CategoryID   string `db:"category_id" json:"category_id"`
	BrandName    string `db:"brand_name" json:"brand_name"`
	BrandID      string `db:"brand_id" json:"brand_id"`
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
func (m ProductModel) All(categoryID, brandID string) (products []ProductDetail, err error) {
	query := `select p.id, p.name, p.origin, p.description, p.image_url, p.price, p.stock, COALESCE(p.material, '') material, COALESCE(p.size,'') size, COALESCE(p.barcode,'') barcode,
				COALESCE(c.name, '') category_name, COALESCE(c.id, '00000000-0000-0000-0000-000000000000') category_id,
				COALESCE(b.name, '') brand_name, COALESCE(b.id, '00000000-0000-0000-0000-000000000000') brand_id
				from products p 
				left join product_brands pb on p.id = pb.product_id 
				left join brands b on pb.brand_id = b.id 
				left join product_categories pc on p.id = pc.product_id 
				left join categories c on c.id = pc.category_id WHERE p.deleted_at is null `

	if categoryID != "" {
		query += fmt.Sprintf("AND c.id = '%s'\n", categoryID)
	}

	if brandID != "" {
		query += fmt.Sprintf("AND b.id = '%s'\n", brandID)
	}

	_, err = db.GetDB().Select(&products, query)

	return products, err
}

// All ...
func (m ProductModel) FindByCategoryID(categoryID string) (products []Product, err error) {
	_, err = db.GetDB().Select(&products, "SELECT * FROM public.products LEFT JOIN ON where deleted_at is null")
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
