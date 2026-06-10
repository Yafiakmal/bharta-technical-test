package main

import (
	"context"
	"log"
	"order-service/handler"
	"order-service/repository"
	"order-service/service"
	"os"

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

	db := client.Database("order_db")
	repo := repository.NewOrderRepository(db)
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, X-Admin-Key", // Tambahkan X-Admin-Key
	}))

	// Order routes
	app.Post("/orders", h.Create)
	app.Get("/orders", h.FindAll)
	app.Get("/orders/:id", h.GetByID)
	app.Put("/orders/:id", h.Update)
	app.Delete("/orders/:id", h.Delete)

	// Force delete routes (for testing/cleanup)
	app.Delete("/orders/force/:id", h.ForceDelete)
	app.Delete("/orders/admin/delete-all", h.DeleteAllOrders)

	log.Fatal(app.Listen(":8002"))
}
