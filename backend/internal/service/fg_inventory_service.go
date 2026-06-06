package service

import (
	"context"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/repository"
)

type FGInventoryService struct {
	repo *repository.FGInventoryRepo
}

func NewFGInventoryService(repo *repository.FGInventoryRepo) *FGInventoryService {
	return &FGInventoryService{repo: repo}
}

func (s *FGInventoryService) ListStockSummary(ctx context.Context) ([]*domain.FGStockSummary, error) {
	return s.repo.ListStockSummary(ctx)
}

func (s *FGInventoryService) ListStockLots(ctx context.Context, finishedGoodID string) ([]*domain.FGStockLot, error) {
	return s.repo.ListStockLots(ctx, finishedGoodID)
}

func (s *FGInventoryService) GetStockLotDetail(ctx context.Context, id string) (*domain.FGStockLotDetail, error) {
	return s.repo.GetStockLotDetail(ctx, id)
}

func (s *FGInventoryService) ListAdjustments(ctx context.Context) ([]*domain.FGAdjustment, error) {
	return s.repo.ListAdjustments(ctx)
}

func (s *FGInventoryService) GetAdjustment(ctx context.Context, id string) (*domain.FGAdjustment, error) {
	return s.repo.GetAdjustment(ctx, id)
}

func (s *FGInventoryService) CreateAdjustment(ctx context.Context, userID string, req *dto.CreateFGAdjustmentRequest) (*domain.FGAdjustment, error) {
	return s.repo.CreateAdjustment(ctx, userID, req)
}

func (s *FGInventoryService) ApproveAdjustment(ctx context.Context, id, approverID string) (*domain.FGAdjustment, error) {
	return s.repo.ApproveAdjustment(ctx, id, approverID)
}

func (s *FGInventoryService) CancelAdjustment(ctx context.Context, id string) error {
	return s.repo.CancelAdjustment(ctx, id)
}
