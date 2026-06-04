package services

import (
	"fmt"

	"backend/api/models"
	"backend/api/repository"
)

// GetWishlist fetches all wishlist items for a user
func GetWishlist(userID string) ([]models.Wishlist, error) {
	wishlist, err := repository.GetWishlistByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wishlist: %w", err)
	}
	return wishlist, nil
}

// AddToWishlist adds a variant to the user's wishlist
func AddToWishlist(userID string, req models.AddToWishlistRequest) (*models.Wishlist, error) {
	// Check the variant actually exists
	variant, err := repository.GetVariantByID(req.VariantID)
	if err != nil {
		return nil, fmt.Errorf("failed to check variant: %w", err)
	}
	if variant == nil {
		return nil, fmt.Errorf("variant not found")
	}

	wishlist, err := repository.AddToWishlist(userID, req.VariantID)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

// RemoveFromWishlist removes a variant from the user's wishlist
func RemoveFromWishlist(userID, variantID string) error {
	if err := repository.RemoveFromWishlist(userID, variantID); err != nil {
		return err
	}
	return nil
}

// ClearWishlist removes all items from the user's wishlist
func ClearWishlist(userID string) error {
	if err := repository.ClearWishlist(userID); err != nil {
		return fmt.Errorf("failed to clear wishlist: %w", err)
	}
	return nil
}
