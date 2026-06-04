package handlers

import (
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateReview godoc
// POST /api/v1/reviews
func CreateReview(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	review, err := services.CreateReview(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "review submitted successfully, it will appear after approval", review)
}

// GetProductReviews godoc
// GET /api/v1/reviews/:product_id
func GetProductReviews(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing product id")
		return
	}

	reviews, err := services.GetReviewsByProductID(productID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", reviews)
}

// DeleteReview godoc
// DELETE /api/v1/reviews/:id
func DeleteReview(c *gin.Context) {
	userID := c.GetString("user_id")
	reviewID := c.Param("id")

	if reviewID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing review id")
		return
	}

	if err := services.DeleteReview(reviewID, userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "review deleted successfully", nil)
}

// GetPendingReviews godoc
// GET /api/v1/admin/reviews/pending
func GetPendingReviews(c *gin.Context) {
	reviews, err := services.GetPendingReviews()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", reviews)
}

// ApproveReview godoc
// PATCH /api/v1/admin/reviews/:id/approve
func ApproveReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing review id")
		return
	}

	var req models.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.ApproveReview(reviewID, req.IsApproved); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	message := "review approved successfully"
	if !req.IsApproved {
		message = "review rejected successfully"
	}

	utils.SuccessResponse(c, http.StatusOK, message, nil)
}
