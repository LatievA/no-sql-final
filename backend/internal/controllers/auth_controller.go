package controllers

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// Register godoc
// @Summary Register a new customer
// @Description Create a new customer account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.CustomerInput true "Customer registration data"
// @Success 201 {object} models.AuthResponse
// @Failure 400 {object} models.APIResponse
// @Router /auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var input models.CustomerInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response, err := c.authService.Register(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login to the system
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "Login credentials"
// @Success 200 {object} models.AuthResponse
// @Failure 401 {object} models.APIResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var input models.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response, err := c.authService.Login(ctx.Request.Context(), input)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetMe godoc
// @Summary Get current user
// @Description Get the currently authenticated user's profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse{data=models.CustomerResponse}
// @Failure 401 {object} models.APIResponse
// @Router /auth/me [get]
func (c *AuthController) GetMe(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	user, err := c.authService.GetCurrentUser(ctx.Request.Context(), userID.(string))
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    user,
	})
}
