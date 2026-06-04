package repository

import (
	"database/sql"
	"fmt"

	"backend/api/models"
	"backend/pkg/database"
)

// GetAddressByID fetches a single address that belongs to a user
func GetAddressByID(addressID, userID string) (*models.Address, error) {
	var a models.Address

	err := database.DB.QueryRow(`
		SELECT id, user_id, label, full_name, phone, county, town, street, is_default, created_at
		FROM addresses WHERE id = $1 AND user_id = $2`,
		addressID, userID,
	).Scan(
		&a.ID, &a.UserID, &a.Label, &a.FullName, &a.Phone,
		&a.County, &a.Town, &a.Street, &a.IsDefault, &a.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}

	return &a, nil
}

// GetAddressesByUserID fetches all addresses for a user
func GetAddressesByUserID(userID string) ([]models.Address, error) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, label, full_name, phone, county, town, street, is_default, created_at
		FROM addresses WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch addresses: %w", err)
	}
	defer rows.Close()

	addresses := []models.Address{}
	for rows.Next() {
		var a models.Address
		rows.Scan(
			&a.ID, &a.UserID, &a.Label, &a.FullName, &a.Phone,
			&a.County, &a.Town, &a.Street, &a.IsDefault, &a.CreatedAt,
		)
		addresses = append(addresses, a)
	}

	return addresses, nil
}

// CreateAddress inserts a new address for a user
func CreateAddress(userID string, req models.CreateAddressRequest) (*models.Address, error) {
	// If this is set as default unset all other defaults first
	if req.IsDefault {
		database.DB.Exec(`
			UPDATE addresses SET is_default = FALSE WHERE user_id = $1`, userID,
		)
	}

	var a models.Address
	err := database.DB.QueryRow(`
		INSERT INTO addresses (user_id, label, full_name, phone, county, town, street, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, label, full_name, phone, county, town, street, is_default, created_at`,
		userID, req.Label, req.FullName, req.Phone,
		req.County, req.Town, req.Street, req.IsDefault,
	).Scan(
		&a.ID, &a.UserID, &a.Label, &a.FullName, &a.Phone,
		&a.County, &a.Town, &a.Street, &a.IsDefault, &a.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create address: %w", err)
	}

	return &a, nil
}

// UpdateAddress updates an existing address
func UpdateAddress(addressID, userID string, req models.UpdateAddressRequest) error {
	if req.IsDefault {
		database.DB.Exec(`
			UPDATE addresses SET is_default = FALSE WHERE user_id = $1`, userID,
		)
	}

	_, err := database.DB.Exec(`
		UPDATE addresses SET
			label      = $1,
			full_name  = $2,
			phone      = $3,
			county     = $4,
			town       = $5,
			street     = $6,
			is_default = $7
		WHERE id = $8 AND user_id = $9`,
		req.Label, req.FullName, req.Phone,
		req.County, req.Town, req.Street,
		req.IsDefault, addressID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to update address: %w", err)
	}

	return nil
}

// DeleteAddress removes an address belonging to a user
func DeleteAddress(addressID, userID string) error {
	result, err := database.DB.Exec(`
		DELETE FROM addresses WHERE id = $1 AND user_id = $2`,
		addressID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete address: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("address not found or does not belong to you")
	}

	return nil
}
