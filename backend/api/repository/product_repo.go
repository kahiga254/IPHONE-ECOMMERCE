package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"backend/api/models"
	"backend/pkg/database"
)

// GetAllProducts fetches a paginated and filtered list of products
func GetAllProducts(q models.ProductFilterQuery) ([]models.Product, int, error) {
	where := "WHERE p.is_active = TRUE"
	args := []interface{}{}
	idx := 1

	if q.Search != "" {
		where += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", idx, idx)
		args = append(args, "%"+q.Search+"%")
		idx++
	}
	if q.CategorySlug != "" {
		where += fmt.Sprintf(" AND c.slug = $%d", idx)
		args = append(args, q.CategorySlug)
		idx++
	}
	if q.MinPrice > 0 {
		where += fmt.Sprintf(" AND p.base_price >= $%d", idx)
		args = append(args, q.MinPrice)
		idx++
	}
	if q.MaxPrice > 0 {
		where += fmt.Sprintf(" AND p.base_price <= $%d", idx)
		args = append(args, q.MaxPrice)
		idx++
	}
	if q.Featured {
		where += " AND p.is_featured = TRUE"
	}

	// Count total matching products for pagination
	var total int
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		` + where
	if err := database.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	// Allowed sort columns to prevent SQL injection
	sortColumn := "p.created_at"
	if q.SortBy == "price" {
		sortColumn = "p.base_price"
	} else if q.SortBy == "name" {
		sortColumn = "p.name"
	}

	sortOrder := "DESC"
	if q.Order == "asc" {
		sortOrder = "ASC"
	}

	offset := (q.Page - 1) * q.Limit
	args = append(args, q.Limit, offset)

	query := fmt.Sprintf(`
		SELECT p.id, p.name, p.slug, p.description, p.category_id, p.base_price,
		       p.discount_price, p.is_featured, p.is_active, p.meta_title,
		       p.meta_description, p.created_at, p.updated_at,
		       COALESCE(AVG(r.rating), 0) AS avg_rating
		FROM products p
		LEFT JOIN categories c  ON p.category_id = c.id
		LEFT JOIN reviews r     ON r.product_id = p.id AND r.is_approved = TRUE
		%s
		GROUP BY p.id
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		where, sortColumn, sortOrder, idx, idx+1,
	)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch products: %w", err)
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Slug, &p.Description, &p.CategoryID,
			&p.BasePrice, &p.DiscountPrice, &p.IsFeatured, &p.IsActive,
			&p.MetaTitle, &p.MetaDescription, &p.CreatedAt, &p.UpdatedAt,
			&p.AvgRating,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		p.Variants = getVariantsByProductID(p.ID)
		products = append(products, p)
	}

	return products, total, nil
}

// GetProductBySlug fetches a single product with all its variants, specs and reviews
func GetProductBySlug(slug string) (*models.Product, error) {
	var p models.Product

	err := database.DB.QueryRow(`
		SELECT p.id, p.name, p.slug, p.description, p.category_id, p.base_price,
		       p.discount_price, p.is_featured, p.is_active, p.meta_title,
		       p.meta_description, p.created_at, p.updated_at,
		       COALESCE(AVG(r.rating), 0) AS avg_rating
		FROM products p
		LEFT JOIN reviews r ON r.product_id = p.id AND r.is_approved = TRUE
		WHERE p.slug = $1 AND p.is_active = TRUE
		GROUP BY p.id`, slug,
	).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.CategoryID,
		&p.BasePrice, &p.DiscountPrice, &p.IsFeatured, &p.IsActive,
		&p.MetaTitle, &p.MetaDescription, &p.CreatedAt, &p.UpdatedAt,
		&p.AvgRating,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	p.Variants = getVariantsByProductID(p.ID)
	p.Specs = getSpecsByProductID(p.ID)

	return &p, nil
}

// GetProductByID fetches a single product by its UUID
func GetProductByID(id string) (*models.Product, error) {
	var p models.Product

	err := database.DB.QueryRow(`
		SELECT id, name, slug, description, category_id, base_price,
		       discount_price, is_featured, is_active, meta_title,
		       meta_description, created_at, updated_at
		FROM products WHERE id = $1`, id,
	).Scan(
		&p.ID, &p.Name, &p.Slug, &p.Description, &p.CategoryID,
		&p.BasePrice, &p.DiscountPrice, &p.IsFeatured, &p.IsActive,
		&p.MetaTitle, &p.MetaDescription, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}

	return &p, nil
}

// CreateProduct inserts a new product and returns the created product
func CreateProduct(req models.CreateProductRequest) (*models.Product, error) {
	var p models.Product

	err := database.DB.QueryRow(`
		INSERT INTO products (name, slug, description, category_id, base_price, discount_price, is_featured, meta_title, meta_description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, slug, base_price, created_at, updated_at`,
		req.Name, req.Slug, req.Description, req.CategoryID, req.BasePrice,
		req.DiscountPrice, req.IsFeatured, req.MetaTitle, req.MetaDescription,
	).Scan(&p.ID, &p.Name, &p.Slug, &p.BasePrice, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Insert variants
	for _, v := range req.Variants {
		imagesJSON, _ := json.Marshal(v.Images)
		_, err := database.DB.Exec(`
			INSERT INTO product_variants (product_id, sku, color, storage, price, stock, images)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			p.ID, v.SKU, v.Color, v.Storage, v.Price, v.Stock, imagesJSON,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant %s: %w", v.SKU, err)
		}
	}

	// Insert specs
	for _, s := range req.Specs {
		_, err := database.DB.Exec(`
			INSERT INTO product_specs (product_id, spec_key, spec_value)
			VALUES ($1, $2, $3)`,
			p.ID, s.Key, s.Value,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create spec %s: %w", s.Key, err)
		}
	}

	return &p, nil
}

// UpdateProduct updates an existing product's fields
func UpdateProduct(id string, req models.UpdateProductRequest) error {
	_, err := database.DB.Exec(`
		UPDATE products SET
			name             = $1,
			description      = $2,
			base_price       = $3,
			discount_price   = $4,
			is_featured      = $5,
			is_active        = $6,
			meta_title       = $7,
			meta_description = $8,
			updated_at       = NOW()
		WHERE id = $9`,
		req.Name, req.Description, req.BasePrice, req.DiscountPrice,
		req.IsFeatured, req.IsActive, req.MetaTitle, req.MetaDescription, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	return nil
}

// DeleteProduct soft deletes a product by setting is_active to false
func DeleteProduct(id string) error {
	_, err := database.DB.Exec(`UPDATE products SET is_active = FALSE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}

// ─── Private Helpers ──────────────────────────────────────────────────────────

// getVariantsByProductID fetches all variants for a given product
func getVariantsByProductID(productID string) []models.Variant {
	rows, err := database.DB.Query(`
		SELECT id, product_id, sku, color, storage, price, stock, images, created_at
		FROM product_variants WHERE product_id = $1`, productID,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	variants := []models.Variant{}
	for rows.Next() {
		var v models.Variant
		var imagesJSON []byte
		rows.Scan(
			&v.ID, &v.ProductID, &v.SKU, &v.Color,
			&v.Storage, &v.Price, &v.Stock, &imagesJSON, &v.CreatedAt,
		)
		json.Unmarshal(imagesJSON, &v.Images)
		variants = append(variants, v)
	}
	return variants
}

// getSpecsByProductID fetches all specs for a given product
func getSpecsByProductID(productID string) []models.Spec {
	rows, err := database.DB.Query(`
		SELECT id, product_id, spec_key, spec_value
		FROM product_specs WHERE product_id = $1`, productID,
	)
	if err != nil {
		return nil
	}
	defer rows.Close()

	specs := []models.Spec{}
	for rows.Next() {
		var s models.Spec
		rows.Scan(&s.ID, &s.ProductID, &s.Key, &s.Value)
		specs = append(specs, s)
	}
	return specs
}

// GetVariantByID fetches a single variant by its UUID
func GetVariantByID(id string) (*models.Variant, error) {
	var v models.Variant
	var imagesJSON []byte

	err := database.DB.QueryRow(`
		SELECT id, product_id, sku, color, storage, price, stock, images, created_at
		FROM product_variants WHERE id = $1`, id,
	).Scan(
		&v.ID, &v.ProductID, &v.SKU, &v.Color,
		&v.Storage, &v.Price, &v.Stock, &imagesJSON, &v.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get variant by id: %w", err)
	}

	json.Unmarshal(imagesJSON, &v.Images)
	return &v, nil
}
