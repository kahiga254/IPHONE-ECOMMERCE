// api/repository/category_repo.go
package repository

import (
	"database/sql"
	"fmt"

	"backend/api/models"
	"backend/pkg/database"

	"github.com/google/uuid"
)

// CreateCategory inserts a new category
func CreateCategory(category *models.Category) error {
	if category.ID == "" {
		category.ID = uuid.New().String()
	}

	_, err := database.DB.Exec(`
		INSERT INTO categories (id, name, slug, description, image_url, parent_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		category.ID, category.Name, category.Slug, category.Description, category.ImageURL, category.ParentID,
	)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// GetCategoryByID fetches a category by ID
func GetCategoryByID(id string) (*models.Category, error) {
	var category models.Category
	var parentID sql.NullString
	var description sql.NullString
	var imageURL sql.NullString

	err := database.DB.QueryRow(`
		SELECT id, name, slug, COALESCE(description, ''), COALESCE(image_url, ''), parent_id, created_at
		FROM categories 
		WHERE id = $1`, id,
	).Scan(
		&category.ID, &category.Name, &category.Slug, &description,
		&imageURL, &parentID, &category.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if description.Valid {
		category.Description = description.String
	}
	if imageURL.Valid {
		category.ImageURL = imageURL.String
	}
	if parentID.Valid {
		category.ParentID = &parentID.String
	}

	return &category, nil
}

// GetCategoryBySlug fetches a category by slug
func GetCategoryBySlug(slug string) (*models.Category, error) {
	var category models.Category
	var parentID sql.NullString
	var description sql.NullString
	var imageURL sql.NullString

	err := database.DB.QueryRow(`
		SELECT id, name, slug, COALESCE(description, ''), COALESCE(image_url, ''), parent_id, created_at
		FROM categories 
		WHERE slug = $1`, slug,
	).Scan(
		&category.ID, &category.Name, &category.Slug, &description,
		&imageURL, &parentID, &category.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category by slug: %w", err)
	}

	if description.Valid {
		category.Description = description.String
	}
	if imageURL.Valid {
		category.ImageURL = imageURL.String
	}
	if parentID.Valid {
		category.ParentID = &parentID.String
	}

	return &category, nil
}

// GetAllCategories fetches all categories with pagination
func GetAllCategories(limit, offset int) ([]models.Category, int, error) {
	// Get total count
	var total int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM categories`).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count categories: %w", err)
	}

	// Get paginated categories
	rows, err := database.DB.Query(`
		SELECT id, name, slug, COALESCE(description, ''), COALESCE(image_url, ''), parent_id, created_at
		FROM categories
		ORDER BY name ASC
		LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		var parentID sql.NullString
		var description sql.NullString
		var imageURL sql.NullString

		err := rows.Scan(
			&category.ID, &category.Name, &category.Slug, &description,
			&imageURL, &parentID, &category.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan category: %w", err)
		}

		if description.Valid {
			category.Description = description.String
		}
		if imageURL.Valid {
			category.ImageURL = imageURL.String
		}
		if parentID.Valid {
			category.ParentID = &parentID.String
		}

		categories = append(categories, category)
	}

	return categories, total, nil
}

// UpdateCategory updates an existing category
func UpdateCategory(id string, category *models.Category) error {
	_, err := database.DB.Exec(`
		UPDATE categories 
		SET name = $1, slug = $2, description = $3, image_url = $4, parent_id = $5, updated_at = NOW()
		WHERE id = $6`,
		category.Name, category.Slug, category.Description, category.ImageURL, category.ParentID, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

// DeleteCategory removes a category
func DeleteCategory(id string) error {
	// First check if category has child categories
	var childCount int
	err := database.DB.QueryRow(`SELECT COUNT(*) FROM categories WHERE parent_id = $1`, id).Scan(&childCount)
	if err != nil {
		return fmt.Errorf("failed to check child categories: %w", err)
	}

	if childCount > 0 {
		return fmt.Errorf("cannot delete category with child categories")
	}

	// Delete the category
	_, err = database.DB.Exec(`DELETE FROM categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}
