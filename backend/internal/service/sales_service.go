package service

import (
	"context"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/repository"
)

type SalesService struct {
	repo *repository.SalesRepo
}

func NewSalesService(repo *repository.SalesRepo) *SalesService {
	return &SalesService{repo: repo}
}

// Vehicles
func (s *SalesService) ListVehicles(ctx context.Context) ([]*domain.Vehicle, error) {
	return s.repo.ListVehicles(ctx)
}

// Sales Orders
func (s *SalesService) ListSalesOrders(ctx context.Context) ([]*domain.SalesOrder, error) {
	return s.repo.ListSalesOrders(ctx)
}
func (s *SalesService) GetSalesOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	return s.repo.GetSalesOrder(ctx, id)
}
func (s *SalesService) CreateSalesOrder(ctx context.Context, userID string, req *dto.CreateSalesOrderRequest) (*domain.SalesOrder, error) {
	return s.repo.CreateSalesOrder(ctx, userID, req)
}
func (s *SalesService) UpdateSalesOrder(ctx context.Context, id string, req *dto.UpdateSalesOrderRequest) (*domain.SalesOrder, error) {
	return s.repo.UpdateSalesOrder(ctx, id, req)
}
func (s *SalesService) ConfirmSalesOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	return s.repo.ConfirmSalesOrder(ctx, id)
}
func (s *SalesService) CancelSalesOrder(ctx context.Context, id string) (*domain.SalesOrder, error) {
	return s.repo.CancelSalesOrder(ctx, id)
}

// Delivery Orders
func (s *SalesService) ListDeliveryOrders(ctx context.Context) ([]*domain.DeliveryOrder, error) {
	return s.repo.ListDeliveryOrders(ctx)
}
func (s *SalesService) GetDeliveryOrder(ctx context.Context, id string) (*domain.DeliveryOrder, error) {
	return s.repo.GetDeliveryOrder(ctx, id)
}
func (s *SalesService) CreateDeliveryOrder(ctx context.Context, req *dto.CreateDeliveryOrderRequest) (*domain.DeliveryOrder, error) {
	return s.repo.CreateDeliveryOrder(ctx, req)
}
func (s *SalesService) DispatchDeliveryOrder(ctx context.Context, id string, userID string) (*domain.DeliveryOrder, error) {
	return s.repo.DispatchDeliveryOrder(ctx, id, userID)
}
func (s *SalesService) DeliverDeliveryOrder(ctx context.Context, id string) (*domain.DeliveryOrder, error) {
	return s.repo.DeliverDeliveryOrder(ctx, id)
}

// Invoices
func (s *SalesService) ListInvoices(ctx context.Context) ([]*domain.Invoice, error) {
	return s.repo.ListInvoices(ctx)
}
func (s *SalesService) GetInvoice(ctx context.Context, id string) (*domain.Invoice, error) {
	return s.repo.GetInvoice(ctx, id)
}
func (s *SalesService) CreateInvoice(ctx context.Context, userID string, req *dto.CreateInvoiceRequest) (*domain.Invoice, error) {
	return s.repo.CreateInvoice(ctx, userID, req)
}
func (s *SalesService) MarkInvoicePaid(ctx context.Context, id string, req *dto.MarkPaidRequest) (*domain.Invoice, error) {
	return s.repo.MarkInvoicePaid(ctx, id, req)
}
