package handler

import (
	"context"
	"order-service/service"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var body struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if body.ProductID == "" || body.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "productId and quantity required"})
	}

	order, err := h.service.CreateOrder(context.Background(), body.ProductID, body.Quantity)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		if err.Error() == "insufficient stock" {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(order)
}

func (h *OrderHandler) FindAll(c *fiber.Ctx) error {
	orders, err := h.service.FindAll(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}

func (h *OrderHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	order, err := h.service.FindByID(context.Background(), id)
	if err != nil {
		if err.Error() == "order not found" {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(order)
}

func (h *OrderHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}
	if body.Quantity <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "quantity must be greater than 0"})
	}

	order, err := h.service.UpdateOrder(context.Background(), id, body.Quantity)
	if err != nil {
		if err.Error() == "order not found" {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		if err.Error() == "insufficient stock" {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(order)
}

func (h *OrderHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteOrder(context.Background(), id); err != nil {
		if err.Error() == "order not found" {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "order deleted successfully"})
}

// ForceDelete - Delete order without any checks or stock restoration (for testing/cleanup)
func (h *OrderHandler) ForceDelete(c *fiber.Ctx) error {
	id := c.Params("id")

	// Optional: add secret key or admin token for security
	secretKey := c.Get("X-Admin-Key")
	if secretKey == "" || secretKey != "your-secret-key-here" {
		// You can remove this check if you don't need authentication for force delete
		// For testing, you can comment this block
		return c.Status(403).JSON(fiber.Map{"error": "unauthorized: admin key required"})
	}

	if err := h.service.ForceDeleteOrder(context.Background(), id); err != nil {
		if err.Error() == "order not found" {
			return c.Status(404).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "order force deleted successfully"})
}

// DeleteAllOrders - Delete all orders (for testing cleanup)
func (h *OrderHandler) DeleteAllOrders(c *fiber.Ctx) error {
	// Optional: add secret key or admin token for security
	secretKey := c.Get("X-Admin-Key")
	if secretKey == "" || secretKey != "your-secret-key-here" {
		return c.Status(403).JSON(fiber.Map{"error": "unauthorized: admin key required"})
	}

	count, err := h.service.DeleteAllOrders(context.Background())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "all orders deleted successfully",
		"deleted": count,
	})
}
