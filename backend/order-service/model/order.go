package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Order struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID    bson.ObjectID `bson:"product_id" json:"product_id"`
	ProductName  string        `bson:"product_name" json:"product_name"`
	ProductPrice float64       `bson:"product_price" json:"product_price"`
	Quantity     int           `bson:"quantity" json:"quantity"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at"`
}
