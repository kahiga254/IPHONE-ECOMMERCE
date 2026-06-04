package handlers

import (
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetAddresses godoc
// GET /api/v1/addresses
func GetAddresses(c *gin.Context) {
	userID := c.GetString("user_id")

	addresses, err := services.GetAddresses(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", addresses)
}

// CreateAddress godoc
// POST /api/v1/addresses
func CreateAddress(c *gin.Context) {
	userID := c.GetString("user_id")

	var req models.CreateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	address, err := services.CreateAddress(userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "address created successfully", address)
}

// UpdateAddress godoc
// PUT /api/v1/addresses/:id
func UpdateAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("id")

	if addressID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing address id")
		return
	}

	var req models.UpdateAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.UpdateAddress(addressID, userID, req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "address updated successfully", nil)
}

// DeleteAddress godoc
// DELETE /api/v1/addresses/:id
func DeleteAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("id")

	if addressID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing address id")
		return
	}

	if err := services.DeleteAddress(addressID, userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "address deleted successfully", nil)
}

// SetDefaultAddress godoc
// PATCH /api/v1/addresses/:id/default
func SetDefaultAddress(c *gin.Context) {
	userID := c.GetString("user_id")
	addressID := c.Param("id")

	if addressID == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing address id")
		return
	}

	if err := services.SetDefaultAddress(addressID, userID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "default address updated successfully", nil)
}
