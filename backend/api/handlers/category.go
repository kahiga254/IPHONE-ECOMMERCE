package handlers

import (
	"fmt"
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func GetAllCategories(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	
	var pageInt, limitInt int = 1, 10
	
	if p, err := parseInt(page); err == nil && p > 0 {
		pageInt = p
	}
	if l, err := parseInt(limit); err == nil && l > 0 && l <= 100 {
		limitInt = l
	}
	
	categories, total, err := services.GetAllCategories(pageInt, limitInt)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	totalPages := (total + limitInt - 1) / limitInt
	
	utils.SuccessResponse(c, http.StatusOK, "", gin.H{
		"categories":   categories,
		"total":        total,
		"page":         pageInt,
		"limit":        limitInt,
		"total_pages":  totalPages,
	})
}

func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	
	category, err := services.CreateCategory(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusCreated, "category created successfully", category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "category id required")
		return
	}
	
	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	
	err := services.UpdateCategory(id, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "category updated successfully", nil)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "category id required")
		return
	}
	
	err := services.DeleteCategory(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "category deleted successfully", nil)
}

func parseInt(s string) (int, error) {
	var result int
	_, err := fmt.Sscan(s, &result)
	return result, err
}
