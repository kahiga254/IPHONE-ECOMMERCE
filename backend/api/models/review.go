package models

import "time"

type Review struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProductID   string    `json:"product_id"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	IsApproved  bool      `json:"is_approved"`
	UserName    string    `json:"user_name,omitempty"`
	ProductName string    `json:"product_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateReviewRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required,min=10"`
}
