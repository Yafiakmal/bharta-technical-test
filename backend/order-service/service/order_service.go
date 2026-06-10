package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"order-service/model"
	"order-service/repository"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func (s *OrderService) getProduct(productID string) (*ProductResponse, error) {
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8001"
	}

	resp, err := http.Get(fmt.Sprintf("%s/products/%s", productServiceURL, productID))
	if err != nil {
		return nil, fmt.Errorf("failed to reach product service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("product not found")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("product service error: status %d", resp.StatusCode)
	}

	var product ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *OrderService) updateStock(productID string, newStock int) error {
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8001"
	}

	body := fmt.Sprintf(`{"stock": %d}`, newStock)
	req, _ := http.NewRequest(http.MethodPatch,
		fmt.Sprintf("%s/products/%s/stock", productServiceURL, productID),
		strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to update stock, status: %d", resp.StatusCode)
	}
	return nil
}

func (s *OrderService) CreateOrder(ctx context.Context, productID string, quantity int) (*model.Order, error) {
	product, err := s.getProduct(productID)
	if err != nil {
		return nil, err
	}

	if product.Stock < quantity {
		return nil, fmt.Errorf("insufficient stock")
	}

	if err := s.updateStock(productID, product.Stock-quantity); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	productObjID, _ := bson.ObjectIDFromHex(productID)
	order := &model.Order{
		ProductID:    productObjID,
		ProductName:  product.Name,  // Simpan nama product
		ProductPrice: product.Price, // Simpan harga product
		Quantity:     quantity,
		CreatedAt:    time.Now(),
	}

	return s.repo.Create(ctx, order)
}

func (s *OrderService) FindAll(ctx context.Context) ([]model.Order, error) {
	return s.repo.FindAll(ctx)
}

func (s *OrderService) FindByID(ctx context.Context, id string) (*model.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}
	return order, nil
}

func (s *OrderService) UpdateOrder(ctx context.Context, id string, newQuantity int) (*model.Order, error) {
	existingOrder, err := s.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	productID := existingOrder.ProductID.Hex()
	product, err := s.getProduct(productID)
	if err != nil {
		if err.Error() == "product not found" {
			return nil, fmt.Errorf("cannot update order: product no longer exists")
		}
		return nil, err
	}

	stockDiff := newQuantity - existingOrder.Quantity

	if stockDiff > 0 && product.Stock < stockDiff {
		return nil, fmt.Errorf("insufficient stock")
	}

	newStock := product.Stock - stockDiff
	if err := s.updateStock(productID, newStock); err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	return s.repo.Update(ctx, id, newQuantity)
}

func (s *OrderService) DeleteOrder(ctx context.Context, id string) error {
	existingOrder, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}

	productID := existingOrder.ProductID.Hex()
	product, err := s.getProduct(productID)

	if err != nil {
		if err.Error() == "product not found" {
			log.Printf("Warning: Product %s not found for order %s, skipping stock restoration\n", productID, id)
		} else {
			log.Printf("Warning: Failed to get product %s for order %s: %v\n", productID, id, err)
		}
	} else {
		newStock := product.Stock + existingOrder.Quantity
		if err := s.updateStock(productID, newStock); err != nil {
			log.Printf("Warning: Failed to restore stock for product %s: %v\n", productID, err)
		}
	}

	return s.repo.Delete(ctx, id)
}

// ForceDeleteOrder - Delete order without any stock restoration or checks (for testing)
func (s *OrderService) ForceDeleteOrder(ctx context.Context, id string) error {
	// Optional: check if order exists
	_, err := s.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Directly delete without any stock operations
	return s.repo.ForceDelete(ctx, id)
}

// DeleteAllOrders - Delete all orders (for testing cleanup)
func (s *OrderService) DeleteAllOrders(ctx context.Context) (int64, error) {
	return s.repo.DeleteAll(ctx)
}
