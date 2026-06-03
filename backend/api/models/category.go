package models

import "time"

type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	ParentID    *string   `json:"parent_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateCategoryRequest is used when an admin creates a new category
type CreateCategoryRequest struct {
	Name        string  `json:"name"        binding:"required,min=2,max=100"`
	Slug        string  `json:"slug"        binding:"required"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	ParentID    *string `json:"parent_id"`
}

// UpdateCategoryRequest is used when an admin updates an existing category
type UpdateCategoryRequest struct {
	Name        string  `json:"name"        binding:"omitempty,min=2,max=100"`
	Slug        string  `json:"slug"        binding:"omitempty"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	ParentID    *string `json:"parent_id"`
}
