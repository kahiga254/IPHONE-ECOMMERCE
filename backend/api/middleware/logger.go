package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a middleware that logs every incoming request
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Record the time the request came in
		start := time.Now()

		// Get request details before the handler runs
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		// Let the handler run
		c.Next()

		// After the handler — collect response details
		statusCode := c.Writer.Status()
		latency := time.Since(start)
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Pick a color based on the status code
		statusColor := colorForStatus(statusCode)
		methodColor := colorForMethod(method)
		resetColor := "\033[0m"

		// Format and print the log line
		fmt.Printf(
			"%s | %s%d%s | %s%-7s%s | %s | %v %s\n",
			start.Format("2006/01/02 15:04:05"),
			statusColor, statusCode, resetColor,
			methodColor, method, resetColor,
			clientIP,
			latency,
			path,
		)

		// Print any error messages on a separate line
		if errorMessage != "" {
			fmt.Printf("  └─ error: %s\n", errorMessage)
		}
	}
}

// colorForStatus returns an ANSI color code based on the HTTP status code
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[32m" // green for success
	case code >= 300 && code < 400:
		return "\033[36m" // cyan for redirects
	case code >= 400 && code < 500:
		return "\033[33m" // yellow for client errors
	default:
		return "\033[31m" // red for server errors
	}
}

// colorForMethod returns an ANSI color code based on the HTTP method
func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[34m" // blue
	case "POST":
		return "\033[32m" // green
	case "PUT", "PATCH":
		return "\033[33m" // yellow
	case "DELETE":
		return "\033[31m" // red
	default:
		return "\033[37m" // white
	}
}
