package handlers

import (
	"fmt"
	"net/http"

	"backend/api/auth"
	"backend/api/models"
	"backend/api/repository"
	"backend/api/services"
	"backend/config"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Register godoc
// POST /api/v1/auth/register
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := services.Register(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "account created successfully", resp)
}

// Login godoc
// POST /api/v1/auth/login
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := services.Login(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "login successful", resp)
}

// SendOTP godoc
// POST /api/v1/auth/otp/send
func SendOTP(c *gin.Context) {
	var req models.PhoneLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.LoginWithPhone(req.Phone); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP sent successfully", nil)
}

// VerifyOTP godoc
// POST /api/v1/auth/otp/verify
func VerifyOTP(c *gin.Context) {
	var req models.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := services.VerifyPhoneOTP(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "OTP verified successfully", resp)
}

// GoogleLogin godoc
// GET /api/v1/auth/google
func GoogleLogin(c *gin.Context) {
	// Generate a random state string to prevent CSRF
	state := "random-state-string" // TODO: generate and store in session
	url := auth.GetGoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback godoc
// GET /api/v1/auth/google/callback
func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing authorization code")
		return
	}

	resp, err := services.GoogleLogin(code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Redirect to frontend with tokens as query params
	redirectURL := fmt.Sprintf(
		"%s/auth/callback?access_token=%s&refresh_token=%s",
		config.App.FrontendURL,
		resp.AccessToken,
		resp.RefreshToken,
	)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// RefreshToken godoc
// POST /api/v1/auth/refresh
func RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := services.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "token refreshed successfully", resp)
}

// Logout godoc
// POST /api/v1/auth/logout
func Logout(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.Logout(req.RefreshToken); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "logged out successfully", nil)
}

// VerifyEmail godoc
// GET /api/v1/auth/verify-email/:token
func VerifyEmail(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing token")
		return
	}

	if err := services.VerifyEmail(token); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "email verified successfully", nil)
}

// GetMe godoc
// GET /api/v1/auth/me
func GetMe(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := repository.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		utils.ErrorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	user.PasswordHash = nil
	utils.SuccessResponse(c, http.StatusOK, "", user)
}