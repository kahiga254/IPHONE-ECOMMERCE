package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		return "", err
	}

	req.SetBasicAuth(m.consumerKey, m.consumerSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (m *MpesaService) generatePassword() string {
	timestamp := time.Now().Format("20060102150405")
	// Password = Base64(Shortcode + Passkey + Timestamp)
	return timestamp // Will be properly encoded in the request
}

func (m *MpesaService) InitiateSTKPush(phoneNumber string, amount float64, accountRef string) (*STKPushResponse, error) {
	// Get access token
	token, err := m.getAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	timestamp := time.Now().Format("20060102150405")
	password := m.shortcode + m.passkey + timestamp

	// The actual password needs to be base64 encoded
	// For now, we'll use a placeholder - this will be implemented properly

	request := STKPushRequest{
		BusinessShortCode: m.shortcode,
		Password:          password,
		Timestamp:         timestamp,
		TransactionType:   "CustomerPayBillOnline",
		Amount:            fmt.Sprintf("%.0f", amount),
		PartyA:            phoneNumber,
		PartyB:            m.shortcode,
		PhoneNumber:       phoneNumber,
		CallBackURL:       config.App.MPesaCallbackURL,
		AccountReference:  accountRef,
		TransactionDesc:   "Payment for order " + accountRef,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	url := m.getBaseURL() + "/mpesa/stkpush/v1/processrequest"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result STKPushResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.ResponseCode != "0" {
		return nil, fmt.Errorf("STK push failed: %s", result.ResponseDescription)
	}

	return &result, nil
}
