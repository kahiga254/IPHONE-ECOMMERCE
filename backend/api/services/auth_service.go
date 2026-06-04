package services

import (
	"fmt"

	"backend/api/auth"
	"backend/api/models"
	"backend/api/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user account with email and password
func Register(req models.RegisterRequest) (*models.AuthResponse, error) {
	// Check if email is already taken
	existing, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Check if phone is already taken
	existingPhone, err := repository.GetUserByPhone(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check phone: %w", err)
	}
	if existingPhone != nil {
		return nil, fmt.Errorf("phone number already registered")
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Build user model
	user := &models.User{
		Name:     req.Name,
		Email:    &req.Email,
		Phone:    &req.Phone,
		Provider: "local",
	}

	// Save user to database
	created, err := repository.CreateUser(user, string(hash))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate email verification token and save it
	token := uuid.NewString()
	if err := repository.SaveVerificationToken(created.ID, token, "email_verify"); err != nil {
		return nil, fmt.Errorf("failed to save verification token: %w", err)
	}

	// TODO: Send verification email
	// email.SendVerificationEmail(*created.Email, token)

	// Generate JWT tokens
	accessToken, err := auth.GenerateAccessToken(created.ID, created.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(created.ID, created.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token
	if err := repository.SaveRefreshToken(created.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         created,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Login authenticates a user with email and password
func Login(req models.LoginRequest) (*models.AuthResponse, error) {
	// Fetch user by email
	user, err := repository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check account is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Check user registered with email not OAuth
	if user.Provider != "local" {
		return nil, fmt.Errorf("please login with %s", user.Provider)
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Generate JWT tokens
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token
	if err := repository.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	// Clear password hash before returning
	user.PasswordHash = nil

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// LoginWithPhone sends an OTP to the given phone number
func LoginWithPhone(phone string) error {
	// Generate OTP
	code, err := auth.GenerateOTP()
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	// Save OTP to database
	if err := auth.SaveOTP(phone, code); err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	// TODO: Send OTP via Africa's Talking
	// sms.SendOTP(phone, code)
	fmt.Printf("OTP for %s: %s\n", phone, code) // dev only — remove in production

	return nil
}

// VerifyPhoneOTP verifies the OTP and logs the user in or creates their account
func VerifyPhoneOTP(req models.VerifyOTPRequest) (*models.AuthResponse, error) {
	// Verify OTP
	valid, err := auth.VerifyOTP(req.Phone, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to verify OTP: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid or expired OTP")
	}

	// Find existing user or create a new one
	user, err := repository.GetUserByPhone(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if user == nil {
		// Auto create account for new phone users
		newUser := &models.User{
			Name:       "User " + req.Phone[len(req.Phone)-4:],
			Phone:      &req.Phone,
			Provider:   "phone",
			IsVerified: true,
		}
		user, err = repository.CreateUser(newUser, "")
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Generate JWT tokens
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token
	if err := repository.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GoogleLogin handles the OAuth callback and logs the user in
func GoogleLogin(code string) (*models.AuthResponse, error) {
	// Exchange code for Google access token
	googleAccessToken, err := auth.ExchangeGoogleCode(code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange google code: %w", err)
	}

	// Fetch user info from Google
	googleUser, err := auth.GetGoogleUser(googleAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get google user: %w", err)
	}

	// Upsert user in our database
	user, err := repository.UpsertOAuthUser(
		googleUser.Name,
		googleUser.Email,
		googleUser.AvatarURL,
		"google",
		googleUser.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert google user: %w", err)
	}

	// Generate JWT tokens
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save refresh token
	if err := repository.SaveRefreshToken(user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshToken rotates the refresh token and returns a new access token
func RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	// Validate the refresh token signature and expiry
	claims, err := auth.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	// Check token exists in database
	userID, err := repository.GetRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify refresh token: %w", err)
	}
	if userID == "" {
		return nil, fmt.Errorf("refresh token expired or already used")
	}

	// Delete old refresh token — rotation prevents reuse
	if err := repository.DeleteRefreshToken(refreshToken); err != nil {
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	// Generate new token pair
	newAccessToken, err := auth.GenerateAccessToken(claims.UserID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := auth.GenerateRefreshToken(claims.UserID, claims.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Save new refresh token
	if err := repository.SaveRefreshToken(claims.UserID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	return &models.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// Logout deletes the refresh token from the database
func Logout(refreshToken string) error {
	if err := repository.DeleteRefreshToken(refreshToken); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

// VerifyEmail marks a user's email as verified using the token
func VerifyEmail(token string) error {
	userID, err := repository.GetVerificationToken(token, "email_verify")
	if err != nil {
		return fmt.Errorf("failed to get verification token: %w", err)
	}
	if userID == "" {
		return fmt.Errorf("invalid or expired verification token")
	}

	if err := repository.VerifyUserEmail(userID); err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}

	if err := repository.MarkTokenUsed(token); err != nil {
		return fmt.Errorf("failed to mark token used: %w", err)
	}

	return nil
}
