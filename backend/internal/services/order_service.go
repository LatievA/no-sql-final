package services

import (
	"bicycle-store/internal/models"
	"bicycle-store/internal/repositories"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService struct {
	orderRepo    *repositories.OrderRepository
	bicycleRepo  *repositories.BicycleRepository
	customerRepo *repositories.CustomerRepository
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo:    repositories.NewOrderRepository(),
		bicycleRepo:  repositories.NewBicycleRepository(),
		customerRepo: repositories.NewCustomerRepository(),
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, input models.OrderInput) (*models.Order, error) {
	custID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, errors.New("invalid customer ID")
	}

	// Get customer info
	customer, err := s.customerRepo.GetByID(ctx, custID)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	// Build order items
	var items []models.OrderItem
	var totalAmount float64

	for _, itemInput := range input.Items {
		bicycleID, err := primitive.ObjectIDFromHex(itemInput.BicycleID)
		if err != nil {
			return nil, errors.New("invalid bicycle ID")
		}

		bicycle, err := s.bicycleRepo.GetByID(ctx, bicycleID)
		if err != nil {
			return nil, errors.New("bicycle not found: " + itemInput.BicycleID)
		}

		// Check stock
		if bicycle.StockQuantity < itemInput.Quantity {
			return nil, errors.New("insufficient stock for: " + bicycle.ModelName)
		}

		item := models.OrderItem{
			BicycleID:              bicycleID,
			ModelName:              bicycle.ModelName,
			Brand:                  bicycle.Brand,
			Quantity:               itemInput.Quantity,
			PriceAtPurchase:        bicycle.Price,
			SelectedCustomizations: itemInput.SelectedCustomizations,
		}

		items = append(items, item)
		totalAmount += bicycle.Price * float64(itemInput.Quantity)
	}

	order := &models.Order{
		CustomerID:      custID,
		CustomerName:    customer.Name,
		Items:           items,
		TotalAmount:     totalAmount,
		DeliveryAddress: input.DeliveryAddress,
		PaymentMethod:   input.PaymentMethod,
	}

	// Use transaction to create order and decrement stock
	if err := s.orderRepo.CreateWithTransaction(ctx, order); err != nil {
		return nil, err
	}

	// Add loyalty points to customer (1 point per 1000 spent)
	points := int(totalAmount / 1000)
	if points > 0 {
		s.customerRepo.UpdateLoyaltyPoints(ctx, custID, points)
	}

	return order, nil
}

func (s *OrderService) GetOrders(ctx context.Context, filter models.OrderFilter) ([]models.Order, int64, error) {
	return s.orderRepo.GetAll(ctx, filter)
}

func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerID string, filter models.OrderFilter) ([]models.Order, int64, error) {
	custID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, 0, errors.New("invalid customer ID")
	}

	return s.orderRepo.GetByCustomerID(ctx, custID, filter)
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	id, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	return s.orderRepo.GetByID(ctx, id)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, status string) (*models.Order, error) {
	id, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, errors.New("invalid order ID")
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"shipped":   true,
		"delivered": true,
		"cancelled": true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid order status")
	}

	// If cancelling, use transaction to restore stock
	if status == "cancelled" {
		if err := s.orderRepo.CancelOrderWithTransaction(ctx, id); err != nil {
			return nil, err
		}
		return s.orderRepo.GetByID(ctx, id)
	}

	return s.orderRepo.UpdateStatus(ctx, id, status)
}

func (s *OrderService) AddReview(ctx context.Context, customerID, bicycleID string, input models.ReviewInput) (*models.Bicycle, error) {
	custID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		return nil, errors.New("invalid customer ID")
	}

	bicID, err := primitive.ObjectIDFromHex(bicycleID)
	if err != nil {
		return nil, errors.New("invalid bicycle ID")
	}

	// Get customer name
	customer, err := s.customerRepo.GetByID(ctx, custID)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	review := models.Review{
		ReviewID:     primitive.NewObjectID(),
		CustomerID:   custID,
		CustomerName: customer.Name,
		Rating:       input.Rating,
		Comment:      input.Comment,
		ReviewDate:   time.Now(),
	}

	return s.bicycleRepo.AddReview(ctx, bicID, review)
}
