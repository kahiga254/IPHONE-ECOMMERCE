// api/services/review_service.go
package services

import (
	"fmt"
	"time"

	"backend/api/models"
	"backend/api/repository"

	"github.com/google/uuid"
)

// CreateReview creates a new product review
func CreateReview(userID string, req models.CreateReviewRequest) (*models.Review, error) {
	// Create review
	review := &models.Review{
		ID:         uuid.New().String(),
		UserID:     userID,
		ProductID:  req.ProductID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsApproved: false, // Reviews need admin approval by default
		CreatedAt:  time.Now(),
	}

	err := repository.CreateReview(review)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Get user details for response
	user, _ := repository.GetUserByID(userID)
	if user != nil {
		review.User = user
	}

	return review, nil
}

// GetReviewsByProductID returns all approved reviews for a product
func GetReviewsByProductID(productID string, limit, offset int) ([]models.Review, int, error) {
	reviews, err := repository.GetReviewsByProductID(productID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get reviews: %w", err)
	}

	total, err := repository.CountReviewsByProductID(productID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count reviews: %w", err)
	}

	return reviews, total, nil
}

// DeleteReview deletes a review (user can delete their own, admin can delete any)
func DeleteReview(reviewID, userID, role string) error {
	// Get review to check ownership
	review, err := repository.GetReviewByID(reviewID)
	if err != nil {
		return fmt.Errorf("failed to get review: %w", err)
	}
	if review == nil {
		return fmt.Errorf("review not found")
	}

	// Check if user is authorized (owner or admin)
	if review.UserID != userID && role != "admin" {
		return fmt.Errorf("unauthorized to delete this review")
	}

	err = repository.DeleteReview(reviewID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

// GetPendingReviews returns all unapproved reviews (admin only)
func GetPendingReviews(limit, offset int) ([]models.Review, int, error) {
	reviews, err := repository.GetPendingReviews(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pending reviews: %w", err)
	}

	total, err := repository.CountPendingReviews()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pending reviews: %w", err)
	}

	return reviews, total, nil
}

// ApproveReview approves a review (admin only)
func ApproveReview(reviewID string) error {
	// Check if review exists
	review, err := repository.GetReviewByID(reviewID)
	if err != nil {
		return fmt.Errorf("failed to get review: %w", err)
	}
	if review == nil {
		return fmt.Errorf("review not found")
	}

	err = repository.ApproveReview(reviewID)
	if err != nil {
		return fmt.Errorf("failed to approve review: %w", err)
	}

	return nil
}
