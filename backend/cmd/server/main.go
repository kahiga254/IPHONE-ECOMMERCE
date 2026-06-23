package main

import (
	"log"

	"backend/api/handlers"
	"backend/api/middleware"
	"backend/config"
	"backend/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	// Step 1: Load environment variables
	config.Load()

	// Step 2: Connect to PostgreSQL and run migrations
	database.Connect()
	database.Migrate()

	// Step 3: Set Gin mode based on environment
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Step 4: Create the Gin router
	r := gin.New()

	// Step 5: Register global middleware
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery()) // recovers from panics and returns 500

	// ─── API Version 1 ────────────────────────────────────────────────────────
	v1 := r.Group("/api/v1")

	// ─── Public Auth Routes ───────────────────────────────────────────────────
	auth := v1.Group("/auth")
	{
		auth.POST("/register", middleware.AuthRateLimit(), handlers.Register)
		auth.POST("/login", middleware.AuthRateLimit(), handlers.Login)
		auth.POST("/otp/send", middleware.OTPRateLimit(), handlers.SendOTP)
		auth.POST("/otp/verify", middleware.OTPRateLimit(), handlers.VerifyOTP)
		auth.GET("/google", handlers.GoogleLogin)
		auth.GET("/google/callback", handlers.GoogleCallback)
		auth.POST("/refresh", handlers.RefreshToken)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/verify-email/:token", handlers.VerifyEmail)
	}

	// ─── Protected Auth Routes ────────────────────────────────────────────────
	authProtected := v1.Group("/auth")
	authProtected.Use(middleware.AuthRequired())
	{
		authProtected.GET("/me", handlers.GetMe)
	}

	// ─── Public Product Routes ────────────────────────────────────────────────
	products := v1.Group("/products")
	{
		products.GET("", handlers.GetAllProducts)
		products.GET("/:slug", handlers.GetProduct)
	}

	// ─── Public Category Routes ───────────────────────────────────────────────
	categories := v1.Group("/categories")
	{
		categories.GET("", handlers.GetAllCategories)
	}

	// ─── Public Review Routes ─────────────────────────────────────────────────
	reviews := v1.Group("/reviews")
	{
		reviews.GET("/:product_id", handlers.GetProductReviews)
	}

	// ─── Protected User Routes ────────────────────────────────────────────────
	user := v1.Group("")
	user.Use(middleware.AuthRequired())
	{
		// Orders
		user.POST("/orders", handlers.CreateOrder)
		user.GET("/orders", handlers.GetMyOrders)
		user.GET("/orders/:id", handlers.GetOrder)
		user.PATCH("/orders/:id/cancel", handlers.CancelOrder)

		// Payments
		user.POST("/payments/mpesa/stkpush", handlers.InitiatePayment)
		user.GET("/payments/:order_id/status", handlers.QueryPaymentStatus)

		// Wishlist
		user.GET("/wishlist", handlers.GetWishlist)
		user.POST("/wishlist", handlers.AddToWishlist)
		user.DELETE("/wishlist/:variant_id", handlers.RemoveFromWishlist)
		user.DELETE("/wishlist", handlers.ClearWishlist)

		// Addresses
		user.GET("/addresses", handlers.GetAddresses)
		user.POST("/addresses", handlers.CreateAddress)
		user.PUT("/addresses/:id", handlers.UpdateAddress)
		user.DELETE("/addresses/:id", handlers.DeleteAddress)
		user.PATCH("/addresses/:id/default", handlers.SetDefaultAddress)

		// Reviews
		user.POST("/reviews", handlers.CreateReview)
		user.DELETE("/reviews/:id", handlers.DeleteReview)
	}

	v1.POST("/orders/guest", handlers.CreateGuestOrder)
	v1.POST("/payments/mpesa/guest/stkpush", handlers.InitiateGuestPayment)
	v1.POST("/payments/mpesa/callback", handlers.MpesaCallback)
	v1.GET("/payments/guest/:order_id/status", handlers.QueryGuestPaymentStatus)

	// ─── Admin Routes ─────────────────────────────────────────────────────────
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		// Products
		admin.POST("/products", handlers.CreateProduct)
		admin.PUT("/products/:id", handlers.UpdateProduct)
		admin.DELETE("/products/:id", handlers.DeleteProduct)
		admin.PATCH("/products/:id/status", handlers.ToggleProductStatus)
		admin.PUT("/variants/:id", handlers.UpdateVariant)

		// Categories
		admin.POST("/categories", handlers.CreateCategory)
		admin.PUT("/categories/:id", handlers.UpdateCategory)
		admin.DELETE("/categories/:id", handlers.DeleteCategory)

		// Orders
		admin.GET("/orders", handlers.GetAllOrders)
		admin.PATCH("/orders/:id/status", handlers.UpdateOrderStatus)

		// Reviews
		admin.GET("/reviews/pending", handlers.GetPendingReviews)
		admin.PATCH("/reviews/:id/approve", handlers.ApproveReview)
	}

	// Step 6: Start the server
	port := ":" + config.App.Port
	log.Printf("🚀 Server running on port %s", port)
	log.Printf("🌍 Environment: %s", config.App.Env)
	log.Printf("📦 API base URL: http://localhost%s/api/v1", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
