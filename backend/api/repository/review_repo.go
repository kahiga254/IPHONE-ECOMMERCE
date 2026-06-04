// api/repository/review_repo.go
package repository

import (
	"database/sql"
	"fmt"

	"backend/api/models"
	"backend/pkg/database"
)

// CreateReview inserts a new review
func CreateReview(review *models.Review) error {
	_, err := database.DB.Exec(`
		INSERT INTO reviews (id, user_id, product_id, rating, comment, is_approved, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		review.ID, review.UserID, review.ProductID, review.Rating, review.Comment, review.IsApproved, review.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	return nil
}

// GetReviewByID fetches a review by ID
func GetReviewByID(reviewID string) (*models.Review, error) {
	var review models.Review
	var comment sql.NullString

	err := database.DB.QueryRow(`
		SELECT id, user_id, product_id, rating, COALESCE(comment, ''), is_approved, created_at
		FROM reviews 
		WHERE id = $1`, reviewID,
	).Scan(
		&review.ID, &review.UserID, &review.ProductID, &review.Rating,
		&comment, &review.IsApproved, &review.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	review.Comment = comment.String

	// Fetch user data
	user, err := GetUserByID(review.UserID)
	if err == nil && user != nil {
		review.User = user
	}

	return &review, nil
}

// GetReviewsByProductID fetches approved reviews for a product
func GetReviewsByProductID(productID string, limit, offset int) ([]models.Review, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.rating, COALESCE(r.comment, ''), r.is_approved, r.created_at
		FROM reviews r
		WHERE r.product_id = $1 AND r.is_approved = true
		ORDER BY r.created_at DESC
		LIMIT $2 OFFSET $3`,
		productID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var comment sql.NullString
		err := rows.Scan(
			&review.ID, &review.UserID, &review.ProductID, &review.Rating,
			&comment, &review.IsApproved, &review.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		review.Comment = comment.String

		// Fetch user data for each review
		user, _ := GetUserByID(review.UserID)
		if user != nil {
			review.User = user
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// CountReviewsByProductID returns total number of approved reviews for a product
func CountReviewsByProductID(productID string) (int, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM reviews 
		WHERE product_id = $1 AND is_approved = true`,
		productID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count reviews: %w", err)
	}
	return count, nil
}

// DeleteReview removes a review
func DeleteReview(reviewID string) error {
	_, err := database.DB.Exec(`DELETE FROM reviews WHERE id = $1`, reviewID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}

// GetPendingReviews fetches unapproved reviews
func GetPendingReviews(limit, offset int) ([]models.Review, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.rating, COALESCE(r.comment, ''), r.is_approved, r.created_at
		FROM reviews r
		WHERE r.is_approved = false
		ORDER BY r.created_at ASC
		LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reviews: %w", err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var comment sql.NullString
		err := rows.Scan(
			&review.ID, &review.UserID, &review.ProductID, &review.Rating,
			&comment, &review.IsApproved, &review.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		review.Comment = comment.String

		// Fetch user data for each review
		user, _ := GetUserByID(review.UserID)
		if user != nil {
			review.User = user
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// CountPendingReviews returns total number of unapproved reviews
func CountPendingReviews() (int, error) {
	var count int
	err := database.DB.QueryRow(`
		SELECT COUNT(*) FROM reviews WHERE is_approved = false`,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count pending reviews: %w", err)
	}
	return count, nil
}

// ApproveReview marks a review as approved
func ApproveReview(reviewID string) error {
	_, err := database.DB.Exec(`
		UPDATE reviews SET is_approved = true 
		WHERE id = $1`,
		reviewID,
	)
	if err != nil {
		return fmt.Errorf("failed to approve review: %w", err)
	}
	return nil
}
