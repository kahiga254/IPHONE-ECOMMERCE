package utils

import (
	"github.com/gin-gonic/gin"
)

// APIResponse is the standard response shape for every endpoint
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a 200 or custom status with data
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response with a status code and message
func ErrorResponse(c *gin.Context, statusCode int, errMessage string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   errMessage,
	})
}

// PaginatedResponse is the standard shape for paginated list responses
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}
