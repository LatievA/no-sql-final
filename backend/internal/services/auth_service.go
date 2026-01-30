package services

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/repositories"
	"bicycle-store/internal/utils"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	customerRepo *repositories.CustomerRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		customerRepo: repositories.NewCustomerRepository(),
	}
}

func (s *AuthService) Register(ctx context.Context, input models.CustomerInput) (*models.AuthResponse, error) {
	// Check if email already exists
	existingCustomer, err := s.customerRepo.GetByEmail(ctx, input.Email)
	if err == nil && existingCustomer != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Create customer
	customer := &models.Customer{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
		Role:     "customer",
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, err
	}

	// Generate token
	token, err := utils.GenerateToken(customer.ID, customer.Email, customer.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Success: true,
		Token:   token,
		User:    customer.ToResponse(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input models.LoginInput) (*models.AuthResponse, error) {
	customer, err := s.customerRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Check password
	if !utils.CheckPassword(input.Password, customer.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate token
	token, err := utils.GenerateToken(customer.ID, customer.Email, customer.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		Success: true,
		Token:   token,
		User:    customer.ToResponse(),
	}, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID string) (*models.CustomerResponse, error) {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	customer, err := s.customerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response := customer.ToResponse()
	return &response, nil
}
