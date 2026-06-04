// api/middleware/cors.go
package middleware

import (
	"time"

	"backend/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a configured CORS middleware for the Gin router
func CORS() gin.HandlerFunc {
	// Allow all origins in development, specific in production
	allowOrigins := []string{config.App.FrontendURL}
	if config.App.Env == "development" {
		allowOrigins = []string{"*"}
	}

	return cors.New(cors.Config{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
