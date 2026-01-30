package repositories

import (
	"bicycle-store/internal/database"
	"bicycle-store/internal/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BicycleRepository struct{}

func NewBicycleRepository() *BicycleRepository {
	return &BicycleRepository{}
}

func (r *BicycleRepository) GetAll(ctx context.Context, filter models.BicycleFilter) ([]models.Bicycle, int64, error) {
	collection := database.GetCollection("bicycles")

	// Build filter query
	query := bson.M{}

	if filter.CategoryID != "" {
		categoryID, err := primitive.ObjectIDFromHex(filter.CategoryID)
		if err == nil {
			query["category_id"] = categoryID
		}
	}

	if filter.MinPrice > 0 {
		query["price"] = bson.M{"$gte": filter.MinPrice}
	}

	if filter.MaxPrice > 0 {
		if _, exists := query["price"]; exists {
			query["price"].(bson.M)["$lte"] = filter.MaxPrice
		} else {
			query["price"] = bson.M{"$lte": filter.MaxPrice}
		}
	}

	if filter.Brand != "" {
		query["brand"] = bson.M{"$regex": filter.Brand, "$options": "i"}
	}

	if filter.Search != "" {
		query["$or"] = []bson.M{
			{"model_name": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"brand": bson.M{"$regex": filter.Search, "$options": "i"}},
			{"description": bson.M{"$regex": filter.Search, "$options": "i"}},
		}
	}

	// Count total documents
	total, err := collection.CountDocuments(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	// Build options for pagination and sorting
	skip := (filter.Page - 1) * filter.Limit
	sortOrder := -1
	if filter.Order == "asc" {
		sortOrder = 1
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(filter.Limit)).
		SetSort(bson.D{{Key: filter.Sort, Value: sortOrder}})

	cursor, err := collection.Find(ctx, query, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var bicycles []models.Bicycle
	if err := cursor.All(ctx, &bicycles); err != nil {
		return nil, 0, err
	}

	return bicycles, total, nil
}

func (r *BicycleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Bicycle, error) {
	collection := database.GetCollection("bicycles")

	var bicycle models.Bicycle
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&bicycle)
	if err != nil {
		return nil, err
	}

	return &bicycle, nil
}

func (r *BicycleRepository) Create(ctx context.Context, input models.BicycleInput) (*models.Bicycle, error) {
	collection := database.GetCollection("bicycles")

	categoryID, err := primitive.ObjectIDFromHex(input.CategoryID)
	if err != nil {
		return nil, err
	}

	bicycle := models.Bicycle{
		ModelName:            input.ModelName,
		Brand:                input.Brand,
		Price:                input.Price,
		StockQuantity:        input.StockQuantity,
		CategoryID:           categoryID,
		Specifications:       input.Specifications,
		CustomizationOptions: input.CustomizationOptions,
		Description:          input.Description,
		ImageURL:             input.ImageURL,
		Reviews:              []models.Review{},
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	result, err := collection.InsertOne(ctx, bicycle)
	if err != nil {
		return nil, err
	}

	bicycle.ID = result.InsertedID.(primitive.ObjectID)
	return &bicycle, nil
}

func (r *BicycleRepository) Update(ctx context.Context, id primitive.ObjectID, input models.BicycleInput) (*models.Bicycle, error) {
	collection := database.GetCollection("bicycles")

	categoryID, err := primitive.ObjectIDFromHex(input.CategoryID)
	if err != nil {
		return nil, err
	}

	// Using $set for updating specific fields
	update := bson.M{
		"$set": bson.M{
			"model_name":            input.ModelName,
			"brand":                 input.Brand,
			"price":                 input.Price,
			"stock_quantity":        input.StockQuantity,
			"category_id":           categoryID,
			"specifications":        input.Specifications,
			"customization_options": input.CustomizationOptions,
			"description":           input.Description,
			"image_url":             input.ImageURL,
			"updated_at":            time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var bicycle models.Bicycle
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&bicycle)
	if err != nil {
		return nil, err
	}

	return &bicycle, nil
}

func (r *BicycleRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := database.GetCollection("bicycles")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AddReview uses $push to add a review to the reviews array
func (r *BicycleRepository) AddReview(ctx context.Context, bicycleID primitive.ObjectID, review models.Review) (*models.Bicycle, error) {
	collection := database.GetCollection("bicycles")

	update := bson.M{
		"$push": bson.M{
			"reviews": review,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var bicycle models.Bicycle
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": bicycleID}, update, opts).Decode(&bicycle)
	if err != nil {
		return nil, err
	}

	return &bicycle, nil
}

// RemoveReview uses $pull to remove a review from the reviews array
func (r *BicycleRepository) RemoveReview(ctx context.Context, bicycleID, reviewID primitive.ObjectID) error {
	collection := database.GetCollection("bicycles")

	update := bson.M{
		"$pull": bson.M{
			"reviews": bson.M{"review_id": reviewID},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": bicycleID}, update)
	return err
}

// UpdateStock uses $inc to increment or decrement stock
func (r *BicycleRepository) UpdateStock(ctx context.Context, id primitive.ObjectID, quantity int) error {
	collection := database.GetCollection("bicycles")

	update := bson.M{
		"$inc": bson.M{
			"stock_quantity": quantity,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

// UpdateReview uses positional $ operator to update a specific review
func (r *BicycleRepository) UpdateReview(ctx context.Context, bicycleID, reviewID primitive.ObjectID, rating int, comment string) error {
	collection := database.GetCollection("bicycles")

	update := bson.M{
		"$set": bson.M{
			"reviews.$.rating":      rating,
			"reviews.$.comment":     comment,
			"reviews.$.review_date": time.Now(),
			"updated_at":            time.Now(),
		},
	}

	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": bicycleID, "reviews.review_id": reviewID},
		update,
	)
	return err
}
