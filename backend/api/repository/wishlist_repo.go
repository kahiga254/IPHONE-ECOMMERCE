package repository

import (
	"fmt"

	"backend/api/models"
	"backend/pkg/database"
	"encoding/json"
)

// GetWishlistByUserID fetches all wishlist items for a given user
func GetWishlistByUserID(userID string) ([]models.Wishlist, error) {
	rows, err := database.DB.Query(`
		SELECT w.id, w.user_id, w.variant_id, w.created_at,
		       pv.sku, pv.color, pv.storage, pv.price, pv.stock, pv.images,
		       p.name, p.slug
		FROM wishlists w
		JOIN product_variants pv ON w.variant_id = pv.id
		JOIN products p          ON pv.product_id = p.id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wishlist: %w", err)
	}
	defer rows.Close()

	wishlist := []models.Wishlist{}
	for rows.Next() {
		var w models.Wishlist
		var v models.Variant
		var p models.Product
		var imagesJSON []byte

		err := rows.Scan(
			&w.ID, &w.UserID, &w.VariantID, &w.CreatedAt,
			&v.SKU, &v.Color, &v.Storage, &v.Price, &v.Stock, &imagesJSON,
			&p.Name, &p.Slug,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan wishlist item: %w", err)
		}

		json.Unmarshal(imagesJSON, &v.Images)
		v.ID = w.VariantID
		w.Variant = &v
		wishlist = append(wishlist, w)
	}

	return wishlist, nil
}

// AddToWishlist inserts a new wishlist item for a user
func AddToWishlist(userID, variantID string) (*models.Wishlist, error) {
	// Check if item already exists in wishlist
	var count int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM wishlists
		WHERE user_id = $1 AND variant_id = $2`,
		userID, variantID,
	).Scan(&count)

	if count > 0 {
		return nil, fmt.Errorf("item already in wishlist")
	}

	var w models.Wishlist
	err := database.DB.QueryRow(`
		INSERT INTO wishlists (user_id, variant_id)
		VALUES ($1, $2)
		RETURNING id, user_id, variant_id, created_at`,
		userID, variantID,
	).Scan(&w.ID, &w.UserID, &w.VariantID, &w.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add to wishlist: %w", err)
	}

	return &w, nil
}

// RemoveFromWishlist deletes a wishlist item by variant ID and user ID
func RemoveFromWishlist(userID, variantID string) error {
	result, err := database.DB.Exec(`
		DELETE FROM wishlists
		WHERE user_id = $1 AND variant_id = $2`,
		userID, variantID,
	)
	if err != nil {
		return fmt.Errorf("failed to remove from wishlist: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found in wishlist")
	}

	return nil
}

// ClearWishlist deletes all wishlist items for a user
func ClearWishlist(userID string) error {
	_, err := database.DB.Exec(`
		DELETE FROM wishlists WHERE user_id = $1`, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to clear wishlist: %w", err)
	}
	return nil
}
