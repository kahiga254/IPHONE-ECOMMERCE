package models

import "time"

type Payment struct {
	ID                string    `json:"id"`
	OrderID           string    `json:"order_id"`
	UserID            string    `json:"user_id"`
	Amount            float64   `json:"amount"`
	Phone             string    `json:"phone"`
	Provider          string    `json:"provider"`
	MpesaReference    string    `json:"mpesa_reference,omitempty"`
	CheckoutRequestID string    `json:"checkout_request_id,omitempty"`
	MerchantRequestID string    `json:"merchant_request_id,omitempty"`
	Status            string    `json:"status"`
	FailureReason     string    `json:"failure_reason,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// InitiatePaymentRequest is what the frontend sends to trigger an STK push
type InitiatePaymentRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Phone   string `json:"phone"    binding:"required"`
}

// StkPushResponse is what Daraja sends back immediately after we call their API
type StkPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

// MpesaCallback is what Daraja sends to our callback URL after the user pays
type MpesaCallback struct {
	Body struct {
		StkCallback struct {
			MerchantRequestID string `json:"MerchantRequestID"`
			CheckoutRequestID string `json:"CheckoutRequestID"`
			ResultCode        int    `json:"ResultCode"`
			ResultDesc        string `json:"ResultDesc"`
			CallbackMetadata  struct {
				Item []struct {
					Name  string      `json:"Name"`
					Value interface{} `json:"Value"`
				} `json:"Item"`
			} `json:"CallbackMetadata"`
		} `json:"stkCallback"`
	} `json:"Body"`
}

// MpesaQueryResponse is what Daraja returns when we query a payment status
type MpesaQueryResponse struct {
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResultCode          string `json:"ResultCode"`
	ResultDesc          string `json:"ResultDesc"`
}

// PaymentStatus constants
const (
	PaymentStatusPending = "pending"
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
)
