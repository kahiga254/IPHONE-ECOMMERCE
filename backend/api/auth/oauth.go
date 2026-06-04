package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"backend/config"
)

// GoogleUser holds the user info we get back from Google
type GoogleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"picture"`
	Verified  bool   `json:"verified_email"`
}

// GetGoogleAuthURL returns the Google OAuth2 consent page URL
func GetGoogleAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", config.App.GoogleClientID)
	params.Set("redirect_uri", config.App.GoogleRedirectURL)
	params.Set("response_type", "code")
	params.Set("scope", "openid email profile")
	params.Set("state", state)
	params.Set("access_type", "offline")

	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

// ExchangeGoogleCode exchanges the authorization code for an access token
func ExchangeGoogleCode(code string) (string, error) {
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"code":          {code},
		"client_id":     {config.App.GoogleClientID},
		"client_secret": {config.App.GoogleClientSecret},
		"redirect_uri":  {config.App.GoogleRedirectURL},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	accessToken, ok := result["access_token"].(string)
	if !ok || accessToken == "" {
		return "", fmt.Errorf("no access token in response")
	}

	return accessToken, nil
}

// GetGoogleUser uses the access token to fetch the user's profile from Google
func GetGoogleUser(accessToken string) (*GoogleUser, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch google user: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var googleUser GoogleUser
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse google user: %w", err)
	}

	if googleUser.Email == "" {
		return nil, fmt.Errorf("no email returned from google")
	}

	return &googleUser, nil
}
