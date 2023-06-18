package models

import (
	"errors"
	"toy-store/db"

	"github.com/google/uuid"
)

// Article ...

type Order struct {
	UserID   uuid.UUID `db:"user_id" json:"user_id"`
	Products []Product `json:"products"`
	ID       uuid.UUID `db:"id, primarykey" json:"id"`
	BaseModel
}

type OrderDetail struct {
	OrderID       uuid.UUID `json:"order_id" db:"order_id"`
	ProductID     uuid.UUID `json:"product_id" db:"product_id"`
	OrderQuantity int       `json:"order_quantity" db:"order_quantity"`
	Price         float64   `json:"price"`
	Name          string    `json:"name"`
	Origin        string    `json:"origin"`
	ImageUrl      string    `json:"image_url" db:"image_url"`
}

type OrderDetailResponse struct {
	OrderID     uuid.UUID     `json:"order_id"`
	Total       float64       `json:"total"`
	Products    []ProductItem `json:"products"`
	RedirectURL string        `json:"redirect_url,omitempty"`
}

// ArticleModel ...
type OrderModel struct{}

// Create ...
func (m OrderModel) Create(userID string, total float64, address string) (orderID string, err error) {
	err = db.GetDB().QueryRow("INSERT INTO public.orders(user_id, total, address_shipping) VALUES($1, $2, $3) RETURNING id", userID, total, address).Scan(&orderID)

	if err != nil {
		return orderID, err
	}

	return orderID, nil
}

// One ...
func (m OrderModel) Detail(id string) (detailResponse OrderDetailResponse, err error) {
	details := []OrderDetail{}
	_, err = db.GetDB().Select(&details, `select po.order_id, po.product_id, po.order_quantity, po.price, p."name", p.origin, p.image_url 
	FROM product_orders po
	LEFT JOIN products p ON po.product_id = p.id 
	where po.order_id = $1 AND po.deleted_at IS NULL`, id)
	if len(details) == 0 {
		detailResponse.OrderID = uuid.MustParse(id)
		detailResponse.Total = 0
		return
	}

	detailResponse.OrderID = details[0].OrderID
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
func (m OrderModel) GetOrderByUserID(userID string) (order Order, err error) {
	err = db.GetDB().SelectOne(&order, "SELECT id, user_id, created_at, updated_at, deleted_at FROM public.orders WHERE user_id = $1 and deleted_at IS NULL", userID)

	return order, err
}

// All ...
func (m OrderModel) All() (orders []Order, err error) {
	// rows, err := db.GetDB().Raw("SELECT * FROM public.orders where deleted_at is null").Rows()
	// defer rows.Close()

	// for rows.Next() {
	// 	// ScanRows scan a row into user
	// 	var order Order
	// 	db.GetDB().ScanRows(rows, &order)
	// 	orders = append(orders, order)
	// 	// do something
	// }
	return orders, err
}

// Update ...
// func (m OrderModel) Update(id string, form forms.CreateOrderForm) (err error) {
// 	//METHOD 1
// 	//Check the article by ID using this way
// 	// _, err = m.One(userID, id)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// err = db.GetDB().Exec("UPDATE public.orders SET name=? WHERE id=?", form.UserID, id).Error

// 	return err
// }

// Delete ...
func (m OrderModel) Delete(id string) (err error) {
	operation, err := db.GetDB().Exec(" FROM public.order_products WHERE id=$1", id)
	if err != nil {
		return err
	}

	success, _ := operation.RowsAffected()
	if success == 0 {
		return errors.New("no records were deleted")
	}

	return err
}
