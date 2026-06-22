package services

import (
	"fmt"

	"backend/api/models"
	"backend/api/repository"
	"backend/pkg/mpesa"
)

// InitiatePayment creates a payment record and triggers an STK push to the user's phone
func InitiatePayment(userID string, req models.InitiatePaymentRequest) (*models.Payment, error) {
	// Fetch order — guests have no userID so fetch by ID only
	var order *models.Order
	var err error
	if userID == "" {
		order, err = repository.GetOrderByIDOnly(req.OrderID)
	} else {
		order, err = repository.GetOrderByID(req.OrderID, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	// Only pending orders can be paid
	if order.Status != "pending" {
		return nil, fmt.Errorf("this order has already been paid or cancelled")
	}

	// Check if a successful payment already exists for this order
	existing, err := repository.GetPaymentByOrderID(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing payment: %w", err)
	}
	if existing != nil && existing.Status == models.PaymentStatusSuccess {
		return nil, fmt.Errorf("this order has already been paid")
	}

	// Normalize phone number to Safaricom format 2547XXXXXXXX
	phone, err := mpesa.NormalizePhone(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	// Create a pending payment record in the database
	// For guests userID is empty string — payments table allows null user_id
	payment, err := repository.CreatePayment(req.OrderID, userID, phone, order.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// Trigger STK push via Daraja
	stkResp, err := mpesa.InitiateSTKPush(phone, order.Total, order.OrderNumber)
	if err != nil {
		// Mark payment as failed if STK push could not be initiated
		repository.MarkPaymentFailed(payment.ID, "STK push failed: "+err.Error(), nil)
		return nil, fmt.Errorf("failed to initiate STK push: %w", err)
	}

	// Save Daraja's checkout request ID and merchant request ID
	err = repository.UpdatePaymentCheckoutIDs(
		payment.ID,
		stkResp.CheckoutRequestID,
		stkResp.MerchantRequestID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment checkout ids: %w", err)
	}

	payment.CheckoutRequestID = stkResp.CheckoutRequestID
	payment.MerchantRequestID = stkResp.MerchantRequestID

	return payment, nil
}

// HandleMpesaCallback processes the callback from Daraja after the user pays
func HandleMpesaCallback(callback models.MpesaCallback) error {
	stk := callback.Body.StkCallback

	// Find the payment record using the checkout request ID
	payment, err := repository.GetPaymentByCheckoutRequestID(stk.CheckoutRequestID)
	if err != nil {
		return fmt.Errorf("failed to fetch payment: %w", err)
	}
	if payment == nil {
		return fmt.Errorf("payment not found for checkout request id: %s", stk.CheckoutRequestID)
	}

	// ResultCode 0 means success, anything else is a failure
	if stk.ResultCode == 0 {
		// Extract M-Pesa receipt number from callback metadata
		mpesaRef := extractCallbackValue(stk.CallbackMetadata.Item, "MpesaReceiptNumber")

		// Mark payment as successful
		if err := repository.MarkPaymentSuccess(stk.CheckoutRequestID, mpesaRef, stk); err != nil {
			return fmt.Errorf("failed to mark payment success: %w", err)
		}

		// Update order status to confirmed
		if err := repository.UpdateOrderStatus(payment.OrderID, "confirmed"); err != nil {
			return fmt.Errorf("failed to confirm order: %w", err)
		}

		// TODO: Send payment confirmation SMS and email
		// sms.SendPaymentConfirmation(payment.Phone, mpesaRef, payment.Amount)

	} else {
		// Payment failed — save the reason from Daraja
		if err := repository.MarkPaymentFailed(stk.CheckoutRequestID, stk.ResultDesc, stk); err != nil {
			return fmt.Errorf("failed to mark payment failed: %w", err)
		}
	}

	return nil
}

// QueryPaymentStatus manually queries Daraja for the payment status
// This is used when the callback is delayed or never arrives
func QueryPaymentStatus(orderID, userID string) (*models.Payment, error) {
	// Fetch the payment record
	payment, err := repository.GetPaymentByOrderID(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("no payment found for this order")
	}

	// If already resolved no need to query Daraja
	if payment.Status != models.PaymentStatusPending {
		return payment, nil
	}

	// Query Daraja for the current status
	queryResp, err := mpesa.QuerySTKStatus(payment.CheckoutRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payment status: %w", err)
	}

	// ResultCode 0 means success
	if queryResp.ResultCode == "0" {
		if err := repository.MarkPaymentSuccess(payment.CheckoutRequestID, "", queryResp); err != nil {
			return nil, fmt.Errorf("failed to mark payment success: %w", err)
		}
		if err := repository.UpdateOrderStatus(payment.OrderID, "confirmed"); err != nil {
			return nil, fmt.Errorf("failed to confirm order: %w", err)
		}
		payment.Status = models.PaymentStatusSuccess
	} else {
		if err := repository.MarkPaymentFailed(payment.CheckoutRequestID, queryResp.ResultDesc, queryResp); err != nil {
			return nil, fmt.Errorf("failed to mark payment failed: %w", err)
		}
		payment.Status = models.PaymentStatusFailed
	}

	return payment, nil
}

// ─── Private Helpers ──────────────────────────────────────────────────────────

// extractCallbackValue extracts a value from the Daraja callback metadata items by name
func extractCallbackValue(items []struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value"`
}, name string) string {
	for _, item := range items {
		if item.Name == name {
			if v, ok := item.Value.(string); ok {
				return v
			}
			return fmt.Sprintf("%v", item.Value)
		}
	}
	return ""
}
