package repository

import (
	"database/sql"
	"fmt"
	"time"

	"backend/api/models"
	"backend/pkg/database"
)

// CreateUser inserts a new user into the database and returns the created user
func CreateUser(user *models.User, passwordHash string) (*models.User, error) {
	var created models.User

	err := database.DB.QueryRow(`
		INSERT INTO users (name, email, phone, password_hash, provider, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, email, phone, avatar_url, provider, is_verified, is_active, role, loyalty_points, created_at, updated_at`,
		user.Name, user.Email, user.Phone, passwordHash, user.Provider, user.IsVerified,
	).Scan(
		&created.ID, &created.Name, &created.Email, &created.Phone,
		&created.AvatarURL, &created.Provider, &created.IsVerified,
		&created.IsActive, &created.Role, &created.LoyaltyPoints,
		&created.CreatedAt, &created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &created, nil
}

// GetUserByID fetches a single user by their UUID
func GetUserByID(id string) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(`
		SELECT id, name, email, phone, password_hash, avatar_url, provider,
		       is_verified, is_active, role, loyalty_points, created_at, updated_at
		FROM users WHERE id = $1`, id,
	).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash,
		&user.AvatarURL, &user.Provider, &user.IsVerified, &user.IsActive,
		&user.Role, &user.LoyaltyPoints, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetUserByEmail fetches a single user by their email address
func GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(`
		SELECT id, name, email, phone, password_hash, avatar_url, provider,
		       is_verified, is_active, role, loyalty_points, created_at, updated_at
		FROM users WHERE email = $1`, email,
	).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash,
		&user.AvatarURL, &user.Provider, &user.IsVerified, &user.IsActive,
		&user.Role, &user.LoyaltyPoints, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetUserByPhone fetches a single user by their phone number
func GetUserByPhone(phone string) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(`
		SELECT id, name, email, phone, password_hash, avatar_url, provider,
		       is_verified, is_active, role, loyalty_points, created_at, updated_at
		FROM users WHERE phone = $1`, phone,
	).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash,
		&user.AvatarURL, &user.Provider, &user.IsVerified, &user.IsActive,
		&user.Role, &user.LoyaltyPoints, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by phone: %w", err)
	}

	return &user, nil
}

// UpsertOAuthUser creates a new OAuth user or updates their info if they already exist
func UpsertOAuthUser(name, email, avatarURL, provider, providerID string) (*models.User, error) {
	var user models.User

	err := database.DB.QueryRow(`
		INSERT INTO users (name, email, avatar_url, provider, provider_id, is_verified)
		VALUES ($1, $2, $3, $4, $5, TRUE)
		ON CONFLICT (email) DO UPDATE SET
			name       = EXCLUDED.name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()
		RETURNING id, name, email, phone, avatar_url, provider, is_verified, is_active, role, loyalty_points, created_at, updated_at`,
		name, email, avatarURL, provider, providerID,
	).Scan(
		&user.ID, &user.Name, &user.Email, &user.Phone,
		&user.AvatarURL, &user.Provider, &user.IsVerified,
		&user.IsActive, &user.Role, &user.LoyaltyPoints,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert oauth user: %w", err)
	}

	return &user, nil
}

// UpdateProfile updates a user's name, phone and avatar
func UpdateProfile(id, name, phone, avatarURL string) error {
	_, err := database.DB.Exec(`
		UPDATE users SET name = $1, phone = $2, avatar_url = $3, updated_at = NOW()
		WHERE id = $4`,
		name, phone, avatarURL, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}

// VerifyUserEmail marks a user's email as verified
func VerifyUserEmail(userID string) error {
	_, err := database.DB.Exec(`UPDATE users SET is_verified = TRUE WHERE id = $1`, userID)
	if err != nil {
		return fmt.Errorf("failed to verify email: %w", err)
	}
	return nil
}

// UpdatePassword updates a user's hashed password
func UpdatePassword(userID, passwordHash string) error {
	_, err := database.DB.Exec(`
		UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`,
		passwordHash, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}
	return nil
}

// SaveVerificationToken stores an email verification or password reset token
func SaveVerificationToken(userID, token, tokenType string) error {
	_, err := database.DB.Exec(`
		INSERT INTO verification_tokens (user_id, token, type, expires_at)
		VALUES ($1, $2, $3, $4)`,
		userID, token, tokenType, time.Now().Add(24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("failed to save verification token: %w", err)
	}
	return nil
}

// GetVerificationToken fetches a valid unused token
func GetVerificationToken(token, tokenType string) (string, error) {
	var userID string
	var expiresAt time.Time

	err := database.DB.QueryRow(`
		SELECT user_id, expires_at FROM verification_tokens
		WHERE token = $1 AND type = $2 AND used = FALSE`,
		token, tokenType,
	).Scan(&userID, &expiresAt)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get verification token: %w", err)
	}
	if time.Now().After(expiresAt) {
		return "", nil
	}

	return userID, nil
}

// MarkTokenUsed marks a verification token as used so it cannot be reused
func MarkTokenUsed(token string) error {
	_, err := database.DB.Exec(`UPDATE verification_tokens SET used = TRUE WHERE token = $1`, token)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}
	return nil
}

// SaveRefreshToken stores a refresh token in the database
func SaveRefreshToken(userID, token string) error {
	_, err := database.DB.Exec(`
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`,
		userID, token, time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

// GetRefreshToken fetches a valid refresh token from the database
func GetRefreshToken(token string) (string, error) {
	var userID string

	err := database.DB.QueryRow(`
		SELECT user_id FROM refresh_tokens
		WHERE token = $1 AND expires_at > NOW()`, token,
	).Scan(&userID)

	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	return userID, nil
}

// DeleteRefreshToken removes a refresh token on logout
func DeleteRefreshToken(token string) error {
	_, err := database.DB.Exec(`DELETE FROM refresh_tokens WHERE token = $1`, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}
