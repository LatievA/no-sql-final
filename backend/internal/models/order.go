package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SelectedCustomization struct {
	Name  string `bson:"name" json:"name"`   // e.g., "frame_color"
	Value string `bson:"value" json:"value"` // e.g., "Red"
}

type OrderItem struct {
	BicycleID              primitive.ObjectID      `bson:"bicycle_id" json:"bicycle_id"`
	ModelName              string                  `bson:"model_name" json:"model_name"`
	Brand                  string                  `bson:"brand" json:"brand"`
	Quantity               int                     `bson:"quantity" json:"quantity"`
	PriceAtPurchase        float64                 `bson:"price_at_purchase" json:"price_at_purchase"`
	SelectedCustomizations []SelectedCustomization `bson:"selected_customizations" json:"selected_customizations"`
}

type DeliveryAddress struct {
	Street     string `bson:"street" json:"street"`
	City       string `bson:"city" json:"city"`
	PostalCode string `bson:"postal_code" json:"postal_code"`
	Phone      string `bson:"phone" json:"phone"`
}

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerID      primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	CustomerName    string             `bson:"customer_name" json:"customer_name"`
	OrderDate       time.Time          `bson:"order_date" json:"order_date"`
	Status          string             `bson:"status" json:"status"` // pending, confirmed, shipped, delivered, cancelled
	Items           []OrderItem        `bson:"items" json:"items"`
	TotalAmount     float64            `bson:"total_amount" json:"total_amount"`
	DeliveryAddress DeliveryAddress    `bson:"delivery_address" json:"delivery_address"`
	PaymentMethod   string             `bson:"payment_method" json:"payment_method"`
	PaymentStatus   string             `bson:"payment_status" json:"payment_status"` // pending, paid, refunded
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}

type OrderItemInput struct {
	BicycleID              string                  `json:"bicycle_id" binding:"required"`
	Quantity               int                     `json:"quantity" binding:"required,min=1"`
	SelectedCustomizations []SelectedCustomization `json:"selected_customizations"`
}

type OrderInput struct {
	Items           []OrderItemInput `json:"items" binding:"required,min=1"`
	DeliveryAddress DeliveryAddress  `json:"delivery_address" binding:"required"`
	PaymentMethod   string           `json:"payment_method" binding:"required"`
}

type OrderStatusInput struct {
	Status string `json:"status" binding:"required"`
}

type OrderFilter struct {
	Status     string `form:"status"`
	CustomerID string `form:"customer_id"`
	Page       int    `form:"page,default=1"`
	Limit      int    `form:"limit,default=10"`
}
