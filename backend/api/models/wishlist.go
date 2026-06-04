package models

import "time"

type Wishlist struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	VariantID string    `json:"variant_id"`
	Variant   *Variant  `json:"variant,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// AddToWishlistRequest is used when a user adds a variant to their wishlist
type AddToWishlistRequest struct {
	VariantID string `json:"variant_id" binding:"required"`
}
