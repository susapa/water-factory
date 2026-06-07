package service

import (
	"context"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepo
}

func NewDashboardService(repo *repository.DashboardRepo) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetSummary(ctx context.Context) (*domain.DashboardSummary, error) {
	return s.repo.GetSummary(ctx)
}
