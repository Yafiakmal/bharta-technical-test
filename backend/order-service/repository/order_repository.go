package repository

import (
	"context"
	"order-service/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db *mongo.Database) *OrderRepository {
	return &OrderRepository{collection: db.Collection("orders")}
}

func (r *OrderRepository) Create(ctx context.Context, o *model.Order) (*model.Order, error) {
	o.ID = bson.NewObjectID()
	_, err := r.collection.InsertOne(ctx, o)
	return o, err
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]model.Order, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var orders []model.Order
	err = cursor.All(ctx, &orders)
	return orders, err
}

func (r *OrderRepository) FindByID(ctx context.Context, id string) (*model.Order, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order model.Order
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) Update(ctx context.Context, id string, quantity int) (*model.Order, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"quantity": quantity}},
	)
	if err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return nil // Order not found
	}
	return nil
}

// ForceDelete - Delete order directly without any checks
func (r *OrderRepository) ForceDelete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return nil // Order not found
	}
	return nil
}

// DeleteAll - Delete all orders (for testing cleanup)
func (r *OrderRepository) DeleteAll(ctx context.Context) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
