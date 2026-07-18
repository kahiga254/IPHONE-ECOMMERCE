package services

import (
	"fmt"
	"strings"
	"time"

	"backend/api/models"
	"backend/api/repository"
	"backend/pkg/cache"
)

// GetAllProducts fetches a paginated and filtered list of products
func GetAllProducts(q models.ProductFilterQuery) (*models.PaginatedResponse, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 12
	}

	// Build cache key from query params
	cacheKey := fmt.Sprintf("products:page=%d:limit=%d:search=%s:category=%s:min=%.0f:max=%.0f:sort=%s:order=%s",
		q.Page, q.Limit, q.Search, q.CategorySlug, q.MinPrice, q.MaxPrice, q.SortBy, q.Order)

	// Return cached response if available
	if cached, found := cache.Cache.Get(cacheKey); found {
		return cached.(*models.PaginatedResponse), nil
	}

	products, total, err := repository.GetAllProducts(q)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	totalPages := (total + q.Limit - 1) / q.Limit

	response := &models.PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       q.Page,
		Limit:      q.Limit,
		TotalPages: totalPages,
	}

	// Cache for 5 minutes
	cache.Cache.Set(cacheKey, response, 5*time.Minute)

	return response, nil
}

// GetProductBySlug fetches a single product by its slug
func GetProductBySlug(slug string) (*models.Product, error) {
	cacheKey := "product:" + slug

	if cached, found := cache.Cache.Get(cacheKey); found {
		return cached.(*models.Product), nil
	}

	product, err := repository.GetProductBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	if product == nil {
		return nil, fmt.Errorf("product not found")
	}

	// Cache for 5 minutes
	cache.Cache.Set(cacheKey, product, 5*time.Minute)

	return product, nil
}

// CreateProduct validates and creates a new product
func CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	req.Slug = normalizeSlug(req.Slug)

	if req.DiscountPrice != nil && *req.DiscountPrice >= req.BasePrice {
		return nil, fmt.Errorf("discount price must be less than base price")
	}

	skus := map[string]bool{}
	for _, v := range req.Variants {
		if skus[v.SKU] {
			return nil, fmt.Errorf("duplicate SKU found: %s", v.SKU)
		}
		skus[v.SKU] = true
	}

	product, err := repository.CreateProduct(req)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return nil, fmt.Errorf("a product with this slug already exists")
		}
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Invalidate products cache
	cache.Cache.DeleteByPrefix("products:")

	return product, nil
}

// UpdateProduct validates and updates an existing product
func UpdateProduct(id string, req models.UpdateProductRequest) error {
	existing, err := repository.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	if req.DiscountPrice != nil && *req.DiscountPrice >= req.BasePrice {
		return fmt.Errorf("discount price must be less than base price")
	}

	err = repository.UpdateProduct(id, req)
	if err != nil {
		return err
	}

	// Invalidate cache for this product and all product lists
	cache.Cache.DeleteByPrefix("products:")
	cache.Cache.Delete("product:" + existing.Slug)

	return nil
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

	err = repository.DeleteProduct(id)
	if err != nil {
		return err
	}

	// Invalidate cache
	cache.Cache.DeleteByPrefix("products:")
	cache.Cache.Delete("product:" + existing.Slug)

	return nil
}

// GetAllCategories fetches all product categories with pagination
func GetAllCategories(page, limit int) ([]models.Category, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	cacheKey := fmt.Sprintf("categories:page=%d:limit=%d", page, limit)

	if cached, found := cache.Cache.Get(cacheKey); found {
		data := cached.(map[string]interface{})
		return data["categories"].([]models.Category), data["total"].(int), nil
	}

	offset := (page - 1) * limit
	categories, total, err := repository.GetAllCategories(limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch categories: %w", err)
	}

	// Cache categories for 30 minutes
	cache.Cache.Set(cacheKey, map[string]interface{}{
		"categories": categories,
		"total":      total,
	}, 30*time.Minute)

	return categories, total, nil
}

// CreateCategory validates and creates a new category
func CreateCategory(req models.CreateCategoryRequest) (*models.Category, error) {
	req.Slug = normalizeSlug(req.Slug)

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

	// Invalidate categories cache
	cache.Cache.DeleteByPrefix("categories:")

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

	err = repository.UpdateCategory(id, existing)
	if err != nil {
		return err
	}

	// Invalidate categories cache
	cache.Cache.DeleteByPrefix("categories:")

	return nil
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

	err = repository.DeleteCategory(id)
	if err != nil {
		return err
	}

	cache.Cache.DeleteByPrefix("categories:")

	return nil
}

// ─── Private Helpers ──────────────────────────────────────────────────────────

func normalizeSlug(slug string) string {
	slug = strings.TrimSpace(slug)
	slug = strings.ToLower(slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}

func ToggleProductStatus(id string, isActive bool) error {
	existing, err := repository.GetProductByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("product not found")
	}

	err = repository.UpdateProductStatus(id, isActive)
	if err != nil {
		return err
	}

	// Invalidate cache
	cache.Cache.DeleteByPrefix("products:")
	cache.Cache.Delete("product:" + existing.Slug)

	return nil
}
