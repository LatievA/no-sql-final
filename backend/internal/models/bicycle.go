package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Specifications struct {
	FrameMaterial string `bson:"frame_material" json:"frame_material"`
	WheelSize     string `bson:"wheel_size" json:"wheel_size"`
	GearCount     int    `bson:"gear_count" json:"gear_count"`
	BrakeType     string `bson:"brake_type" json:"brake_type"`
	Suspension    string `bson:"suspension" json:"suspension"`
	Weight        string `bson:"weight" json:"weight"`
	MaxLoad       string `bson:"max_load" json:"max_load"`
}

type CustomizationOption struct {
	Name    string   `bson:"name" json:"name"`       // e.g., "frame_color", "wheel_size", "accessories"
	Options []string `bson:"options" json:"options"` // e.g., ["Red", "Blue", "Black"]
}

type Review struct {
	ReviewID     primitive.ObjectID `bson:"review_id" json:"review_id"`
	CustomerID   primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	CustomerName string             `bson:"customer_name" json:"customer_name"`
	Rating       int                `bson:"rating" json:"rating"`
	Comment      string             `bson:"comment" json:"comment"`
	ReviewDate   time.Time          `bson:"review_date" json:"review_date"`
}

type Bicycle struct {
	ID                   primitive.ObjectID    `bson:"_id,omitempty" json:"id"`
	ModelName            string                `bson:"model_name" json:"model_name" binding:"required"`
	Brand                string                `bson:"brand" json:"brand" binding:"required"`
	Price                float64               `bson:"price" json:"price" binding:"required"`
	StockQuantity        int                   `bson:"stock_quantity" json:"stock_quantity"`
	CategoryID           primitive.ObjectID    `bson:"category_id" json:"category_id" binding:"required"`
	Specifications       Specifications        `bson:"specifications" json:"specifications"`
	CustomizationOptions []CustomizationOption `bson:"customization_options" json:"customization_options"`
	Description          string                `bson:"description" json:"description"`
	ImageURL             string                `bson:"image_url" json:"image_url"`
	Reviews              []Review              `bson:"reviews" json:"reviews"`
	CreatedAt            time.Time             `bson:"created_at" json:"created_at"`
	UpdatedAt            time.Time             `bson:"updated_at" json:"updated_at"`
}

type BicycleInput struct {
	ModelName            string                `json:"model_name" binding:"required"`
	Brand                string                `json:"brand" binding:"required"`
	Price                float64               `json:"price" binding:"required"`
	StockQuantity        int                   `json:"stock_quantity"`
	CategoryID           string                `json:"category_id" binding:"required"`
	Specifications       Specifications        `json:"specifications"`
	CustomizationOptions []CustomizationOption `json:"customization_options"`
	Description          string                `json:"description"`
	ImageURL             string                `json:"image_url"`
}

type ReviewInput struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}

type BicycleFilter struct {
	CategoryID string  `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	Brand      string  `form:"brand"`
	Search     string  `form:"search"`
	Page       int     `form:"page,default=1"`
	Limit      int     `form:"limit,default=10"`
	Sort       string  `form:"sort,default=created_at"`
	Order      string  `form:"order,default=desc"`
}
