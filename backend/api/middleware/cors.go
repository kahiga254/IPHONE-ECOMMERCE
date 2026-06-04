package middleware

import (
	"time"

	"github.com/adams/applestore-backend/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS returns a configured CORS middleware for the Gin router
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// Only allow requests from our Next.js frontend
		AllowOrigins: []string{
			config.App.FrontendURL,
		},

		// HTTP methods the frontend is allowed to use
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		// Headers the frontend is allowed to send
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
		},

		// Headers the frontend is allowed to read from responses
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},

		// Allow the frontend to send cookies and auth headers
		AllowCredentials: true,

		// Cache preflight response for 12 hours
		// Browser will not send a preflight OPTIONS request again for 12 hours
		MaxAge: 12 * time.Hour,
	})
}
