package mpesa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"backend/config"
)

// baseURL returns the correct Daraja base URL based on the environment
func baseURL() string {
	if config.App.MpesaEnv == "production" {
		return "https://api.safaricom.co.ke"
	}
	return "https://sandbox.safaricom.co.ke"
}

// accessTokenResponse holds the response from the Daraja auth endpoint
type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// GetAccessToken fetches a fresh OAuth access token from Daraja
func GetAccessToken() (string, error) {
	url := baseURL() + "/oauth/v1/generate?grant_type=client_credentials"

	// Base64 encode consumer key and secret
	credentials := base64.StdEncoding.EncodeToString(
		[]byte(config.App.MpesaConsumerKey + ":" + config.App.MpesaConsumerSecret),
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Authorization", "Basic "+credentials)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call daraja auth: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read auth response: %w", err)
	}

	var tokenResp accessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("empty access token received from daraja")
	}

	return tokenResp.AccessToken, nil
}

// GeneratePassword generates the Daraja STK push password
// It is a base64 encoding of ShortCode + Passkey + Timestamp
func GeneratePassword(timestamp string) string {
	raw := config.App.MpesaShortCode + config.App.MpesaPasskey + timestamp
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

// GenerateTimestamp returns the current time in Daraja's required format YYYYMMDDHHmmss
func GenerateTimestamp() string {
	return time.Now().Format("20060102150405")
}

// NormalizePhone converts a Kenyan phone number to the format Daraja expects: 2547XXXXXXXX
func NormalizePhone(phone string) (string, error) {
	// Remove all spaces
	phone = strings.ReplaceAll(phone, " ", "")

	// Handle +2547XXXXXXXX
	if strings.HasPrefix(phone, "+254") {
		phone = strings.TrimPrefix(phone, "+")
		return validatePhone(phone)
	}

	// Handle 2547XXXXXXXX
	if strings.HasPrefix(phone, "254") {
		return validatePhone(phone)
	}

	// Handle 07XXXXXXXX
	if strings.HasPrefix(phone, "0") {
		phone = "254" + phone[1:]
		return validatePhone(phone)
	}

	// Handle 7XXXXXXXX
	if strings.HasPrefix(phone, "7") || strings.HasPrefix(phone, "1") {
		phone = "254" + phone
		return validatePhone(phone)
	}

	return "", fmt.Errorf("unrecognized phone number format: %s", phone)
}

// validatePhone checks the final phone number is exactly 12 digits
func validatePhone(phone string) (string, error) {
	if len(phone) != 12 {
		return "", fmt.Errorf("invalid phone number length: %s", phone)
	}
	return phone, nil
}
