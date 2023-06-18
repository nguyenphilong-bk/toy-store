package models

import (
	"github.com/go-gorp/gorp"
	"github.com/google/uuid"
)

// Article ...

type OrderProduct struct {
	OrderID        uuid.UUID `db:"order_id" json:"order_id"`
	ProductID     uuid.UUID `db:"product_id" json:"product_id"`
	ID            uuid.UUID `db:"id, primarykey" json:"id"`
	OrderQuantity int       `db:"order_quantity" json:"order_quantity"`
	BaseModel
}

// ArticleModel ...
type OrderProductModel struct{}

// Create ...
func (m OrderProductModel) BatchCreate(tx *gorp.Transaction, orderID string, cartInfo CartDetailResponse) (ids []string, err error) {
	for _, product := range cartInfo.Products {
		var returnID string
		err = tx.QueryRow("INSERT INTO public.product_orders(order_id, product_id, order_quantity, price) VALUES($1, $2, $3, $4) RETURNING id", orderID, product.ProductID, product.OrderQuantity, product.Price).Scan(&returnID)
		if err != nil {
			return ids, err
		}
		ids = append(ids, returnID)
	}

	return ids, nil
}

// One ...
// func (m OrderModel) One(id string) (order Order, err error) {
// 	// err = db.GetDB().Raw("SELECT * FROM public.orders as b WHERE b.id=? AND deleted_at IS NULL LIMIT 1", id).Scan(&order).Error

// 	// if order.ID == uuid.Nil {
// 	// 	return order, errors.New("not found")
// 	// }

// 	return order, err
// }

// // One ...
// func (m OrderModel) GetOrderByUserID(userID string) (order Order, err error) {
// 	err = db.GetDB().SelectOne(&order, "SELECT * FROM public.orders WHERE user_id = $1 and deleted_at IS NULL", userID)

// 	return order, err
// }

// // All ...
// func (m OrderModel) All() (orders []Order, err error) {
// 	// rows, err := db.GetDB().Raw("SELECT * FROM public.orders where deleted_at is null").Rows()
// 	// defer rows.Close()

// 	// for rows.Next() {
// 	// 	// ScanRows scan a row into user
// 	// 	var order Order
// 	// 	db.GetDB().ScanRows(rows, &order)
// 	// 	orders = append(orders, order)
// 	// 	// do something
// 	// }
// 	return orders, err
// }

// Update ...
func (m OrderProductModel) Delete(tx *gorp.Transaction, order_id string) (err error) {
	_, err = tx.Exec("UPDATE public.order_products SET deleted_at = current_timestamp WHERE order_id = $1", order_id)

	return err
}
