package services

import (
	"fmt"

	"backend/api/models"
	"backend/api/repository"
	"github.com/google/uuid"
)

func CreateReview(userID string, req models.CreateReviewRequest) (*models.Review, error) {
	review := &models.Review{
		ID:         uuid.New().String(),
		UserID:     userID,
		ProductID:  req.ProductID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsApproved: true, // Auto-approve reviews
	}

	err := repository.CreateReview(review)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

func GetProductReviews(productID string) ([]models.Review, error) {
	reviews, err := repository.GetReviewsByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	return reviews, nil
}

func GetPendingReviews() ([]models.Review, error) {
	reviews, err := repository.GetPendingReviews()
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reviews: %w", err)
	}
	return reviews, nil
}

func ApproveReview(reviewID string) error {
	err := repository.ApproveReview(reviewID)
	if err != nil {
		return fmt.Errorf("failed to approve review: %w", err)
	}
	return nil
}

func DeleteReview(reviewID string) error {
	err := repository.DeleteReview(reviewID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}
