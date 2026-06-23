// config/config.go
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig holds all application configuration
type AppConfig struct {
	// Server
	Port string
	Env  string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret string

	// Frontend
	FrontendURL string

	// Google OAuth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	// M-Pesa
	MPesaConsumerKey    string
	MPesaConsumerSecret string
	MPesaPasskey        string
	MPesaShortcode      string
	MPesaStoreNumber    string
	MPesaEnvironment    string
	MPesaCallbackURL    string // ← ADD THIS LINE

	// Africa's Talking SMS
	ATSAPIKey   string
	ATSUsername string
	ATSSenderID string
}

var App *AppConfig

// Load loads environment variables and initializes config
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  No .env file found, using system environment variables")
	}

	App = &AppConfig{
		// Server
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "ecommerce_db"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this"),

		// Frontend
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),

		// Google OAuth
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/api/v1/auth/google/callback"),

		// M-Pesa
		MPesaConsumerKey:    getEnv("MPESA_CONSUMER_KEY", ""),
		MPesaConsumerSecret: getEnv("MPESA_CONSUMER_SECRET", ""),
		MPesaPasskey:        getEnv("MPESA_PASSKEY", ""),
		MPesaShortcode:      getEnv("MPESA_SHORTCODE", "174379"),
		MPesaStoreNumber:    getEnv("MPESA_STORE_NUMBER", ""),
		MPesaEnvironment:    getEnv("MPESA_ENVIRONMENT", "sandbox"),
		MPesaCallbackURL:    getEnv("MPESA_CALLBACK_URL", "https://yourdomain.com/api/v1/payments/mpesa/callback"), // ADD THIS LINE

		// Africa's Talking
		ATSAPIKey:   getEnv("ATS_API_KEY", ""),
		ATSUsername: getEnv("ATS_USERNAME", ""),
		ATSSenderID: getEnv("ATS_SENDER_ID", "Ecommerce"),
	}

	log.Println("✅ Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
