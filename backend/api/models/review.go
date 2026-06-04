package models

import "time"

type Review struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ProductID  string    `json:"product_id"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment,omitempty"`
	IsApproved bool      `json:"is_approved"`
	User       *User     `json:"user,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// CreateReviewRequest is used when a user submits a review for a product
type CreateReviewRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Rating    int    `json:"rating"     binding:"required,min=1,max=5"`
	Comment   string `json:"comment"    binding:"omitempty,min=10,max=1000"`
}

// UpdateReviewRequest is used by admin to approve or reject a review
type UpdateReviewRequest struct {
	IsApproved bool `json:"is_approved" binding:"required"`
}
