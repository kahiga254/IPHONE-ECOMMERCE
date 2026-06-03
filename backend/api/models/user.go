package models

import "time"

type User struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         *string   `json:"email,omitempty"`
	Phone         *string   `json:"phone,omitempty"`
	PasswordHash  *string   `json:"-"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	Provider      string    `json:"provider"`
	ProviderID    *string   `json:"-"`
	IsVerified    bool      `json:"is_verified"`
	IsActive      bool      `json:"is_active"`
	Role          string    `json:"role"`
	LoyaltyPoints int       `json:"loyalty_points"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// RegisterRequest is used when a new user signs up with email and password
type RegisterRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Phone    string `json:"phone"    binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest is used for email/password login
type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// PhoneLoginRequest is used to request an OTP for phone login
type PhoneLoginRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// VerifyOTPRequest is used to verify the OTP sent to a phone
type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code"  binding:"required,len=6"`
}

// RefreshTokenRequest is used to get a new access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest is used to trigger a password reset email
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest is used to set a new password using a reset token
type ResetPasswordRequest struct {
	Token    string `json:"token"    binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

// UpdateProfileRequest is used to update a user's name, phone or avatar
type UpdateProfileRequest struct {
	Name      string `json:"name"       binding:"omitempty,min=2,max=100"`
	Phone     string `json:"phone"      binding:"omitempty"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

// AuthResponse is returned after a successful login or register
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
