package middleware

import (
	"strings"

	"backend/api/auth"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequired validates JWT token and adds user info to context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, 401, "missing authorization header")
			c.Abort()
			return
		}

		// Extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, 401, "invalid authorization format. Use: Bearer <token>")
			c.Abort()
			return
		}

		token := parts[1]

		// Validate the token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			utils.ErrorResponse(c, 401, "invalid or expired token")
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminRequired checks if the authenticated user has admin role
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, 403, "access denied")
			c.Abort()
			return
		}

		if role != "admin" {
			utils.ErrorResponse(c, 403, "admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth tries to authenticate but doesn't fail if no token
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		token := parts[1]
		claims, err := auth.ValidateToken(token)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.Role)
		}

		c.Next()
	}
}
