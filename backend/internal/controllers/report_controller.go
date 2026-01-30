package controllers

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	repo *repositories.ReportRepository
}

func NewReportController() *ReportController {
	return &ReportController{
		repo: repositories.NewReportRepository(),
	}
}

// GetSalesByCategory godoc
// @Summary Get sales by category
// @Description Get sales statistics grouped by category (Admin only)
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.SalesByCategory}
// @Router /reports/sales-by-category [get]
func (c *ReportController) GetSalesByCategory(ctx *gin.Context) {
	results, err := c.repo.GetSalesByCategory(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate report",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    results,
	})
}

// GetTopSellingBicycles godoc
// @Summary Get top selling bicycles
// @Description Get the top N best-selling bicycles (Admin only)
// @Tags reports
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Number of results" default(10)
// @Success 200 {object} models.APIResponse
// @Router /reports/top-selling [get]
func (c *ReportController) GetTopSellingBicycles(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	results, err := c.repo.GetTopSellingBicycles(ctx.Request.Context(), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to generate report",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    results,
	})
}
