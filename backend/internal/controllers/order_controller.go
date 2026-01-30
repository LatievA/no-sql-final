package controllers

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService *services.OrderService
}

func NewOrderController() *OrderController {
	return &OrderController{
		orderService: services.NewOrderService(),
	}
}

// GetAll godoc
// @Summary Get all orders (Admin)
// @Description Get a paginated list of all orders (Admin only)
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Param customer_id query string false "Filter by customer ID"
// @Success 200 {object} models.PaginatedResponse{data=[]models.Order}
// @Router /orders [get]
func (c *OrderController) GetAll(ctx *gin.Context) {
	var filter models.OrderFilter
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

	orders, total, err := c.orderService.GetOrders(ctx.Request.Context(), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch orders",
		})
		return
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

	ctx.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       orders,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetMyOrders godoc
// @Summary Get current user's orders
// @Description Get orders for the authenticated customer
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} models.PaginatedResponse{data=[]models.Order}
// @Router /orders/my [get]
func (c *OrderController) GetMyOrders(ctx *gin.Context) {
	customerID, _ := ctx.Get("userID")

	var filter models.OrderFilter
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

	orders, total, err := c.orderService.GetOrdersByCustomer(ctx.Request.Context(), customerID.(string), filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch orders",
		})
		return
	}

	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

	ctx.JSON(http.StatusOK, models.PaginatedResponse{
		Success:    true,
		Data:       orders,
		Page:       filter.Page,
		Limit:      filter.Limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetByID godoc
// @Summary Get order by ID
// @Description Get a single order by its ID
// @Tags orders
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} models.APIResponse{data=models.Order}
// @Failure 404 {object} models.APIResponse
// @Router /orders/{id} [get]
func (c *OrderController) GetByID(ctx *gin.Context) {
	orderID := ctx.Param("id")

	order, err := c.orderService.GetOrderByID(ctx.Request.Context(), orderID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Order not found",
		})
		return
	}

	// Check if user is admin or order owner
	role, _ := ctx.Get("role")
	userID, _ := ctx.Get("userID")
	if role != "admin" && order.CustomerID.Hex() != userID.(string) {
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Error:   "Access denied",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    order,
	})
}

// Create godoc
// @Summary Create a new order
// @Description Create a new order for the authenticated customer
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.OrderInput true "Order data"
// @Success 201 {object} models.APIResponse{data=models.Order}
// @Failure 400 {object} models.APIResponse
// @Router /orders [post]
func (c *OrderController) Create(ctx *gin.Context) {
	customerID, _ := ctx.Get("userID")

	var input models.OrderInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	order, err := c.orderService.CreateOrder(ctx.Request.Context(), customerID.(string), input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Order created successfully",
		Data:    order,
	})
}

// UpdateStatus godoc
// @Summary Update order status
// @Description Update the status of an order (Admin only)
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param input body models.OrderStatusInput true "New status"
// @Success 200 {object} models.APIResponse{data=models.Order}
// @Failure 400 {object} models.APIResponse
// @Router /orders/{id}/status [patch]
func (c *OrderController) UpdateStatus(ctx *gin.Context) {
	orderID := ctx.Param("id")

	var input models.OrderStatusInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	order, err := c.orderService.UpdateOrderStatus(ctx.Request.Context(), orderID, input.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Order status updated successfully",
		Data:    order,
	})
}
