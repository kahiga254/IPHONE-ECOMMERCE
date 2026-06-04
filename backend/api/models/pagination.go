// api/models/pagination.go
package models

// PaginatedResponse is a generic wrapper for paginated API responses
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Limit      int         `json:"limit"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
}

// PaginationParams holds pagination query parameters
type PaginationParams struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

// GetOffset calculates the SQL OFFSET value
func (p *PaginationParams) GetOffset() int {
	if p.Page < 1 {
		p.Page = 1
	}
	return (p.Page - 1) * p.Limit
}

// GetLimit returns the limit with a maximum cap
func (p *PaginationParams) GetLimit() int {
	if p.Limit < 1 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	return p.Limit
}
