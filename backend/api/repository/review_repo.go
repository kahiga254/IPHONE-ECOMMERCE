package repository

import (
	"fmt"

	"backend/api/models"
	"backend/pkg/database"
)

func CreateReview(review *models.Review) error {
	_, err := database.DB.Exec(`
		INSERT INTO reviews (id, user_id, product_id, rating, comment, is_approved)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		review.ID, review.UserID, review.ProductID, review.Rating, review.Comment, review.IsApproved,
	)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	return nil
}

func GetReviewsByProductID(productID string) ([]models.Review, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.rating, COALESCE(r.comment, ''), r.is_approved, r.created_at,
		       COALESCE(u.name, '') as user_name
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.product_id = $1 AND r.is_approved = true
		ORDER BY r.created_at DESC`,
		productID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var userName string
		err := rows.Scan(
			&review.ID, &review.UserID, &review.ProductID, &review.Rating,
			&review.Comment, &review.IsApproved, &review.CreatedAt,
			&userName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		review.UserName = userName
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func GetPendingReviews() ([]models.Review, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.rating, COALESCE(r.comment, ''), r.is_approved, r.created_at,
		       COALESCE(u.name, '') as user_name, COALESCE(p.name, '') as product_name
		FROM reviews r
		LEFT JOIN users u ON r.user_id = u.id
		LEFT JOIN products p ON r.product_id = p.id
		WHERE r.is_approved = false
		ORDER BY r.created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending reviews: %w", err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		var userName, productName string
		err := rows.Scan(
			&review.ID, &review.UserID, &review.ProductID, &review.Rating,
			&review.Comment, &review.IsApproved, &review.CreatedAt,
			&userName, &productName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		review.UserName = userName
		review.ProductName = productName
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func ApproveReview(reviewID string) error {
	_, err := database.DB.Exec(`
		UPDATE reviews SET is_approved = true, updated_at = NOW()
		WHERE id = $1`,
		reviewID,
	)
	if err != nil {
		return fmt.Errorf("failed to approve review: %w", err)
	}
	return nil
}

func DeleteReview(reviewID string) error {
	_, err := database.DB.Exec(`DELETE FROM reviews WHERE id = $1`, reviewID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}
