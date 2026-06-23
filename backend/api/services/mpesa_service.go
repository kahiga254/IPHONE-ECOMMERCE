package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/config"
)

type MpesaService struct {
	consumerKey    string
	consumerSecret string
	passkey        string
	shortcode      string
	environment    string
}

type STKPushRequest struct {
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            string `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

func NewMpesaService() *MpesaService {
	return &MpesaService{
		consumerKey:    config.App.MPesaConsumerKey,
		consumerSecret: config.App.MPesaConsumerSecret,
		passkey:        config.App.MPesaPasskey,
		shortcode:      config.App.MPesaShortcode,
		environment:    config.App.MPesaEnvironment,
	}
}

func (m *MpesaService) getBaseURL() string {
	if m.environment == "production" {
		return "https://api.safaricom.co.ke"
	}
	return "https://sandbox.safaricom.co.ke"
}

func (m *MpesaService) getAccessToken() (string, error) {
	url := m.getBaseURL() + "/oauth/v1/generate?grant_type=client_credentials"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(m.consumerKey, m.consumerSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.AccessToken == "" {
		return "", fmt.Errorf("empty access token received")
	}

	return result.AccessToken, nil
}

func (m *MpesaService) generatePassword(timestamp string) string {
	// Password = Base64(Shortcode + Passkey + Timestamp)
	data := m.shortcode + m.passkey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func normalizePhone(phone string) string {
	// Remove spaces
	phone = strings.ReplaceAll(phone, " ", "")

	// Remove + if present

	phone = strings.TrimPrefix(phone, "+")

	// If starts with 0, replace with 254
	if strings.HasPrefix(phone, "0") {
		phone = "254" + phone[1:]
	}

	// If starts with 7, add 254
	if strings.HasPrefix(phone, "7") {
		phone = "254" + phone
	}

	return phone
}

func (m *MpesaService) InitiateSTKPush(phoneNumber string, amount float64, accountRef string) (*STKPushResponse, error) {
	// Get access token
	token, err := m.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// Normalize phone number
	phone := normalizePhone(phoneNumber)
	if len(phone) != 12 {
		return nil, fmt.Errorf("invalid phone number format: %s (must be 12 digits)", phoneNumber)
	}

	timestamp := time.Now().Format("20060102150405")
	password := m.generatePassword(timestamp)

	request := STKPushRequest{
		BusinessShortCode: m.shortcode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            fmt.Sprintf("%.0f", amount),
		PartyA:            phone,
		PartyB:            m.shortcode,
		PhoneNumber:       phone,
		CallBackURL:       config.App.MPesaCallbackURL,
		AccountReference:  accountRef,
		TransactionDesc:   "Payment for order " + accountRef,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := m.getBaseURL() + "/mpesa/stkpush/v1/processrequest"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result STKPushResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.ResponseCode != "0" {
		return nil, fmt.Errorf("STK push failed: %s", result.ResponseDescription)
	}

	return &result, nil
}
