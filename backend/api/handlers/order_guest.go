package handlers

import (
	"log"
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func CreateGuestOrder(c *gin.Context) {
	var req models.CreateGuestOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	order, err := services.CreateGuestOrder(req)
	if err != nil {
		log.Printf("❌ CreateGuestOrder error: %v", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "order created successfully", order)
}
