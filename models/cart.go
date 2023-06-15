package models

import (
	"errors"
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

type CartDetail struct {
	CartID        uuid.UUID `json:"cart_id" db:"cart_id"`
	ProductID     uuid.UUID `json:"product_id" db:"product_id"`
	OrderQuantity int       `json:"order_quantity" db:"order_quantity"`
	Price         float64   `json:"price"`
	Name          string    `json:"name"`
	Origin        string    `json:"origin"`
	ImageUrl      string    `json:"image_url" db:"image_url"`
}

type ProductItem struct {
	ProductID     uuid.UUID `json:"product_id"`
	Price         float64   `json:"price"`
	Name          string    `json:"name"`
	Origin        string    `json:"origin"`
	OrderQuantity int       `json:"order_quantity"`
	ImageUrl      string    `json:"image_url"`
}

type CartDetailResponse struct {
	CartID   uuid.UUID     `json:"cart_id"`
	Total    float64       `json:"total"`
	Products []ProductItem `json:"products"`
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
func (m CartModel) Detail(id string) (detailResponse CartDetailResponse, err error) {
	details := []CartDetail{}
	_, err = db.GetDB().Select(&details, `select cp.cart_id, cp.product_id, cp.order_quantity, p.price, p."name", p.origin, p.image_url 
	FROM cart_products cp
	LEFT JOIN products p ON cp.product_id = p.id 
	where cp.cart_id = $1 AND cp.deleted_at IS NULL`, id)
	if len(details) == 0 {
		detailResponse.CartID = uuid.MustParse(id)
		detailResponse.Total = 0
		return
	}

	detailResponse.CartID = details[0].CartID
	hash := map[uuid.UUID]int{}
	total := 0.
	for _, detail := range details {
		_, ok := hash[detail.ProductID]
		if !ok {
			detailResponse.Products = append(detailResponse.Products, ProductItem{
				ProductID:     detail.ProductID,
				Price:         detail.Price,
				Name:          detail.Name,
				Origin:        detail.Origin,
				OrderQuantity: detail.OrderQuantity,
				ImageUrl:      detail.ImageUrl,
			})
			hash[detail.ProductID] = len(detailResponse.Products) - 1
		} else {
			detailResponse.Products[hash[detail.ProductID]].OrderQuantity += detail.OrderQuantity
		}
		total += detail.Price * float64(detail.OrderQuantity)
	}
	detailResponse.Total = total
	return detailResponse, err
}

// One ...
func (m CartModel) GetCartByUserID(userID string) (cart Cart, err error) {
	err = db.GetDB().SelectOne(&cart, "SELECT id, user_id, created_at, updated_at, deleted_at FROM public.carts WHERE user_id = $1 and deleted_at IS NULL", userID)

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
	operation, err := db.GetDB().Exec(" FROM public.cart_products WHERE id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}
