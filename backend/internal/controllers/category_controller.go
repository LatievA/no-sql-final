package controllers

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryController struct {
	repo *repositories.CategoryRepository
}

func NewCategoryController() *CategoryController {
	return &CategoryController{
		repo: repositories.NewCategoryRepository(),
	}
}

// GetAll godoc
// @Summary Get all categories
// @Description Get a list of all bicycle categories
// @Tags categories
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Category}
// @Router /categories [get]
func (c *CategoryController) GetAll(ctx *gin.Context) {
	categories, err := c.repo.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to fetch categories",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    categories,
	})
}

// GetByID godoc
// @Summary Get category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.APIResponse{data=models.Category}
// @Failure 404 {object} models.APIResponse
// @Router /categories/{id} [get]
func (c *CategoryController) GetByID(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid category ID",
		})
		return
	}

	category, err := c.repo.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Category not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    category,
	})
}

// Create godoc
// @Summary Create a new category
// @Description Create a new bicycle category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.CategoryInput true "Category data"
// @Success 201 {object} models.APIResponse{data=models.Category}
// @Failure 400 {object} models.APIResponse
// @Router /categories [post]
func (c *CategoryController) Create(ctx *gin.Context) {
	var input models.CategoryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := c.repo.Create(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to create category",
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    category,
	})
}

// Update godoc
// @Summary Update a category
// @Description Update an existing category (Admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param input body models.CategoryInput true "Category data"
// @Success 200 {object} models.APIResponse{data=models.Category}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /categories/{id} [put]
func (c *CategoryController) Update(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid category ID",
		})
		return
	}

	var input models.CategoryInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := c.repo.Update(ctx.Request.Context(), id, input)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "Category not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category updated successfully",
		Data:    category,
	})
}

// Delete godoc
// @Summary Delete a category
// @Description Delete a category (Admin only)
// @Tags categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /categories/{id} [delete]
func (c *CategoryController) Delete(ctx *gin.Context) {
	id, err := primitive.ObjectIDFromHex(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "Invalid category ID",
		})
		return
	}

	if err := c.repo.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   "Failed to delete category",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}
