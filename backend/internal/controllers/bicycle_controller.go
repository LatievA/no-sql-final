package controllers

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/repositories"
	"bicycle-store/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BicycleController struct {
	repo         *repositories.BicycleRepository
	orderService *services.OrderService
}

func NewBicycleController() *BicycleController {
	return &BicycleController{
		repo:         repositories.NewBicycleRepository(),
		orderService: services.NewOrderService(),
	}
}

// GetAll godoc
// @Summary Get all bicycles
// @Description Get a paginated list of bicycles with optional filters
// @Tags bicycles
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param category_id query string false "Filter by category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param brand query string false "Filter by brand"
// @Param search query string false "Search in model name, brand, description"
// @Param sort query string false "Sort field" default(created_at)
// @Param order query string false "Sort order (asc/desc)" default(desc)
// @Success 200 {object} models.PaginatedResponse{data=[]models.Bicycle}
// @Router /bicycles [get]
func (c *BicycleController) GetAll(ctx *gin.Context) {
	var filter models.BicycleFilter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Set defaults
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}
	if filter.Sort == "" {
		filter.Sort = "created_at"
	}
	if filter.Order == "" {
		filter.Order = "desc"
	}

	bicycles, total, err := c.repo.GetAll(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch bicycles",
		})
		return
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

	ctx.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       bicycles,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get bicycle by ID
// @Description Get a single bicycle by its ID
// @Tags bicycles
// @Produce json
// @Param id path string true "Bicycle ID"
// @Success 200 {object} models.APIResponse{data=models.Bicycle}
// @Failure 404 {object} models.APIResponse
// @Router /bicycles/{id} [get]
func (c *BicycleController) GetByID(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid bicycle ID",
		})
		return
	}

	bicycle, err := c.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Bicycle not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    bicycle,
	})
}

// Create godoc
// @Summary Create a new bicycle
// @Description Create a new bicycle (Admin only)
// @Tags bicycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.BicycleInput true "Bicycle data"
// @Success 201 {object} models.APIResponse{data=models.Bicycle}
// @Failure 400 {object} models.APIResponse
// @Router /bicycles [post]
func (c *BicycleController) Create(ctx *gin.Context) {
	var input models.BicycleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	bicycle, err := c.repo.Create(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create bicycle: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Bicycle created successfully",
		Data:    bicycle,
	})
}

// Update godoc
// @Summary Update a bicycle
// @Description Update an existing bicycle (Admin only)
// @Tags bicycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bicycle ID"
// @Param input body models.BicycleInput true "Bicycle data"
// @Success 200 {object} models.APIResponse{data=models.Bicycle}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /bicycles/{id} [put]
func (c *BicycleController) Update(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid bicycle ID",
		})
		return
	}

	var input models.BicycleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	bicycle, err := c.repo.Update(ctx.Request.Context(), id, input)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Bicycle not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Bicycle updated successfully",
		Data:    bicycle,
	})
}

// Delete godoc
// @Summary Delete a bicycle
// @Description Delete a bicycle (Admin only)
// @Tags bicycles
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bicycle ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /bicycles/{id} [delete]
func (c *BicycleController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid bicycle ID",
		})
		return
	}

	if err := c.repo.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete bicycle",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Bicycle deleted successfully",
	})
}

// AddReview godoc
// @Summary Add a review to a bicycle
// @Description Add a customer review to a bicycle
// @Tags bicycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bicycle ID"
// @Param input body models.ReviewInput true "Review data"
// @Success 201 {object} models.APIResponse{data=models.Bicycle}
// @Failure 400 {object} models.APIResponse
// @Router /bicycles/{id}/reviews [post]
func (c *BicycleController) AddReview(ctx *gin.Context) {
	bicycleID := ctx.Param("id")
	customerID, _ := ctx.Get("userID")

	var input models.ReviewInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	bicycle, err := c.orderService.AddReview(ctx.Request.Context(), customerID.(string), bicycleID, input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Review added successfully",
		Data:    bicycle,
	})
}

// UpdateStock godoc
// @Summary Update bicycle stock
// @Description Update the stock quantity of a bicycle (Admin only)
// @Tags bicycles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Bicycle ID"
// @Param input body object{quantity=int} true "Stock quantity change (positive or negative)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /bicycles/{id}/stock [patch]
func (c *BicycleController) UpdateStock(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid bicycle ID",
		})
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	if err := c.repo.UpdateStock(ctx.Request.Context(), id, input.Quantity); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update stock",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Stock updated successfully",
	})
}
