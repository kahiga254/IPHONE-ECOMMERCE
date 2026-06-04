package handlers

import (
	"net/http"

	"backend/api/models"
	"backend/api/services"
	"backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetAllProducts godoc
// GET /api/v1/products
func GetAllProducts(c *gin.Context) {
	var q models.ProductFilterQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := services.GetAllProducts(q)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", result)
}

// GetProduct godoc
// GET /api/v1/products/:slug
func GetProduct(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing product slug")
		return
	}

	product, err := services.GetProductBySlug(slug)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", product)
}

// CreateProduct godoc
// POST /api/v1/admin/products
func CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	product, err := services.CreateProduct(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "product created successfully", product)
}

// UpdateProduct godoc
// PUT /api/v1/admin/products/:id
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing product id")
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.UpdateProduct(id, req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "product updated successfully", nil)
}

// DeleteProduct godoc
// DELETE /api/v1/admin/products/:id
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing product id")
		return
	}

	if err := services.DeleteProduct(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "product deleted successfully", nil)
}

// GetAllCategories godoc
// GET /api/v1/categories
func GetAllCategories(c *gin.Context) {
	categories, err := services.GetAllCategories()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "", categories)
}

// CreateCategory godoc
// POST /api/v1/admin/categories
func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := services.CreateCategory(req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "category created successfully", category)
}

// UpdateCategory godoc
// PUT /api/v1/admin/categories/:id
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing category id")
		return
	}

	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := services.UpdateCategory(id, req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "category updated successfully", nil)
}

// DeleteCategory godoc
// DELETE /api/v1/admin/categories/:id
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "missing category id")
		return
	}

	if err := services.DeleteCategory(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "category deleted successfully", nil)
}
