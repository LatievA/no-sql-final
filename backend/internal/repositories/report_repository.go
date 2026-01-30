package repositories

import (
	"bicycle-store/internal/database"
	"bicycle-store/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportRepository struct{}

func NewReportRepository() *ReportRepository {
	return &ReportRepository{}
}

// GetSalesByCategory uses multi-stage aggregation pipeline
// Pipeline: $unwind items -> $lookup bicycles -> $lookup categories -> $group by category -> $sort
func (r *ReportRepository) GetSalesByCategory(ctx context.Context) ([]models.SalesByCategory, error) {
	collection := database.GetCollection("orders")

	pipeline := []bson.M{
		// Stage 1: Only include completed orders
		{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"delivered", "shipped", "confirmed"}},
			},
		},
		// Stage 2: Unwind the items array
		{
			"$unwind": "$items",
		},
		// Stage 3: Lookup bicycle details
		{
			"$lookup": bson.M{
				"from":         "bicycles",
				"localField":   "items.bicycle_id",
				"foreignField": "_id",
				"as":           "bicycle_info",
			},
		},
		// Stage 4: Unwind bicycle info
		{
			"$unwind": bson.M{
				"path":                       "$bicycle_info",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// Stage 5: Lookup category details
		{
			"$lookup": bson.M{
				"from":         "categories",
				"localField":   "bicycle_info.category_id",
				"foreignField": "_id",
				"as":           "category_info",
			},
		},
		// Stage 6: Unwind category info
		{
			"$unwind": bson.M{
				"path":                       "$category_info",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// Stage 7: Group by category
		{
			"$group": bson.M{
				"_id": "$category_info._id",
				"category_name": bson.M{
					"$first": "$category_info.category_name",
				},
				"total_sales": bson.M{
					"$sum": bson.M{
						"$multiply": []interface{}{"$items.price_at_purchase", "$items.quantity"},
					},
				},
				"total_orders": bson.M{
					"$addToSet": "$_id",
				},
				"total_items": bson.M{
					"$sum": "$items.quantity",
				},
			},
		},
		// Stage 8: Project to clean up the output
		{
			"$project": bson.M{
				"_id":           1,
				"category_name": 1,
				"total_sales":   1,
				"total_orders":  bson.M{"$size": "$total_orders"},
				"total_items":   1,
			},
		},
		// Stage 9: Sort by total sales descending
		{
			"$sort": bson.M{
				"total_sales": -1,
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []models.SalesByCategory
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetTopSellingBicycles returns top N best-selling bicycles
func (r *ReportRepository) GetTopSellingBicycles(ctx context.Context, limit int) ([]bson.M, error) {
	collection := database.GetCollection("orders")

	pipeline := []bson.M{
		// Stage 1: Match completed orders
		{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"delivered", "shipped", "confirmed"}},
			},
		},
		// Stage 2: Unwind items
		{
			"$unwind": "$items",
		},
		// Stage 3: Group by bicycle
		{
			"$group": bson.M{
				"_id":         "$items.bicycle_id",
				"model_name":  bson.M{"$first": "$items.model_name"},
				"brand":       bson.M{"$first": "$items.brand"},
				"total_sold":  bson.M{"$sum": "$items.quantity"},
				"total_sales": bson.M{"$sum": bson.M{"$multiply": []interface{}{"$items.price_at_purchase", "$items.quantity"}}},
			},
		},
		// Stage 4: Lookup current bicycle info
		{
			"$lookup": bson.M{
				"from":         "bicycles",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "bicycle_details",
			},
		},
		// Stage 5: Unwind bicycle details
		{
			"$unwind": bson.M{
				"path":                       "$bicycle_details",
				"preserveNullAndEmptyArrays": true,
			},
		},
		// Stage 6: Project final fields
		{
			"$project": bson.M{
				"_id":           1,
				"model_name":    1,
				"brand":         1,
				"total_sold":    1,
				"total_sales":   1,
				"current_price": "$bicycle_details.price",
				"current_stock": "$bicycle_details.stock_quantity",
			},
		},
		// Stage 7: Sort by total sold
		{
			"$sort": bson.M{"total_sold": -1},
		},
		// Stage 8: Limit results
		{
			"$limit": limit,
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetCustomerOrderStats returns order statistics for a specific customer
func (r *ReportRepository) GetCustomerOrderStats(ctx context.Context, customerID primitive.ObjectID) (bson.M, error) {
	collection := database.GetCollection("orders")

	pipeline := []bson.M{
		// Stage 1: Match customer's orders
		{
			"$match": bson.M{
				"customer_id": customerID,
			},
		},
		// Stage 2: Group to calculate stats
		{
			"$group": bson.M{
				"_id":              nil,
				"total_orders":     bson.M{"$sum": 1},
				"total_spent":      bson.M{"$sum": "$total_amount"},
				"average_order":    bson.M{"$avg": "$total_amount"},
				"completed_orders": bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$status", "delivered"}}, 1, 0}}},
				"pending_orders":   bson.M{"$sum": bson.M{"$cond": []interface{}{bson.M{"$eq": []interface{}{"$status", "pending"}}, 1, 0}}},
			},
		},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return bson.M{
			"total_orders":     0,
			"total_spent":      0,
			"average_order":    0,
			"completed_orders": 0,
			"pending_orders":   0,
		}, nil
	}

	return results[0], nil
}
