package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	AddressType string `bson:"address_type" json:"address_type"` // home, work, other
	Street      string `bson:"street" json:"street"`
	City        string `bson:"city" json:"city"`
	PostalCode  string `bson:"postal_code" json:"postal_code"`
	IsDefault   bool   `bson:"is_default" json:"is_default"`
}

type Customer struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name" binding:"required"`
	Email          string             `bson:"email" json:"email" binding:"required,email"`
	Password       string             `bson:"password" json:"-"`
	Phone          string             `bson:"phone" json:"phone"`
	Role           string             `bson:"role" json:"role"` // "admin" or "customer"
	Addresses      []Address          `bson:"addresses" json:"addresses"`
	LoyaltyPoints  int                `bson:"loyalty_points" json:"loyalty_points"`
	RegisteredDate time.Time          `bson:"registered_date" json:"registered_date"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type CustomerInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Phone    string `json:"phone"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AddressInput struct {
	AddressType string `json:"address_type" binding:"required"`
	Street      string `json:"street" binding:"required"`
	City        string `json:"city" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	IsDefault   bool   `json:"is_default"`
}

type CustomerResponse struct {
	ID             primitive.ObjectID `json:"id"`
	Name           string             `json:"name"`
	Email          string             `json:"email"`
	Phone          string             `json:"phone"`
	Role           string             `json:"role"`
	Addresses      []Address          `json:"addresses"`
	LoyaltyPoints  int                `json:"loyalty_points"`
	RegisteredDate time.Time          `json:"registered_date"`
}

func (c *Customer) ToResponse() CustomerResponse {
	return CustomerResponse{
		ID:             c.ID,
		Name:           c.Name,
		Email:          c.Email,
		Phone:          c.Phone,
		Role:           c.Role,
		Addresses:      c.Addresses,
		LoyaltyPoints:  c.LoyaltyPoints,
		RegisteredDate: c.RegisteredDate,
	}
}
