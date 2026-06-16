package handlers

import (
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func CreateReview(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

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

	utils.SuccessResponse(c, http.StatusCreated, "review submitted", review)
}

func GetProductReviews(c *gin.Context) {
	productID := c.Param("product_id")
	if productID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "product_id required")
		return
	}

	reviews, err := services.GetProductReviews(productID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", reviews)
}

func GetPendingReviews(c *gin.Context) {
	reviews, err := services.GetPendingReviews()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", reviews)
}

func ApproveReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "review id required")
		return
	}

	err := services.ApproveReview(reviewID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "review approved", nil)
}

func DeleteReview(c *gin.Context) {
	reviewID := c.Param("id")
	if reviewID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "review id required")
		return
	}

	err := services.DeleteReview(reviewID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "review deleted", nil)
}
