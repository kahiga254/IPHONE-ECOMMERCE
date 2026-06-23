package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/api/models"
	"backend/pkg/database"
)

func CreatePayment(orderID, userID, phone string, amount float64) (*models.Payment, error) {
	var payment models.Payment

	// Use NULL for guest users (empty userID)
	var uid interface{}
	if userID == "" {
		uid = nil
	} else {
		uid = userID
	}

	err := database.DB.QueryRow(`
		INSERT INTO payments (order_id, user_id, amount, phone, provider, status)
		VALUES ($1, $2, $3, $4, 'mpesa', 'pending')
		RETURNING id, order_id, user_id, amount, phone, provider, status, created_at, updated_at`,
		orderID, uid, amount, phone,
	).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount,
		&payment.Phone, &payment.Provider, &payment.Status,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &payment, nil
}

// UpdatePaymentCheckoutIDs saves the checkout request ID and merchant request ID
// returned by Daraja immediately after initiating the STK push
func UpdatePaymentCheckoutIDs(paymentID, checkoutRequestID, merchantRequestID string) error {
	_, err := database.DB.Exec(`
		UPDATE payments
		SET checkout_request_id = $1,
		    merchant_request_id = $2,
		    updated_at          = NOW()
		WHERE id = $3`,
		checkoutRequestID, merchantRequestID, paymentID,
	)
	if err != nil {
		return fmt.Errorf("failed to update checkout ids: %w", err)
	}
	return nil
}

// GetPaymentByCheckoutRequestID fetches a payment using Daraja's checkout request ID
// This is used when processing the M-Pesa callback to find which payment to update
func GetPaymentByCheckoutRequestID(checkoutRequestID string) (*models.Payment, error) {
	var payment models.Payment

	err := database.DB.QueryRow(`
		SELECT id, order_id, user_id, amount, phone, provider,
		       mpesa_reference, checkout_request_id, merchant_request_id,
		       status, failure_reason, created_at, updated_at
		FROM payments WHERE checkout_request_id = $1`, checkoutRequestID,
	).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount,
		&payment.Phone, &payment.Provider, &payment.MpesaReference,
		&payment.CheckoutRequestID, &payment.MerchantRequestID,
		&payment.Status, &payment.FailureReason,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by checkout request id: %w", err)
	}

	return &payment, nil
}

// GetPaymentByOrderID fetches the payment record for a given order
func GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment

	err := database.DB.QueryRow(`
		SELECT id, order_id, user_id, amount, phone, provider,
		       mpesa_reference, checkout_request_id, merchant_request_id,
		       status, failure_reason, created_at, updated_at
		FROM payments WHERE order_id = $1
		ORDER BY created_at DESC LIMIT 1`, orderID,
	).Scan(
		&payment.ID, &payment.OrderID, &payment.UserID, &payment.Amount,
		&payment.Phone, &payment.Provider, &payment.MpesaReference,
		&payment.CheckoutRequestID, &payment.MerchantRequestID,
		&payment.Status, &payment.FailureReason,
		&payment.CreatedAt, &payment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment by order id: %w", err)
	}

	return &payment, nil
}

// MarkPaymentSuccess updates the payment status to success and saves the M-Pesa receipt number
func MarkPaymentSuccess(checkoutRequestID, mpesaReference string, payload interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = database.DB.Exec(`
		UPDATE payments
		SET status          = 'success',
		    mpesa_reference = $1,
		    payload         = $2,
		    updated_at      = NOW()
		WHERE checkout_request_id = $3`,
		mpesaReference, payloadJSON, checkoutRequestID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark payment success: %w", err)
	}
	return nil
}

// MarkPaymentFailed updates the payment status to failed and saves the failure reason
func MarkPaymentFailed(checkoutRequestID, reason string, payload interface{}) error {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	_, err = database.DB.Exec(`
		UPDATE payments
		SET status         = 'failed',
		    failure_reason = $1,
		    payload        = $2,
		    updated_at     = NOW()
		WHERE checkout_request_id = $3`,
		reason, payloadJSON, checkoutRequestID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark payment failed: %w", err)
	}
	return nil
}
