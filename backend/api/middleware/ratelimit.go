package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/adams/applestore-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// visitor holds the request count and the time the window started for a single IP
type visitor struct {
	count       int
	windowStart time.Time
	blocked     bool
	blockedAt   time.Time
}

// store holds all visitors in memory with a mutex for safe concurrent access
var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

// RateLimit returns a middleware that limits requests per IP within a time window
// maxRequests — maximum number of requests allowed within the window
// window      — duration of the time window e.g 1 minute
// blockFor    — how long to block the IP after exceeding the limit
func RateLimit(maxRequests int, window, blockFor time.Duration) gin.HandlerFunc {
	// Start a background goroutine to clean up stale visitors every 5 minutes
	go cleanupVisitors(window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{
				count:       1,
				windowStart: time.Now(),
			}
			mu.Unlock()
			c.Next()
			return
		}

		// Check if the IP is currently blocked
		if v.blocked {
			if time.Since(v.blockedAt) < blockFor {
				mu.Unlock()
				utils.ErrorResponse(c, http.StatusTooManyRequests, "too many requests, please try again later")
				c.Abort()
				return
			}
			// Block duration has passed — reset the visitor
			v.blocked = false
			v.count = 1
			v.windowStart = time.Now()
			mu.Unlock()
			c.Next()
			return
		}

		// Check if the window has expired — reset the count
		if time.Since(v.windowStart) > window {
			v.count = 1
			v.windowStart = time.Now()
			mu.Unlock()
			c.Next()
			return
		}

		// Increment the request count
		v.count++

		// Block the IP if the limit is exceeded
		if v.count > maxRequests {
			v.blocked = true
			v.blockedAt = time.Now()
			mu.Unlock()
			utils.ErrorResponse(c, http.StatusTooManyRequests, "too many requests, please try again later")
			c.Abort()
			return
		}

		mu.Unlock()
		c.Next()
	}
}

// AuthRateLimit is a preset rate limiter specifically for auth routes
// Allows 10 requests per minute, blocks for 15 minutes after exceeding the limit
func AuthRateLimit() gin.HandlerFunc {
	return RateLimit(10, 1*time.Minute, 15*time.Minute)
}

// OTPRateLimit is a stricter rate limiter for OTP sending
// Allows 3 requests per 5 minutes, blocks for 30 minutes after exceeding the limit
func OTPRateLimit() gin.HandlerFunc {
	return RateLimit(3, 5*time.Minute, 30*time.Minute)
}

// cleanupVisitors removes stale visitor entries from memory every 5 minutes
// This prevents the visitors map from growing indefinitely
func cleanupVisitors(window time.Duration) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		mu.Lock()
		for ip, v := range visitors {
			if !v.blocked && time.Since(v.windowStart) > window*2 {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
