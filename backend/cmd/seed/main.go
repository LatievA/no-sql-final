package main

import (
	"bicycle-store/internal/config"
	"bicycle-store/internal/database"
	"bicycle-store/internal/models"
	"bicycle-store/internal/utils"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	if err := database.Connect(cfg.MongoURI, cfg.DBName); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.Disconnect()

	ctx := context.Background()

	log.Println("Starting database seed...")

	// Clear existing data
	database.GetCollection("categories").Drop(ctx)
	database.GetCollection("bicycles").Drop(ctx)
	database.GetCollection("customers").Drop(ctx)
	database.GetCollection("orders").Drop(ctx)

	// Seed Categories
	categories := []models.Category{
		{
			ID:          primitive.NewObjectID(),
			Name:        "Mountain Bike",
			Description: "Off-road bicycles with suspension and durable frames for trail riding",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Road Bike",
			Description: "Lightweight bicycles designed for paved roads and racing",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "City Bike",
			Description: "Comfortable bicycles for urban commuting and casual rides",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "BMX",
			Description: "Compact bikes for tricks, racing, and freestyle riding",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          primitive.NewObjectID(),
			Name:        "Electric Bike",
			Description: "Motor-assisted bicycles for easier commuting and longer distances",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	categoryDocs := make([]interface{}, len(categories))
	for i, cat := range categories {
		categoryDocs[i] = cat
	}
	database.GetCollection("categories").InsertMany(ctx, categoryDocs)
	log.Printf("Inserted %d categories", len(categories))

	// Seed Bicycles
	bicycles := []models.Bicycle{
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "Trailblazer Pro 29",
			Brand:         "RockRider",
			Price:         450000,
			StockQuantity: 15,
			CategoryID:    categories[0].ID, // Mountain Bike
			Specifications: models.Specifications{
				FrameMaterial: "Aluminum",
				WheelSize:     "29 inches",
				GearCount:     21,
				BrakeType:     "Hydraulic Disc",
				Suspension:    "Front Suspension",
				Weight:        "13.5 kg",
				MaxLoad:       "120 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"Red", "Blue", "Black", "Green"}},
				{Name: "saddle_type", Options: []string{"Sport", "Comfort", "Racing"}},
				{Name: "accessories", Options: []string{"Water Bottle Holder", "Phone Mount", "LED Light"}},
			},
			Description: "Professional mountain bike for trail riding with excellent suspension",
			ImageURL:    "/images/trailblazer-pro-29.jpg",
			Reviews: []models.Review{
				{
					ReviewID:     primitive.NewObjectID(),
					CustomerName: "Aidos Bekzhanov",
					Rating:       5,
					Comment:      "Excellent bike for mountain trails! Great suspension.",
					ReviewDate:   time.Now().AddDate(0, -1, 0),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "SpeedMaster 700C",
			Brand:         "VeloRace",
			Price:         620000,
			StockQuantity: 8,
			CategoryID:    categories[1].ID, // Road Bike
			Specifications: models.Specifications{
				FrameMaterial: "Carbon Fiber",
				WheelSize:     "700C",
				GearCount:     22,
				BrakeType:     "Caliper Brakes",
				Suspension:    "None",
				Weight:        "8.2 kg",
				MaxLoad:       "100 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"White", "Red", "Black"}},
				{Name: "handlebar_tape", Options: []string{"Black", "White", "Blue", "Red"}},
				{Name: "pedal_type", Options: []string{"Flat", "Clipless SPD", "Clipless Look"}},
			},
			Description: "Lightweight carbon road bike for competitive racing",
			ImageURL:    "/images/speedmaster-700c.jpg",
			Reviews:     []models.Review{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "City Comfort 26",
			Brand:         "UrbanRide",
			Price:         180000,
			StockQuantity: 25,
			CategoryID:    categories[2].ID, // City Bike
			Specifications: models.Specifications{
				FrameMaterial: "Steel",
				WheelSize:     "26 inches",
				GearCount:     7,
				BrakeType:     "V-Brake",
				Suspension:    "None",
				Weight:        "15 kg",
				MaxLoad:       "110 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"Mint", "Cream", "Black", "Pink"}},
				{Name: "basket", Options: []string{"None", "Front Basket", "Rear Rack"}},
				{Name: "bell_type", Options: []string{"Classic", "Electronic", "None"}},
			},
			Description: "Comfortable city bike perfect for daily commuting",
			ImageURL:    "/images/city-comfort-26.jpg",
			Reviews: []models.Review{
				{
					ReviewID:     primitive.NewObjectID(),
					CustomerName: "Dana Sultanova",
					Rating:       4,
					Comment:      "Great quality, fast delivery. Perfect for my daily commute.",
					ReviewDate:   time.Now().AddDate(0, 0, -15),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "Freestyle X20",
			Brand:         "TrickStar",
			Price:         145000,
			StockQuantity: 12,
			CategoryID:    categories[3].ID, // BMX
			Specifications: models.Specifications{
				FrameMaterial: "Chromoly Steel",
				WheelSize:     "20 inches",
				GearCount:     1,
				BrakeType:     "U-Brake",
				Suspension:    "None",
				Weight:        "11 kg",
				MaxLoad:       "90 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"Matte Black", "Chrome", "Neon Green", "Orange"}},
				{Name: "peg_set", Options: []string{"None", "Front Pegs", "Rear Pegs", "Full Set"}},
				{Name: "grip_color", Options: []string{"Black", "White", "Red", "Blue"}},
			},
			Description: "Professional BMX bike for freestyle tricks and skatepark riding",
			ImageURL:    "/images/freestyle-x20.jpg",
			Reviews:     []models.Review{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "EcoRider E-500",
			Brand:         "PowerPedal",
			Price:         850000,
			StockQuantity: 6,
			CategoryID:    categories[4].ID, // Electric Bike
			Specifications: models.Specifications{
				FrameMaterial: "Aluminum",
				WheelSize:     "27.5 inches",
				GearCount:     9,
				BrakeType:     "Hydraulic Disc",
				Suspension:    "Front Suspension",
				Weight:        "22 kg",
				MaxLoad:       "130 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"Grey", "Black", "White"}},
				{Name: "battery_size", Options: []string{"Standard 400Wh", "Extended 600Wh"}},
				{Name: "display_type", Options: []string{"Basic LCD", "Color LCD", "Smartphone App"}},
			},
			Description: "Electric bike with 500W motor and 80km range",
			ImageURL:    "/images/ecorider-e500.jpg",
			Reviews: []models.Review{
				{
					ReviewID:     primitive.NewObjectID(),
					CustomerName: "Marat Kozhaev",
					Rating:       5,
					Comment:      "Amazing e-bike! The battery lasts forever and the motor is powerful.",
					ReviewDate:   time.Now().AddDate(0, 0, -7),
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "Summit Explorer 27.5",
			Brand:         "RockRider",
			Price:         380000,
			StockQuantity: 18,
			CategoryID:    categories[0].ID, // Mountain Bike
			Specifications: models.Specifications{
				FrameMaterial: "Aluminum",
				WheelSize:     "27.5 inches",
				GearCount:     18,
				BrakeType:     "Mechanical Disc",
				Suspension:    "Hardtail",
				Weight:        "14 kg",
				MaxLoad:       "115 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"Orange", "Blue", "Black"}},
				{Name: "saddle_type", Options: []string{"Standard", "Gel Comfort"}},
				{Name: "tire_type", Options: []string{"All-Terrain", "Mud", "Mixed"}},
			},
			Description: "Versatile mountain bike for beginners and intermediate riders",
			ImageURL:    "/images/summit-explorer-275.jpg",
			Reviews:     []models.Review{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "Aero Elite 700",
			Brand:         "VeloRace",
			Price:         920000,
			StockQuantity: 4,
			CategoryID:    categories[1].ID, // Road Bike
			Specifications: models.Specifications{
				FrameMaterial: "Carbon Fiber",
				WheelSize:     "700C",
				GearCount:     24,
				BrakeType:     "Disc Brakes",
				Suspension:    "None",
				Weight:        "7.5 kg",
				MaxLoad:       "95 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"Stealth Black", "Team Red", "Sky Blue"}},
				{Name: "wheelset", Options: []string{"Alloy", "Carbon", "Aero Carbon"}},
				{Name: "groupset", Options: []string{"Shimano 105", "Shimano Ultegra", "SRAM Rival"}},
			},
			Description: "Top-tier aerodynamic road bike for professional cyclists",
			ImageURL:    "/images/aero-elite-700.jpg",
			Reviews:     []models.Review{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:            primitive.NewObjectID(),
			ModelName:     "Urban Glide 28",
			Brand:         "UrbanRide",
			Price:         220000,
			StockQuantity: 20,
			CategoryID:    categories[2].ID, // City Bike
			Specifications: models.Specifications{
				FrameMaterial: "Aluminum",
				WheelSize:     "28 inches",
				GearCount:     8,
				BrakeType:     "Roller Brakes",
				Suspension:    "Seatpost Suspension",
				Weight:        "14 kg",
				MaxLoad:       "120 kg",
			},
			CustomizationOptions: []models.CustomizationOption{
				{Name: "frame_color", Options: []string{"Navy", "Silver", "Burgundy"}},
				{Name: "lighting", Options: []string{"None", "Front Light", "Full Set"}},
				{Name: "fenders", Options: []string{"None", "Plastic", "Metal"}},
			},
			Description: "Premium city bike with integrated lights and fenders",
			ImageURL:    "/images/urban-glide-28.jpg",
			Reviews:     []models.Review{},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	bicycleDocs := make([]interface{}, len(bicycles))
	for i, bike := range bicycles {
		bicycleDocs[i] = bike
	}
	database.GetCollection("bicycles").InsertMany(ctx, bicycleDocs)
	log.Printf("Inserted %d bicycles", len(bicycles))

	// Seed Customers
	adminPassword, _ := utils.HashPassword("admin123")
	customerPassword, _ := utils.HashPassword("password123")

	customers := []models.Customer{
		{
			ID:       primitive.NewObjectID(),
			Name:     "Admin User",
			Email:    "admin@store.com",
			Password: adminPassword,
			Phone:    "+77001234567",
			Role:     "admin",
			Addresses: []models.Address{
				{
					AddressType: "work",
					Street:      "Admin Street 1",
					City:        "Astana",
					PostalCode:  "010000",
					IsDefault:   true,
				},
			},
			LoyaltyPoints:  0,
			RegisteredDate: time.Now().AddDate(0, -6, 0),
			CreatedAt:      time.Now().AddDate(0, -6, 0),
			UpdatedAt:      time.Now(),
		},
		{
			ID:       primitive.NewObjectID(),
			Name:     "Aidos Bekzhanov",
			Email:    "customer@store.com",
			Password: customerPassword,
			Phone:    "+77012345678",
			Role:     "customer",
			Addresses: []models.Address{
				{
					AddressType: "home",
					Street:      "Mangilik El 55",
					City:        "Astana",
					PostalCode:  "010000",
					IsDefault:   true,
				},
				{
					AddressType: "work",
					Street:      "Kabanbay Batyr 53",
					City:        "Astana",
					PostalCode:  "010001",
					IsDefault:   false,
				},
			},
			LoyaltyPoints:  150,
			RegisteredDate: time.Now().AddDate(0, -3, 0),
			CreatedAt:      time.Now().AddDate(0, -3, 0),
			UpdatedAt:      time.Now(),
		},
	}

	customerDocs := make([]interface{}, len(customers))
	for i, cust := range customers {
		customerDocs[i] = cust
	}
	database.GetCollection("customers").InsertMany(ctx, customerDocs)
	log.Printf("Inserted %d customers", len(customers))

	// Seed Orders
	orders := []models.Order{
		{
			ID:           primitive.NewObjectID(),
			CustomerID:   customers[1].ID,
			CustomerName: customers[1].Name,
			OrderDate:    time.Now().AddDate(0, -1, 0),
			Status:       "delivered",
			Items: []models.OrderItem{
				{
					BicycleID:       bicycles[0].ID,
					ModelName:       bicycles[0].ModelName,
					Brand:           bicycles[0].Brand,
					Quantity:        1,
					PriceAtPurchase: bicycles[0].Price,
					SelectedCustomizations: []models.SelectedCustomization{
						{Name: "frame_color", Value: "Blue"},
						{Name: "saddle_type", Value: "Comfort"},
					},
				},
			},
			TotalAmount: bicycles[0].Price,
			DeliveryAddress: models.DeliveryAddress{
				Street:     "Mangilik El 55",
				City:       "Astana",
				PostalCode: "010000",
				Phone:      "+77012345678",
			},
			PaymentMethod: "card",
			PaymentStatus: "paid",
			CreatedAt:     time.Now().AddDate(0, -1, 0),
			UpdatedAt:     time.Now().AddDate(0, 0, -20),
		},
		{
			ID:           primitive.NewObjectID(),
			CustomerID:   customers[1].ID,
			CustomerName: customers[1].Name,
			OrderDate:    time.Now().AddDate(0, 0, -15),
			Status:       "shipped",
			Items: []models.OrderItem{
				{
					BicycleID:       bicycles[2].ID,
					ModelName:       bicycles[2].ModelName,
					Brand:           bicycles[2].Brand,
					Quantity:        1,
					PriceAtPurchase: bicycles[2].Price,
					SelectedCustomizations: []models.SelectedCustomization{
						{Name: "frame_color", Value: "Mint"},
						{Name: "basket", Value: "Front Basket"},
					},
				},
				{
					BicycleID:              bicycles[3].ID,
					ModelName:              bicycles[3].ModelName,
					Brand:                  bicycles[3].Brand,
					Quantity:               1,
					PriceAtPurchase:        bicycles[3].Price,
					SelectedCustomizations: []models.SelectedCustomization{},
				},
			},
			TotalAmount: bicycles[2].Price + bicycles[3].Price,
			DeliveryAddress: models.DeliveryAddress{
				Street:     "Kabanbay Batyr 53",
				City:       "Astana",
				PostalCode: "010001",
				Phone:      "+77012345678",
			},
			PaymentMethod: "card",
			PaymentStatus: "paid",
			CreatedAt:     time.Now().AddDate(0, 0, -15),
			UpdatedAt:     time.Now().AddDate(0, 0, -10),
		},
		{
			ID:           primitive.NewObjectID(),
			CustomerID:   customers[1].ID,
			CustomerName: customers[1].Name,
			OrderDate:    time.Now().AddDate(0, 0, -2),
			Status:       "pending",
			Items: []models.OrderItem{
				{
					BicycleID:       bicycles[4].ID,
					ModelName:       bicycles[4].ModelName,
					Brand:           bicycles[4].Brand,
					Quantity:        1,
					PriceAtPurchase: bicycles[4].Price,
					SelectedCustomizations: []models.SelectedCustomization{
						{Name: "frame_color", Value: "Black"},
						{Name: "battery_size", Value: "Extended 600Wh"},
						{Name: "display_type", Value: "Color LCD"},
					},
				},
			},
			TotalAmount: bicycles[4].Price,
			DeliveryAddress: models.DeliveryAddress{
				Street:     "Mangilik El 55",
				City:       "Astana",
				PostalCode: "010000",
				Phone:      "+77012345678",
			},
			PaymentMethod: "cash",
			PaymentStatus: "pending",
			CreatedAt:     time.Now().AddDate(0, 0, -2),
			UpdatedAt:     time.Now().AddDate(0, 0, -2),
		},
	}

	orderDocs := make([]interface{}, len(orders))
	for i, order := range orders {
		orderDocs[i] = order
	}
	database.GetCollection("orders").InsertMany(ctx, orderDocs)
	log.Printf("Inserted %d orders", len(orders))

	log.Println("Database seed completed successfully!")
	log.Println("")
	log.Println("Demo accounts:")
	log.Println("  Admin: admin@store.com / admin123")
	log.Println("  Customer: customer@store.com / password123")
}
