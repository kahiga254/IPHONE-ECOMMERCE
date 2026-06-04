package services

import (
	"fmt"
	"strings"

	"backend/api/models"
	"backend/api/repository"
)

// GetAllProducts fetches a paginated and filtered list of products
func GetAllProducts(q models.ProductFilterQuery) (*models.PaginatedResponse, error) {
	// Sanitize page and limit
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 12
	}

	products, total, err := repository.GetAllProducts(q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	totalPages := (total + q.Limit - 1) / q.Limit

	return &models.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetProductBySlug fetches a single product by its slug
func GetProductBySlug(slug string) (*models.Product, error) {
	product, err := repository.GetProductBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}
	return product, nil
}

// CreateProduct validates and creates a new product
func CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	// Normalize slug to lowercase with hyphens
	req.Slug = normalizeSlug(req.Slug)

	// Validate discount price is less than base price
	if req.DiscountPrice != nil && *req.DiscountPrice >= req.BasePrice {
		return nil, fmt.Errorf("discount price must be less than base price")
	}

	// Validate all variant SKUs are unique within the request
	skus := map[string]bool{}
	for _, v := range req.Variants {
		if skus[v.SKU] {
			return nil, fmt.Errorf("duplicate SKU found: %s", v.SKU)
		}
		skus[v.SKU] = true
	}

	product, err := repository.CreateProduct(req)
	if err != nil {
		// Check for duplicate slug error from postgres
		if strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("a product with this slug already exists")
		}
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// UpdateProduct validates and updates an existing product
func UpdateProduct(id string, req models.UpdateProductRequest) error {
	// Check product exists
	existing, err := repository.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	// Validate discount price if provided
	if req.DiscountPrice != nil && *req.DiscountPrice >= req.BasePrice {
		return fmt.Errorf("discount price must be less than base price")
	}

	return repository.UpdateProduct(id, req)
}

// DeleteProduct soft deletes a product by ID
func DeleteProduct(id string) error {
	existing, err := repository.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	return repository.DeleteProduct(id)
}

// GetAllCategories fetches all product categories with pagination
// GetAllCategories fetches all product categories with pagination
func GetAllCategories(page, limit int) ([]models.Category, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	categories, total, err := repository.GetAllCategories(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return categories, total, nil
}

// CreateCategory validates and creates a new category
func CreateCategory(req models.CreateCategoryRequest) (*models.Category, error) {
	req.Slug = normalizeSlug(req.Slug)

	// Convert CreateCategoryRequest to Category model
	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ImageURL:    req.ImageURL,
		ParentID:    req.ParentID,
	}

	err := repository.CreateCategory(category)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("a category with this slug already exists")
		}
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// UpdateCategory validates and updates an existing category
func UpdateCategory(id string, req models.UpdateCategoryRequest) error {
	existing, err := repository.GetCategoryByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch category: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("category not found")
	}

	// Update only provided fields
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Slug != "" {
		existing.Slug = normalizeSlug(req.Slug)
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.ImageURL != "" {
		existing.ImageURL = req.ImageURL
	}
	if req.ParentID != nil {
		existing.ParentID = req.ParentID
	}

	return repository.UpdateCategory(id, existing)
}

// DeleteCategory deletes a category by ID
func DeleteCategory(id string) error {
	existing, err := repository.GetCategoryByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch category: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("category not found")
	}

	return repository.DeleteCategory(id)
}

// ─── Private Helpers ──────────────────────────────────────────────────────────

// normalizeSlug converts a slug to lowercase and replaces spaces with hyphens
func normalizeSlug(slug string) string {
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
