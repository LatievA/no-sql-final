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

type CustomerRepository struct{}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{}
}

func (r *CustomerRepository) GetAll(ctx context.Context) ([]models.Customer, error) {
	collection := database.GetCollection("customers")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var customers []models.Customer
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Customer, error) {
	collection := database.GetCollection("customers")

	var customer models.Customer
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	collection := database.GetCollection("customers")

	var customer models.Customer
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	collection := database.GetCollection("customers")

	customer.RegisteredDate = time.Now()
	customer.CreatedAt = time.Now()
	customer.UpdatedAt = time.Now()
	customer.LoyaltyPoints = 0
	customer.Addresses = []models.Address{}

	if customer.Role == "" {
		customer.Role = "customer"
	}

	result, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return err
	}

	customer.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *CustomerRepository) Update(ctx context.Context, id primitive.ObjectID, input map[string]interface{}) (*models.Customer, error) {
	collection := database.GetCollection("customers")

	input["updated_at"] = time.Now()
	update := bson.M{"$set": input}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var customer models.Customer
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update, opts).Decode(&customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := database.GetCollection("customers")

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// AddAddress uses $push to add an address to the addresses array
func (r *CustomerRepository) AddAddress(ctx context.Context, customerID primitive.ObjectID, address models.Address) (*models.Customer, error) {
	collection := database.GetCollection("customers")

	// If new address is default, unset other defaults first
	if address.IsDefault {
		_, err := collection.UpdateOne(
			ctx,
			bson.M{"_id": customerID},
			bson.M{"$set": bson.M{"addresses.$[].is_default": false}},
		)
		if err != nil {
			return nil, err
		}
	}

	update := bson.M{
		"$push": bson.M{
			"addresses": address,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var customer models.Customer
	err := collection.FindOneAndUpdate(ctx, bson.M{"_id": customerID}, update, opts).Decode(&customer)
	if err != nil {
		return nil, err
	}

	return &customer, nil
}

// RemoveAddress uses $pull to remove an address by type
func (r *CustomerRepository) RemoveAddress(ctx context.Context, customerID primitive.ObjectID, addressType string) error {
	collection := database.GetCollection("customers")

	update := bson.M{
		"$pull": bson.M{
			"addresses": bson.M{"address_type": addressType},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": customerID}, update)
	return err
}

// UpdateLoyaltyPoints uses $inc to increment loyalty points
func (r *CustomerRepository) UpdateLoyaltyPoints(ctx context.Context, customerID primitive.ObjectID, points int) error {
	collection := database.GetCollection("customers")

	update := bson.M{
		"$inc": bson.M{
			"loyalty_points": points,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(ctx, bson.M{"_id": customerID}, update)
	return err
}

// SetDefaultAddress uses positional $ operator to update a specific address
func (r *CustomerRepository) SetDefaultAddress(ctx context.Context, customerID primitive.ObjectID, addressType string) error {
	collection := database.GetCollection("customers")

	// First, unset all defaults
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": customerID},
		bson.M{"$set": bson.M{"addresses.$[].is_default": false}},
	)
	if err != nil {
		return err
	}

	// Then set the specific address as default using positional $ operator
	_, err = collection.UpdateOne(
		ctx,
		bson.M{"_id": customerID, "addresses.address_type": addressType},
		bson.M{"$set": bson.M{"addresses.$.is_default": true, "updated_at": time.Now()}},
	)
	return err
}
