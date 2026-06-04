package models

import "time"

type Address struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Label     string    `json:"label"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	County    string    `json:"county"`
	Town      string    `json:"town"`
	Street    string    `json:"street"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateAddressRequest is used when a user adds a new delivery address
type CreateAddressRequest struct {
	Label     string `json:"label"     binding:"omitempty,oneof=home office other"`
	FullName  string `json:"full_name" binding:"required,min=2,max=100"`
	Phone     string `json:"phone"     binding:"required"`
	County    string `json:"county"    binding:"required"`
	Town      string `json:"town"      binding:"required"`
	Street    string `json:"street"`
	IsDefault bool   `json:"is_default"`
}

// UpdateAddressRequest is used when a user edits an existing address
type UpdateAddressRequest struct {
	Label     string `json:"label"     binding:"omitempty,oneof=home office other"`
	FullName  string `json:"full_name" binding:"omitempty,min=2,max=100"`
	Phone     string `json:"phone"     binding:"omitempty"`
	County    string `json:"county"    binding:"omitempty"`
	Town      string `json:"town"      binding:"omitempty"`
	Street    string `json:"street"`
	IsDefault bool   `json:"is_default"`
}
