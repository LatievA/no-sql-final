package models

// API Response models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int64       `json:"total_pages"`
}

type AuthResponse struct {
	Success bool             `json:"success"`
	Token   string           `json:"token"`
	User    CustomerResponse `json:"user"`
}

// Report models
type SalesByCategory struct {
	CategoryID   interface{} `bson:"_id" json:"category_id"`
	CategoryName string      `bson:"category_name" json:"category_name"`
	TotalSales   float64     `bson:"total_sales" json:"total_sales"`
	TotalOrders  int         `bson:"total_orders" json:"total_orders"`
	TotalItems   int         `bson:"total_items" json:"total_items"`
}
