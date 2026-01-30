package controllers

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerController struct {
	repo *repositories.CustomerRepository
}

func NewCustomerController() *CustomerController {
	return &CustomerController{
		repo: repositories.NewCustomerRepository(),
	}
}

// GetAll godoc
// @Summary Get all customers (Admin)
// @Description Get a list of all customers (Admin only)
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=[]models.CustomerResponse}
// @Router /customers [get]
func (c *CustomerController) GetAll(ctx *gin.Context) {
	customers, err := c.repo.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch customers",
		})
		return
	}

	// Convert to response format (hide passwords)
	var responses []models.CustomerResponse
	for _, customer := range customers {
		responses = append(responses, customer.ToResponse())
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    responses,
	})
}

// GetByID godoc
// @Summary Get customer by ID
// @Description Get a single customer by ID (Admin only)
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} models.APIResponse{data=models.CustomerResponse}
// @Failure 404 {object} models.APIResponse
// @Router /customers/{id} [get]
func (c *CustomerController) GetByID(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid customer ID",
		})
		return
	}

	customer, err := c.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Customer not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    customer.ToResponse(),
	})
}

// AddAddress godoc
// @Summary Add address to customer
// @Description Add a new address to the authenticated customer's profile
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.AddressInput true "Address data"
// @Success 200 {object} models.APIResponse{data=models.CustomerResponse}
// @Failure 400 {object} models.APIResponse
// @Router /customers/addresses [post]
func (c *CustomerController) AddAddress(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	var input models.AddressInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	address := models.Address{
		AddressType: input.AddressType,
		Street:      input.Street,
		City:        input.City,
		PostalCode:  input.PostalCode,
		IsDefault:   input.IsDefault,
	}

	customer, err := c.repo.AddAddress(ctx.Request.Context(), id, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to add address",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Address added successfully",
		Data:    customer.ToResponse(),
	})
}

// RemoveAddress godoc
// @Summary Remove address from customer
// @Description Remove an address from the authenticated customer's profile
// @Tags customers
// @Produce json
// @Security BearerAuth
// @Param type path string true "Address type (home, work, other)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /customers/addresses/{type} [delete]
func (c *CustomerController) RemoveAddress(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	addressType := ctx.Param("type")
	if addressType == "" {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Address type is required",
		})
		return
	}

	if err := c.repo.RemoveAddress(ctx.Request.Context(), id, addressType); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to remove address",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Address removed successfully",
	})
}

// UpdateProfile godoc
// @Summary Update customer profile
// @Description Update the authenticated customer's profile
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body object{name=string,phone=string} true "Profile data"
// @Success 200 {object} models.APIResponse{data=models.CustomerResponse}
// @Failure 400 {object} models.APIResponse
// @Router /customers/profile [put]
func (c *CustomerController) UpdateProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	id, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return
	}

	var input struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	updateData := make(map[string]interface{})
	if input.Name != "" {
		updateData["name"] = input.Name
	}
	if input.Phone != "" {
		updateData["phone"] = input.Phone
	}

	if len(updateData) == 0 {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "No data to update",
		})
		return
	}

	customer, err := c.repo.Update(ctx.Request.Context(), id, updateData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to update profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data:    customer.ToResponse(),
	})
}
