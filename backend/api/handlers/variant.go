package handlers

import (
	"net/http"

	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UpdateVariantRequest struct {
	SKU     string   `json:"sku"`
	Color   string   `json:"color"`
	Storage string   `json:"storage"`
	Price   float64  `json:"price"`
	Stock   int      `json:"stock"`
	Images  []string `json:"images"`
}

func UpdateVariant(c *gin.Context) {
	variantID := c.Param("id")
	if variantID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "variant id required")
		return
	}

	var req UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := services.UpdateVariant(variantID, req.SKU, req.Color, req.Storage, req.Price, req.Stock, req.Images)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "variant updated successfully", nil)
}
