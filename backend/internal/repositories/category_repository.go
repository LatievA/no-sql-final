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

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	collection := database.GetCollection("categories")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Category, error) {
	collection := database.GetCollection("categories")

	var category models.Category
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) Create(ctx context.Context, input models.CategoryInput) (*models.Category, error) {
	collection := database.GetCollection("categories")

	category := models.Category{
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := collection.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	category.ID = result.InsertedID.(primitive.ObjectID)
	return &category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id primitive.ObjectID, input models.CategoryInput) (*models.Category, error) {
	collection := database.GetCollection("categories")

	update := bson.M{
		"$set": bson.M{
			"category_name": input.Name,
			"description":   input.Description,
			"updated_at":    time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var category models.Category
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := database.GetCollection("categories")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
