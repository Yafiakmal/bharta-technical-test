package repository

import (
	"context"
	"product-service/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{collection: db.Collection("products")}
}

func (r *ProductRepository) Create(ctx context.Context, p *model.Product) (*model.Product, error) {
	p.ID = bson.NewObjectID()
	_, err := r.collection.InsertOne(ctx, p)
	return p, err
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var products []model.Product
	err = cursor.All(ctx, &products)
	return products, err
}

func (r *ProductRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.Product, error) {
	var p model.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, id bson.ObjectID, p *model.Product) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": p})
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id bson.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *ProductRepository) UpdateStock(ctx context.Context, id bson.ObjectID, newStock int) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"stock": newStock}})
	return err
}
