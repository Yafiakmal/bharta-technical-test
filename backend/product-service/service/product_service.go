package service

import (
	"context"
	"product-service/model"
	"product-service/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, p *model.Product) (*model.Product, error) {
	return s.repo.Create(ctx, p)
}

func (s *ProductService) FindAll(ctx context.Context) ([]model.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *ProductService) FindByID(ctx context.Context, id bson.ObjectID) (*model.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *ProductService) Update(ctx context.Context, id bson.ObjectID, p *model.Product) error {
	return s.repo.Update(ctx, id, p)
}

func (s *ProductService) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

func (s *ProductService) UpdateStock(ctx context.Context, id bson.ObjectID, newStock int) error {
	return s.repo.UpdateStock(ctx, id, newStock)
}
