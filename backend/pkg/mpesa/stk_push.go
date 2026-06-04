package mpesa

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"backend/api/models"
	"backend/config"
)

// stkPushRequest is the payload we send to Daraja to trigger an STK push
type stkPushRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            int    `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

// InitiateSTKPush sends an STK push request to Daraja and returns the response
func InitiateSTKPush(phone string, amount float64, orderNumber string) (*models.StkPushResponse, error) {
	// Get a fresh Daraja access token
	accessToken, err := GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	timestamp := GenerateTimestamp()
	password := GeneratePassword(timestamp)

	// Build the STK push payload
	payload := stkPushRequest{
		BusinessShortCode: config.App.MPesaShortcode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            int(amount), // Daraja requires whole numbers
		PartyA:            phone,
		PartyB:            config.App.MPesaShortcode,
		PhoneNumber:       phone,
		CallBackURL:       config.App.MPesaCallbackURL,
		AccountReference:  orderNumber,
		TransactionDesc:   "Payment for order " + orderNumber,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal STK push payload: %w", err)
	}

	// Build the request
	url := baseURL() + "/mpesa/stkpush/v1/processrequest"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create STK push request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send the request with a 15 second timeout
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send STK push request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read STK push response: %w", err)
	}

	var stkResp models.StkPushResponse
	if err := json.Unmarshal(body, &stkResp); err != nil {
		return nil, fmt.Errorf("failed to parse STK push response: %w", err)
	}

	// ResponseCode 0 means the STK push was sent successfully
	if stkResp.ResponseCode != "0" {
		return nil, fmt.Errorf("STK push rejected by Daraja: %s", stkResp.ResponseDescription)
	}

	return &stkResp, nil
}

// QuerySTKStatus queries Daraja for the status of an STK push transaction
func QuerySTKStatus(checkoutRequestID string) (*models.MpesaQueryResponse, error) {
	// Get a fresh Daraja access token
	accessToken, err := GetAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	timestamp := GenerateTimestamp()
	password := GeneratePassword(timestamp)

	// Build the query payload
	payload := map[string]string{
		"BusinessShortCode": config.App.MPesaShortcode,
		"Password":          password,
		"Timestamp":         timestamp,
		"CheckoutRequestID": checkoutRequestID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query payload: %w", err)
	}

	// Build the request
	url := baseURL() + "/mpesa/stkpushquery/v1/query"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create query request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send query request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read query response: %w", err)
	}

	var queryResp models.MpesaQueryResponse
	if err := json.Unmarshal(body, &queryResp); err != nil {
		return nil, fmt.Errorf("failed to parse query response: %w", err)
	}

	return &queryResp, nil
}
