package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"backend/api/models"
	"backend/pkg/database"
)

// GenerateOrderNumber creates a unique human readable order number
func GenerateOrderNumber() string {
	return fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
}

// CreateOrder inserts a new order and its items in a single transaction
func CreateOrder(userID string, req models.CreateOrderRequest, subtotal, shippingFee, total float64) (*models.Order, error) {
	// Begin transaction — order + items must be created together or not at all
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	orderNumber := GenerateOrderNumber()

	// Insert the order
	var order models.Order
	err = tx.QueryRow(`
		INSERT INTO orders (user_id, order_number, status, subtotal, shipping_fee, total, address_id, notes)
		VALUES ($1, $2, 'pending', $3, $4, $5, $6, $7)
		RETURNING id, user_id, order_number, status, subtotal, shipping_fee, discount, total, address_id, notes, created_at, updated_at`,
		userID, orderNumber, subtotal, shippingFee, total, req.AddressID, req.Notes,
	).Scan(
		&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
		&order.Subtotal, &order.ShippingFee, &order.Discount, &order.Total,
		&order.AddressID, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Insert each order item and reduce stock
	for _, item := range req.Items {
		// Fetch current variant price and stock
		var unitPrice float64
		var stock int
		err := tx.QueryRow(`
			SELECT price, stock FROM product_variants WHERE id = $1`, item.VariantID,
		).Scan(&unitPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("variant %s not found", item.VariantID)
		}
		if err != nil {
			return nil, fmt.Errorf("failed to fetch variant: %w", err)
		}

		// Check stock availability
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for variant %s", item.VariantID)
		}

		totalPrice := unitPrice * float64(item.Quantity)

		// Insert order item
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, variant_id, quantity, unit_price, total_price)
			VALUES ($1, $2, $3, $4, $5)`,
			order.ID, item.VariantID, item.Quantity, unitPrice, totalPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to insert order item: %w", err)
		}

		// Reduce stock
		_, err = tx.Exec(`
			UPDATE product_variants SET stock = stock - $1 WHERE id = $2`,
			item.Quantity, item.VariantID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to reduce stock: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit order: %w", err)
	}

	return &order, nil
}

// GetOrderByID fetches a single order with its items and address
func GetOrderByID(orderID, userID string) (*models.Order, error) {
	var order models.Order

	err := database.DB.QueryRow(`
		SELECT id, user_id, order_number, status, subtotal, shipping_fee,
		       discount, total, address_id, notes, created_at, updated_at
		FROM orders WHERE id = $1 AND user_id = $2`,
		orderID, userID,
	).Scan(
		&order.ID, &order.UserID, &order.OrderNumber, &order.Status,
		&order.Subtotal, &order.ShippingFee, &order.Discount, &order.Total,
		&order.AddressID, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	order.Items = getOrderItems(order.ID)
	order.Address = getOrderAddress(order.AddressID)

	return &order, nil
}

// GetOrdersByUserID fetches all orders for a given user with pagination
func GetOrdersByUserID(userID string, q models.OrderFilterQuery) ([]models.Order, int, error) {
	where := "WHERE user_id = $1"
	args := []interface{}{userID}
	idx := 2

	if q.Status != "" {
		where += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, q.Status)
		idx++
	}

	var total int
	database.DB.QueryRow("SELECT COUNT(*) FROM orders "+where, args...).Scan(&total)

	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)

	query := fmt.Sprintf(`
		SELECT id, user_id, order_number, status, subtotal, shipping_fee,
		       discount, total, address_id, notes, created_at, updated_at
		FROM orders %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, idx, idx+1,
	)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch orders: %w", err)
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		rows.Scan(
			&o.ID, &o.UserID, &o.OrderNumber, &o.Status,
			&o.Subtotal, &o.ShippingFee, &o.Discount, &o.Total,
			&o.AddressID, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
		)
		o.Items = getOrderItems(o.ID)
		orders = append(orders, o)
	}

	return orders, total, nil
}

// UpdateOrderStatus updates the status of an order
func UpdateOrderStatus(orderID, status string) error {
	_, err := database.DB.Exec(`
		UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`,
		status, orderID,
	)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}

// GetAllOrders fetches all orders for the admin panel with pagination
func GetAllOrders(q models.OrderFilterQuery) ([]models.Order, int, error) {
	where := "WHERE 1=1"
	args := []interface{}{}
	idx := 1

	if q.Status != "" {
		where += fmt.Sprintf(" AND status = $%d", idx)
		args = append(args, q.Status)
		idx++
	}

	var total int
	database.DB.QueryRow("SELECT COUNT(*) FROM orders "+where, args...).Scan(&total)

	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)

	query := fmt.Sprintf(`
		SELECT id, user_id, order_number, status, subtotal, shipping_fee,
		       discount, total, address_id, notes, created_at, updated_at
		FROM orders %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, idx, idx+1,
	)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch all orders: %w", err)
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		rows.Scan(
			&o.ID, &o.UserID, &o.OrderNumber, &o.Status,
			&o.Subtotal, &o.ShippingFee, &o.Discount, &o.Total,
			&o.AddressID, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
		)
		orders = append(orders, o)
	}

	return orders, total, nil
}

// ─── Private Helpers ──────────────────────────────────────────────────────────

// getOrderItems fetches all items for a given order
func getOrderItems(orderID string) []models.OrderItem {
	rows, err := database.DB.Query(`
		SELECT oi.id, oi.order_id, oi.variant_id, oi.quantity, oi.unit_price, oi.total_price,
		       pv.color, pv.storage, pv.images,
		       p.name, p.slug
		FROM order_items oi
		JOIN product_variants pv ON oi.variant_id = pv.id
		JOIN products p          ON pv.product_id = p.id
		WHERE oi.order_id = $1`, orderID,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	items := []models.OrderItem{}
	for rows.Next() {
		var item models.OrderItem
		var variant models.Variant
		var product models.Product
		var imagesJSON []byte

		rows.Scan(
			&item.ID, &item.OrderID, &item.VariantID, &item.Quantity,
			&item.UnitPrice, &item.TotalPrice,
			&variant.Color, &variant.Storage, &imagesJSON,
			&product.Name, &product.Slug,
		)

		json.Unmarshal(imagesJSON, &variant.Images)
		variant.ID = item.VariantID
		item.Variant = &variant
		items = append(items, item)
	}

	return items
}

// getOrderAddress fetches the delivery address for an order
func getOrderAddress(addressID string) *models.Address {
	var a models.Address
	err := database.DB.QueryRow(`
		SELECT id, user_id, label, full_name, phone, county, town, street, is_default, created_at
		FROM addresses WHERE id = $1`, addressID,
	).Scan(
		&a.ID, &a.UserID, &a.Label, &a.FullName, &a.Phone,
		&a.County, &a.Town, &a.Street, &a.IsDefault, &a.CreatedAt,
	)
	if err != nil {
		return nil
	}
	return &a
}
