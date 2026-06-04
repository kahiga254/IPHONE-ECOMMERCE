package handlers

import (
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// POST /api/v1/orders
func CreateOrder(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := services.CreateOrder(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "order created successfully", order)
}

// GetMyOrders godoc
// GET /api/v1/orders
func GetMyOrders(c *gin.Context) {
	userID := c.GetString("user_id")

	var q models.OrderFilterQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := services.GetMyOrders(userID, q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", result)
}

// GetOrder godoc
// GET /api/v1/orders/:id
func GetOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	orderID := c.Param("id")

	if orderID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing order id")
		return
	}

	order, err := services.GetOrderByID(orderID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", order)
}

// CancelOrder godoc
// PATCH /api/v1/orders/:id/cancel
func CancelOrder(c *gin.Context) {
	userID := c.GetString("user_id")
	orderID := c.Param("id")

	if orderID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing order id")
		return
	}

	if err := services.CancelOrder(orderID, userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "order cancelled successfully", nil)
}

// GetAllOrders godoc
// GET /api/v1/admin/orders
func GetAllOrders(c *gin.Context) {
	var q models.OrderFilterQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := services.GetAllOrders(q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", result)
}

// UpdateOrderStatus godoc
// PATCH /api/v1/admin/orders/:id/status
func UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing order id")
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.UpdateOrderStatus(orderID, req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "order status updated successfully", nil)
}
