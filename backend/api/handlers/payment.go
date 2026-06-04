package handlers

import (
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// InitiatePayment godoc
// POST /api/v1/payments/mpesa/stkpush
func InitiatePayment(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	payment, err := services.InitiatePayment(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "STK push sent to your phone, enter your M-Pesa PIN to complete payment", payment)
}

// MpesaCallback godoc
// POST /api/v1/payments/mpesa/callback
// This endpoint is called by Daraja — it must never require authentication
func MpesaCallback(c *gin.Context) {
	var callback models.MpesaCallback
	if err := c.ShouldBindJSON(&callback); err != nil {
		// Always respond with 200 to Daraja even on error
		// If we return non-200 Daraja will keep retrying the callback
		c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Accepted"})
		return
	}

	if err := services.HandleMpesaCallback(callback); err != nil {
		// Log the error but still respond 200 to Daraja
		c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Accepted"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ResultCode": 0, "ResultDesc": "Accepted"})
}

// QueryPaymentStatus godoc
// GET /api/v1/payments/:order_id/status
func QueryPaymentStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	orderID := c.Param("order_id")

	if orderID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing order id")
		return
	}

	payment, err := services.QueryPaymentStatus(orderID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", payment)
}
