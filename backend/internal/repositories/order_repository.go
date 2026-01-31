package repositories

import (
	"bicycle-store/internal/database"
	"bicycle-store/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) GetAll(ctx context.Context, filter models.OrderFilter) ([]models.Order, int64, error) {
	collection := database.GetCollection("orders")

	query := bson.M{}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	if filter.CustomerID != "" {
		customerID, err := primitive.ObjectIDFromHex(filter.CustomerID)
		if err == nil {
			query["customer_id"] = customerID
		}
	}

	// Count total
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Pagination
	skip := (filter.Page - 1) * filter.Limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(filter.Limit)).
		SetSort(bson.D{{Key: "order_date", Value: -1}})

	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Order, error) {
	collection := database.GetCollection("orders")

	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetByCustomerID(ctx context.Context, customerID primitive.ObjectID, filter models.OrderFilter) ([]models.Order, int64, error) {
	collection := database.GetCollection("orders")

	query := bson.M{"customer_id": customerID}

	if filter.Status != "" {
		query["status"] = filter.Status
	}

	// Count total
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Pagination
	skip := (filter.Page - 1) * filter.Limit
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(filter.Limit)).
		SetSort(bson.D{{Key: "order_date", Value: -1}})

	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	collection := database.GetCollection("orders")

	order.OrderDate = time.Now()
	order.Status = "pending"
	order.PaymentStatus = "pending"
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	result, err := collection.InsertOne(ctx, order)
	if err != nil {
		return err
	}

	order.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// CreateWithTransaction creates an order and decrements bicycle stock atomically
func (r *OrderRepository) CreateWithTransaction(ctx context.Context, order *models.Order) error {
	session, err := database.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		ordersCollection := database.GetCollection("orders")
		bicyclesCollection := database.GetCollection("bicycles")

		// Set order timestamps
		order.OrderDate = time.Now()
		order.Status = "pending"
		order.PaymentStatus = "pending"
		order.CreatedAt = time.Now()
		order.UpdatedAt = time.Now()

		// Insert order
		result, err := ordersCollection.InsertOne(sessCtx, order)
		if err != nil {
			return nil, err
		}
		order.ID = result.InsertedID.(primitive.ObjectID)

		// Decrement stock for each item using $inc
		for _, item := range order.Items {
			updateResult, err := bicyclesCollection.UpdateOne(
				sessCtx,
				bson.M{
					"_id":            item.BicycleID,
					"stock_quantity": bson.M{"$gte": item.Quantity},
				},
				bson.M{
					"$inc": bson.M{"stock_quantity": -item.Quantity},
					"$set": bson.M{"updated_at": time.Now()},
				},
			)
			if err != nil {
				return nil, err
			}
			if updateResult.MatchedCount == 0 {
				return nil, mongo.ErrNoDocuments
			}
		}

		return nil, nil
	})

	return err
}

// UpdateStatus uses $set to update order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) (*models.Order, error) {
	collection := database.GetCollection("orders")

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var order models.Order
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// UpdatePaymentStatus uses $set to update payment status
func (r *OrderRepository) UpdatePaymentStatus(ctx context.Context, id primitive.ObjectID, paymentStatus string) error {
	collection := database.GetCollection("orders")

	update := bson.M{
		"$set": bson.M{
			"payment_status": paymentStatus,
			"updated_at":     time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// AddItemToOrder uses $push to add an item to an existing order
func (r *OrderRepository) AddItemToOrder(ctx context.Context, orderID primitive.ObjectID, item models.OrderItem) error {
	collection := database.GetCollection("orders")

	update := bson.M{
		"$push": bson.M{
			"items": item,
		},
		"$inc": bson.M{
			"total_amount": item.PriceAtPurchase * float64(item.Quantity),
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": orderID, "status": "pending"}, update)
	return err
}

// RemoveItemFromOrder uses $pull to remove an item from an order
func (r *OrderRepository) RemoveItemFromOrder(ctx context.Context, orderID, bicycleID primitive.ObjectID) error {
	collection := database.GetCollection("orders")

	// First get the item to calculate the amount to subtract
	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		return err
	}

	var amountToSubtract float64
	for _, item := range order.Items {
		if item.BicycleID == bicycleID {
			amountToSubtract = item.PriceAtPurchase * float64(item.Quantity)
			break
		}
	}

	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{"bicycle_id": bicycleID},
		},
		"$inc": bson.M{
			"total_amount": -amountToSubtract,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": orderID, "status": "pending"}, update)
	return err
}

// UpdateItemQuantity uses positional $ operator to update quantity of a specific item
func (r *OrderRepository) UpdateItemQuantity(ctx context.Context, orderID, bicycleID primitive.ObjectID, newQuantity int) error {
	collection := database.GetCollection("orders")

	// First get the current item to calculate price difference
	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		return err
	}

	var priceAtPurchase float64
	var oldQuantity int
	for _, item := range order.Items {
		if item.BicycleID == bicycleID {
			priceAtPurchase = item.PriceAtPurchase
			oldQuantity = item.Quantity
			break
		}
	}

	priceDifference := priceAtPurchase * float64(newQuantity-oldQuantity)

	update := bson.M{
		"$set": bson.M{
			"items.$.quantity": newQuantity,
			"updated_at":       time.Now(),
		},
		"$inc": bson.M{
			"total_amount": priceDifference,
		},
	}

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": orderID, "items.bicycle_id": bicycleID, "status": "pending"},
		update,
	)
	return err
}

func (r *OrderRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := database.GetCollection("orders")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// CancelOrderWithTransaction cancels an order and restores stock
func (r *OrderRepository) CancelOrderWithTransaction(ctx context.Context, orderID primitive.ObjectID) error {
	session, err := database.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		ordersCollection := database.GetCollection("orders")
		bicyclesCollection := database.GetCollection("bicycles")

		// Get the order
		var order models.Order
		err := ordersCollection.FindOne(sessCtx, bson.M{"_id": orderID}).Decode(&order)
		if err != nil {
			return nil, err
		}

		// Restore stock for each item
		for _, item := range order.Items {
			_, err := bicyclesCollection.UpdateOne(
				sessCtx,
				bson.M{"_id": item.BicycleID},
				bson.M{
					"$inc": bson.M{"stock_quantity": item.Quantity},
					"$set": bson.M{"updated_at": time.Now()},
				},
			)
			if err != nil {
				return nil, err
			}
		}

		// Update order status to cancelled
		_, err = ordersCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": orderID},
			bson.M{
				"$set": bson.M{
					"status":     "cancelled",
					"updated_at": time.Now(),
				},
			},
		)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})

	return err
}
