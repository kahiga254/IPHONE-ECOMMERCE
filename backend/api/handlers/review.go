// api/handlers/review.go (partial - key functions)
package handlers

import (
	"net/http"
	"strconv"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateReview handler
func CreateReview(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	review, err := services.CreateReview(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "review created successfully", review)
}

// GetProductReviews handler
func GetProductReviews(c *gin.Context) {
	productID := c.Param("product_id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reviews, total, err := services.GetReviewsByProductID(productID, limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", gin.H{
		"reviews": reviews,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// DeleteReview handler
func DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("user_role")

	err := services.DeleteReview(reviewID, userID, role)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "review deleted successfully", nil)
}

// GetPendingReviews handler (admin only)
func GetPendingReviews(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	reviews, total, err := services.GetPendingReviews(limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", gin.H{
		"reviews": reviews,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
	})
}

// ApproveReview handler (admin only)
func ApproveReview(c *gin.Context) {
	reviewID := c.Param("id")

	err := services.ApproveReview(reviewID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "review approved successfully", nil)
}
