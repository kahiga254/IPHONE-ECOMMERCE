package mpesa

import (
	"fmt"

	"backend/api/models"
)

// ParseCallback validates and parses the raw Daraja callback payload
func ParseCallback(payload models.MpesaCallback) (*CallbackResult, error) {
	stk := payload.Body.StkCallback

	// Validate required fields are present
	if stk.CheckoutRequestID == "" {
		return nil, fmt.Errorf("missing CheckoutRequestID in callback")
	}
	if stk.MerchantRequestID == "" {
		return nil, fmt.Errorf("missing MerchantRequestID in callback")
	}

	result := &CallbackResult{
		CheckoutRequestID: stk.CheckoutRequestID,
		MerchantRequestID: stk.MerchantRequestID,
		ResultCode:        stk.ResultCode,
		ResultDesc:        stk.ResultDesc,
		Success:           stk.ResultCode == 0,
	}

	// Only extract metadata on successful payment
	// On failure Daraja sends an empty CallbackMetadata
	if result.Success {
		result.MpesaReceiptNumber = extractMetadataValue(stk.CallbackMetadata.Item, "MpesaReceiptNumber")
		result.Amount = extractMetadataFloat(stk.CallbackMetadata.Item, "Amount")
		result.PhoneNumber = extractMetadataValue(stk.CallbackMetadata.Item, "PhoneNumber")
		result.TransactionDate = extractMetadataValue(stk.CallbackMetadata.Item, "TransactionDate")

		// Validate that the receipt number is present on success
		if result.MpesaReceiptNumber == "" {
			return nil, fmt.Errorf("missing MpesaReceiptNumber in successful callback")
		}
	}

	return result, nil
}

// CallbackResult holds the parsed and cleaned data from a Daraja callback
type CallbackResult struct {
	CheckoutRequestID  string
	MerchantRequestID  string
	ResultCode         int
	ResultDesc         string
	Success            bool
	MpesaReceiptNumber string
	Amount             float64
	PhoneNumber        string
	TransactionDate    string
}

// ─── Private Helpers ──────────────────────────────────────────────────────────

// extractMetadataValue extracts a string value from Daraja callback metadata by name
func extractMetadataValue(items []struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value"`
}, name string) string {
	for _, item := range items {
		if item.Name == name {
			if v, ok := item.Value.(string); ok {
				return v
			}
			// Some values come as float64 from JSON unmarshalling
			return fmt.Sprintf("%v", item.Value)
		}
	}
	return ""
}

// extractMetadataFloat extracts a float64 value from Daraja callback metadata by name
func extractMetadataFloat(items []struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value"`
}, name string) float64 {
	for _, item := range items {
		if item.Name == name {
			switch v := item.Value.(type) {
			case float64:
				return v
			case string:
				var f float64
				fmt.Sscanf(v, "%f", &f)
				return f
			}
		}
	}
	return 0
}
