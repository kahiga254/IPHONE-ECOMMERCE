package models

import "time"

type Product struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Slug            string    `json:"slug"`
	Description     string    `json:"description"`
	CategoryID      *string   `json:"category_id,omitempty"`
	BasePrice       float64   `json:"base_price"`
	DiscountPrice   *float64  `json:"discount_price,omitempty"`
	IsFeatured      bool      `json:"is_featured"`
	IsActive        bool      `json:"is_active"`
	MetaTitle       string    `json:"meta_title"`
	MetaDescription string    `json:"meta_description"`
	AvgRating       float64   `json:"avg_rating"`
	Variants        []Variant `json:"variants,omitempty"`
	Specs           []Spec    `json:"specs,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Variant struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	SKU       string    `json:"sku"`
	Color     string    `json:"color"`
	Storage   string    `json:"storage"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	Images    []string  `json:"images"`
	CreatedAt time.Time `json:"created_at"`
}

type Spec struct {
	ID        string `json:"id"`
	ProductID string `json:"product_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// CreateProductRequest is used when an admin creates a new product
type CreateProductRequest struct {
	Name            string                 `json:"name"             binding:"required,min=2,max=255"`
	Slug            string                 `json:"slug"             binding:"required"`
	Description     string                 `json:"description"`
	CategoryID      *string                `json:"category_id"`
	BasePrice       float64                `json:"base_price"       binding:"required,gt=0"`
	DiscountPrice   *float64               `json:"discount_price"   binding:"omitempty,gt=0"`
	IsFeatured      bool                   `json:"is_featured"`
	MetaTitle       string                 `json:"meta_title"`
	MetaDescription string                 `json:"meta_description"`
	Variants        []CreateVariantRequest `json:"variants"         binding:"required,min=1"`
	Specs           []CreateSpecRequest    `json:"specs"`
}

type CreateVariantRequest struct {
	SKU     string   `json:"sku"     binding:"required"`
	Color   string   `json:"color"   binding:"required"`
	Storage string   `json:"storage" binding:"required"`
	Price   float64  `json:"price"   binding:"required,gt=0"`
	Stock   int      `json:"stock"   binding:"required,min=0"`
	Images  []string `json:"images"  binding:"required,min=1"`
}

type CreateSpecRequest struct {
	Key   string `json:"key"   binding:"required"`
	Value string `json:"value" binding:"required"`
}

// UpdateProductRequest is used when an admin updates an existing product
type UpdateProductRequest struct {
	Name            string   `json:"name"             binding:"omitempty,min=2,max=255"`
	Description     string   `json:"description"`
	BasePrice       float64  `json:"base_price"       binding:"omitempty,gt=0"`
	DiscountPrice   *float64 `json:"discount_price"   binding:"omitempty,gt=0"`
	IsFeatured      bool     `json:"is_featured"`
	IsActive        bool     `json:"is_active"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
}

// ProductFilterQuery holds all query params for filtering the product list
type ProductFilterQuery struct {
	Page         int     `form:"page,default=1"`
	Limit        int     `form:"limit,default=12"`
	Search       string  `form:"search"`
	CategorySlug string  `form:"category"`
	MinPrice     float64 `form:"min_price"`
	MaxPrice     float64 `form:"max_price"`
	Color        string  `form:"color"`
	Storage      string  `form:"storage"`
	Featured     bool    `form:"featured"`
	SortBy       string  `form:"sort_by,default=created_at"`
	Order        string  `form:"order,default=desc"`
}
