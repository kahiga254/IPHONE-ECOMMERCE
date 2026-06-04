package repository

import (
	"database/sql"
	"fmt"

	"backend/api/models"
	"backend/pkg/database"
)

// CreateReview inserts a new review into the database
func CreateReview(userID string, req models.CreateReviewRequest) (*models.Review, error) {
	// Check if user has already reviewed this product
	var count int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM reviews
		WHERE user_id = $1 AND product_id = $2`,
		userID, req.ProductID,
	).Scan(&count)

	if count > 0 {
		return nil, fmt.Errorf("you have already reviewed this product")
	}

	// Check if user has actually purchased this product
	var purchased int
	database.DB.QueryRow(`
		SELECT COUNT(*) FROM order_items oi
		JOIN orders o          ON oi.order_id = o.id
		JOIN product_variants pv ON oi.variant_id = pv.id
		WHERE o.user_id = $1
		AND pv.product_id = $2
		AND o.status = 'delivered'`,
		userID, req.ProductID,
	).Scan(&purchased)

	if purchased == 0 {
		return nil, fmt.Errorf("you can only review products you have purchased")
	}

	var review models.Review
	err := database.DB.QueryRow(`
		INSERT INTO reviews (user_id, product_id, rating, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, product_id, rating, comment, is_approved, created_at`,
		userID, req.ProductID, req.Rating, req.Comment,
	).Scan(
		&review.ID, &review.UserID, &review.ProductID,
		&review.Rating, &review.Comment, &review.IsApproved, &review.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return &review, nil
}

// GetReviewsByProductID fetches all approved reviews for a product
func GetReviewsByProductID(productID string) ([]models.Review, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.rating, r.comment,
		       r.is_approved, r.created_at, u.name, u.avatar_url
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = $1 AND r.is_approved = TRUE
		ORDER BY r.created_at DESC`, productID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var r models.Review
		var u models.User
		rows.Scan(
			&r.ID, &r.UserID, &r.ProductID, &r.Rating,
			&r.Comment, &r.IsApproved, &r.CreatedAt,
			&u.Name, &u.AvatarURL,
		)
		r.User = &u
		reviews = append(reviews, r)
	}

	return reviews, nil
}

// GetPendingReviews fetches all unapproved reviews for the admin panel
func GetPendingReviews() ([]models.Review, error) {
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.rating, r.comment,
		       r.is_approved, r.created_at, u.name, u.avatar_url
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.is_approved = FALSE
		ORDER BY r.created_at ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending reviews: %w", err)
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var r models.Review
		var u models.User
		rows.Scan(
			&r.ID, &r.UserID, &r.ProductID, &r.Rating,
			&r.Comment, &r.IsApproved, &r.CreatedAt,
			&u.Name, &u.AvatarURL,
		)
		r.User = &u
		reviews = append(reviews, r)
	}

	return reviews, nil
}

// ApproveReview sets a review's is_approved flag to true or false
func ApproveReview(reviewID string, isApproved bool) error {
	_, err := database.DB.Exec(`
		UPDATE reviews SET is_approved = $1 WHERE id = $2`,
		isApproved, reviewID,
	)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}
	return nil
}

// DeleteReview permanently removes a review from the database
func DeleteReview(reviewID, userID string) error {
	var id string
	err := database.DB.QueryRow(`
		SELECT id FROM reviews WHERE id = $1 AND user_id = $2`,
		reviewID, userID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return fmt.Errorf("review not found or does not belong to you")
	}

	_, err = database.DB.Exec(`DELETE FROM reviews WHERE id = $1`, reviewID)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}
