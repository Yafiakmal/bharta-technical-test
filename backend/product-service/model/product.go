package model

import "go.mongodb.org/mongo-driver/v2/bson"

type Product struct {
	ID    bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string        `bson:"name" json:"name"`
	Price float64       `bson:"price" json:"price"`
	Stock int           `bson:"stock" json:"stock"`
}
