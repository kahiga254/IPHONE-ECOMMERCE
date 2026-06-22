package services

import (
	"fmt"

	"backend/api/models"
	"backend/api/repository"
	"backend/pkg/database"

	"github.com/google/uuid"
)

// shippingFee returns the shipping fee based on the county
func shippingFee(county string) float64 {
	switch county {
	case "Nairobi":
		return 200
	case "Kiambu", "Machakos", "Kajiado":
		return 300
	default:
		return 500
	}
}

// CreateOrder validates the request, calculates totals and creates the order
func CreateOrder(userID string, req models.CreateOrderRequest) (*models.Order, error) {
	// Verify the address belongs to the user
	address, err := repository.GetAddressByID(req.AddressID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch address: %w", err)
	}
	if address == nil {
		return nil, fmt.Errorf("address not found")
	}

	// Calculate subtotal by fetching each variant price from the DB
	subtotal := 0.0
	for _, item := range req.Items {
		var price float64
		var stock int

		err := database.DB.QueryRow(`
			SELECT price, stock FROM product_variants WHERE id = $1`,
			item.VariantID,
		).Scan(&price, &stock)
		if err != nil {
			return nil, fmt.Errorf("variant %s not found", item.VariantID)
		}
		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for variant %s", item.VariantID)
		}

		subtotal += price * float64(item.Quantity)
	}

	// Calculate shipping fee based on delivery county
	fee := shippingFee(address.County)

	// Calculate total
	total := subtotal + fee

	// Create the order
	order, err := repository.CreateOrder(userID, req, subtotal, fee, total)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// TODO: Send order confirmation SMS and email
	// sms.SendOrderConfirmation(address.Phone, order.OrderNumber, total)
	// email.SendOrderConfirmation(userEmail, order)

	return order, nil
}

// GetOrderByID fetches a single order belonging to a user
func GetOrderByID(orderID, userID string) (*models.Order, error) {
	order, err := repository.GetOrderByID(orderID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

// GetMyOrders fetches all orders for a given user with pagination
func GetMyOrders(userID string, q models.OrderFilterQuery) (*models.PaginatedResponse, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}

	orders, total, err := repository.GetOrdersByUserID(userID, q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch orders: %w", err)
	}

	totalPages := (total + q.Limit - 1) / q.Limit

	return &models.PaginatedResponse{
		Data:       orders,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetAllOrders fetches all orders for the admin panel with pagination
func GetAllOrders(q models.OrderFilterQuery) (*models.PaginatedResponse, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}

	orders, total, err := repository.GetAllOrders(q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all orders: %w", err)
	}

	totalPages := (total + q.Limit - 1) / q.Limit

	return &models.PaginatedResponse{
		Data:       orders,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil
}

// UpdateOrderStatus updates an order status — admin only
func UpdateOrderStatus(orderID string, req models.UpdateOrderStatusRequest) error {
	// Fetch order to confirm it exists
	var id string
	err := database.DB.QueryRow(`SELECT id FROM orders WHERE id = $1`, orderID).Scan(&id)
	if err != nil {
		return fmt.Errorf("order not found")
	}

	if err := repository.UpdateOrderStatus(orderID, req.Status); err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// TODO: Notify customer of status change via SMS
	// sms.SendOrderStatusUpdate(phone, orderNumber, req.Status)

	return nil
}

// CancelOrder allows a user to cancel a pending order
func CancelOrder(orderID, userID string) error {
	order, err := repository.GetOrderByID(orderID, userID)
	if err != nil {
		return fmt.Errorf("failed to fetch order: %w", err)
	}
	if order == nil {
		return fmt.Errorf("order not found")
	}

	// Only pending orders can be cancelled
	if order.Status != "pending" {
		return fmt.Errorf("only pending orders can be cancelled")
	}

	// Restore stock for each item
	for _, item := range order.Items {
		_, err := database.DB.Exec(`
			UPDATE product_variants SET stock = stock + $1 WHERE id = $2`,
			item.Quantity, item.VariantID,
		)
		if err != nil {
			return fmt.Errorf("failed to restore stock for variant %s: %w", item.VariantID, err)
		}
	}

	return repository.UpdateOrderStatus(orderID, "cancelled")
}

func CreateGuestOrder(req models.CreateGuestOrderRequest) (*models.Order, error) {
	guestUserID := uuid.New().String() // plain UUID, no "guest-" prefix

	order, err := repository.CreateGuestOrder(guestUserID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create guest order: %w", err)
	}

	return order, nil
}
