package handlers

import (
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetWishlist godoc
// GET /api/v1/wishlist
func GetWishlist(c *gin.Context) {
	userID := c.GetString("user_id")

	wishlist, err := services.GetWishlist(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", wishlist)
}

// AddToWishlist godoc
// POST /api/v1/wishlist
func AddToWishlist(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.AddToWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	wishlist, err := services.AddToWishlist(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "item added to wishlist", wishlist)
}

// RemoveFromWishlist godoc
// DELETE /api/v1/wishlist/:variant_id
func RemoveFromWishlist(c *gin.Context) {
	userID := c.GetString("user_id")
	variantID := c.Param("variant_id")

	if variantID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing variant id")
		return
	}

	if err := services.RemoveFromWishlist(userID, variantID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "item removed from wishlist", nil)
}

// ClearWishlist godoc
// DELETE /api/v1/wishlist
func ClearWishlist(c *gin.Context) {
	userID := c.GetString("user_id")

	if err := services.ClearWishlist(userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "wishlist cleared successfully", nil)
}
