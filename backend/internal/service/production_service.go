package service

import (
	"context"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/repository"
)

type ProductionService struct {
	repo *repository.ProductionRepo
}

func NewProductionService(repo *repository.ProductionRepo) *ProductionService {
	return &ProductionService{repo: repo}
}

func (s *ProductionService) ListOrders(ctx context.Context) ([]*domain.ProductionOrder, error) {
	return s.repo.ListOrders(ctx)
}

func (s *ProductionService) GetOrder(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	return s.repo.GetOrder(ctx, id)
}

func (s *ProductionService) CreateOrder(ctx context.Context, userID string, req *dto.CreateProductionOrderRequest) (*domain.ProductionOrder, error) {
	return s.repo.CreateOrder(ctx, userID, req)
}

func (s *ProductionService) UpdateOrder(ctx context.Context, id string, req *dto.UpdateProductionOrderRequest) (*domain.ProductionOrder, error) {
	return s.repo.UpdateOrder(ctx, id, req)
}

func (s *ProductionService) ConfirmOrder(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	return s.repo.ConfirmOrder(ctx, id)
}

func (s *ProductionService) IssueRM(ctx context.Context, id, userID string, req *dto.IssueRMRequest) (*domain.ProductionOrder, error) {
	return s.repo.IssueRM(ctx, id, userID, req)
}

func (s *ProductionService) StartProduction(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	return s.repo.StartProduction(ctx, id)
}

func (s *ProductionService) RecordYield(ctx context.Context, id, userID string, req *dto.RecordYieldRequest) (*domain.ProductionYield, error) {
	return s.repo.RecordYield(ctx, id, userID, req)
}

func (s *ProductionService) CompleteOrder(ctx context.Context, id string) (*domain.ProductionOrder, error) {
	return s.repo.CompleteOrder(ctx, id)
}

func (s *ProductionService) CancelOrder(ctx context.Context, id string) error {
	return s.repo.CancelOrder(ctx, id)
}

func (s *ProductionService) ListYields(ctx context.Context, orderID string) ([]*domain.ProductionYield, error) {
	return s.repo.ListYields(ctx, orderID)
}
