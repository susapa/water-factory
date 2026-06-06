package service

import (
	"context"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/repository"
)

type RMInventoryService struct {
	repo *repository.RMInventoryRepo
}

func NewRMInventoryService(repo *repository.RMInventoryRepo) *RMInventoryService {
	return &RMInventoryService{repo: repo}
}

func (s *RMInventoryService) ListWarehouseLocations(ctx context.Context) ([]*domain.WarehouseLocation, error) {
	return s.repo.ListWarehouseLocations(ctx)
}

func (s *RMInventoryService) ListGRNs(ctx context.Context) ([]*domain.GRN, error) {
	return s.repo.ListGRNs(ctx)
}

func (s *RMInventoryService) GetGRN(ctx context.Context, id string) (*domain.GRN, error) {
	return s.repo.GetGRN(ctx, id)
}

func (s *RMInventoryService) CreateGRN(ctx context.Context, userID string, req *dto.CreateGRNRequest) (*domain.GRN, error) {
	return s.repo.CreateGRN(ctx, userID, req)
}

func (s *RMInventoryService) UpdateGRN(ctx context.Context, id string, req *dto.UpdateGRNRequest) (*domain.GRN, error) {
	return s.repo.UpdateGRN(ctx, id, req)
}

func (s *RMInventoryService) ConfirmGRN(ctx context.Context, id, userID string) (*domain.GRN, error) {
	return s.repo.ConfirmGRN(ctx, id, userID)
}

func (s *RMInventoryService) CancelGRN(ctx context.Context, id string) error {
	return s.repo.CancelGRN(ctx, id)
}

func (s *RMInventoryService) ListStockSummary(ctx context.Context) ([]*domain.StockSummary, error) {
	return s.repo.ListStockSummary(ctx)
}

func (s *RMInventoryService) ListStockLots(ctx context.Context, rawMaterialID string) ([]*domain.StockLot, error) {
	return s.repo.ListStockLots(ctx, rawMaterialID)
}

func (s *RMInventoryService) GetStockLotDetail(ctx context.Context, id string) (*domain.StockLotDetail, error) {
	return s.repo.GetStockLotDetail(ctx, id)
}

func (s *RMInventoryService) ListAdjustments(ctx context.Context) ([]*domain.Adjustment, error) {
	return s.repo.ListAdjustments(ctx)
}

func (s *RMInventoryService) GetAdjustment(ctx context.Context, id string) (*domain.Adjustment, error) {
	return s.repo.GetAdjustment(ctx, id)
}

func (s *RMInventoryService) CreateAdjustment(ctx context.Context, userID string, req *dto.CreateAdjustmentRequest) (*domain.Adjustment, error) {
	return s.repo.CreateAdjustment(ctx, userID, req)
}

func (s *RMInventoryService) ApproveAdjustment(ctx context.Context, id, approverID string) (*domain.Adjustment, error) {
	return s.repo.ApproveAdjustment(ctx, id, approverID)
}

func (s *RMInventoryService) CancelAdjustment(ctx context.Context, id string) error {
	return s.repo.CancelAdjustment(ctx, id)
}
