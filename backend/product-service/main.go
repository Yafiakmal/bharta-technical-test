package main

import (
	"context"
	"log"
	"os"
	"product-service/handler"
	"product-service/repository"
	"product-service/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	db := client.Database("product_db")
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Content-Type",
	}))

	app.Post("/products", h.Create)
	app.Get("/products", h.FindAll)
	app.Get("/products/:id", h.FindByID)
	app.Put("/products/:id", h.Update)
	app.Delete("/products/:id", h.Delete)
	app.Patch("/products/:id/stock", h.UpdateStock)

	log.Fatal(app.Listen(":8001"))
}
