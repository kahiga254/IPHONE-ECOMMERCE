package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"backend/pkg/database"
)

// GenerateOTP creates a cryptographically secure 6 digit OTP
func GenerateOTP() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

// SaveOTP deletes any existing OTPs for the phone number then saves the new one
func SaveOTP(phone, code string) error {
	// Delete any previous unused OTPs for this phone
	_, err := database.DB.Exec(`DELETE FROM otps WHERE phone = $1`, phone)
	if err != nil {
		return fmt.Errorf("failed to clear old OTPs: %w", err)
	}

	// Save the new OTP with a 5 minute expiry
	_, err = database.DB.Exec(`
		INSERT INTO otps (phone, code, expires_at)
		VALUES ($1, $2, $3)`,
		phone, code, time.Now().Add(5*time.Minute),
	)
	if err != nil {
		return fmt.Errorf("failed to save OTP: %w", err)
	}

	return nil
}

// VerifyOTP checks if the OTP is valid, not expired and not already used
func VerifyOTP(phone, code string) (bool, error) {
	var id string
	var expiresAt time.Time

	err := database.DB.QueryRow(`
		SELECT id, expires_at FROM otps
		WHERE phone = $1 AND code = $2 AND verified = FALSE`,
		phone, code,
	).Scan(&id, &expiresAt)

	if err != nil {
		return false, nil
	}

	// Check if OTP has expired
	if time.Now().After(expiresAt) {
		return false, nil
	}

	// Mark OTP as verified so it cannot be reused
	_, err = database.DB.Exec(`UPDATE otps SET verified = TRUE WHERE id = $1`, id)
	if err != nil {
		return false, fmt.Errorf("failed to mark OTP as verified: %w", err)
	}

	return true, nil
}
