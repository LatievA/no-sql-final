package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

func Connect(uri, dbName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	Client = client
	Database = client.Database(dbName)

	log.Println("Connected to MongoDB successfully")

	// Create indexes
	if err := createIndexes(); err != nil {
		log.Printf("Warning: Failed to create indexes: %v", err)
	}

	return nil
}

func Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return Client.Disconnect(ctx)
}

func GetCollection(name string) *mongo.Collection {
	return Database.Collection(name)
}

func createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Customers collection - unique email index
	customersCollection := GetCollection("customers")
	_, err := customersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		log.Printf("Warning: Failed to create email index: %v", err)
	}

	// Bicycles collection - category_id index
	bicyclesCollection := GetCollection("bicycles")
	_, err = bicyclesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "category_id", Value: 1}},
	})
	if err != nil {
		log.Printf("Warning: Failed to create category_id index: %v", err)
	}

	// Bicycles collection - compound index for filtering
	_, err = bicyclesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "price", Value: 1},
			{Key: "brand", Value: 1},
		},
	})
	if err != nil {
		log.Printf("Warning: Failed to create price-brand index: %v", err)
	}

	// Bicycles - text index for search
	_, err = bicyclesCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "model_name", Value: "text"},
			{Key: "brand", Value: "text"},
			{Key: "description", Value: "text"},
		},
	})
	if err != nil {
		log.Printf("Warning: Failed to create text index: %v", err)
	}

	// Orders collection - compound index for customer and status
	ordersCollection := GetCollection("orders")
	_, err = ordersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "customer_id", Value: 1},
			{Key: "status", Value: 1},
		},
	})
	if err != nil {
		log.Printf("Warning: Failed to create customer_id-status index: %v", err)
	}

	// Orders - order_date index for sorting
	_, err = ordersCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "order_date", Value: -1}},
	})
	if err != nil {
		log.Printf("Warning: Failed to create order_date index: %v", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}
