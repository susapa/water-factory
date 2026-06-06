package service

import (
	"context"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/repository"
)

type MasterDataService struct {
	repo *repository.MasterDataRepo
}

func NewMasterDataService(repo *repository.MasterDataRepo) *MasterDataService {
	return &MasterDataService{repo: repo}
}

func (s *MasterDataService) ListUOM(ctx context.Context) ([]*domain.UnitOfMeasure, error) {
	return s.repo.ListUOM(ctx)
}

func (s *MasterDataService) ListCategories(ctx context.Context) ([]*domain.RawMaterialCategory, error) {
	return s.repo.ListCategories(ctx)
}

// Raw Materials

func (s *MasterDataService) ListRawMaterials(ctx context.Context) ([]*domain.RawMaterial, error) {
	return s.repo.ListRawMaterials(ctx)
}

func (s *MasterDataService) GetRawMaterial(ctx context.Context, id string) (*domain.RawMaterial, error) {
	return s.repo.GetRawMaterial(ctx, id)
}

func (s *MasterDataService) CreateRawMaterial(ctx context.Context, req *dto.CreateRawMaterialRequest) (*domain.RawMaterial, error) {
	return s.repo.CreateRawMaterial(ctx, req)
}

func (s *MasterDataService) UpdateRawMaterial(ctx context.Context, id string, req *dto.UpdateRawMaterialRequest) (*domain.RawMaterial, error) {
	return s.repo.UpdateRawMaterial(ctx, id, req)
}

func (s *MasterDataService) SetRawMaterialActive(ctx context.Context, id string, active bool) error {
	return s.repo.SetRawMaterialActive(ctx, id, active)
}

// Finished Goods

func (s *MasterDataService) ListFinishedGoods(ctx context.Context) ([]*domain.FinishedGood, error) {
	return s.repo.ListFinishedGoods(ctx)
}

func (s *MasterDataService) GetFinishedGood(ctx context.Context, id string) (*domain.FinishedGood, error) {
	return s.repo.GetFinishedGood(ctx, id)
}

func (s *MasterDataService) CreateFinishedGood(ctx context.Context, req *dto.CreateFinishedGoodRequest) (*domain.FinishedGood, error) {
	return s.repo.CreateFinishedGood(ctx, req)
}

func (s *MasterDataService) UpdateFinishedGood(ctx context.Context, id string, req *dto.UpdateFinishedGoodRequest) (*domain.FinishedGood, error) {
	return s.repo.UpdateFinishedGood(ctx, id, req)
}

func (s *MasterDataService) SetFinishedGoodActive(ctx context.Context, id string, active bool) error {
	return s.repo.SetFinishedGoodActive(ctx, id, active)
}

// BOM

func (s *MasterDataService) ListBOMs(ctx context.Context) ([]*domain.BOM, error) {
	return s.repo.ListBOMs(ctx)
}

func (s *MasterDataService) GetBOM(ctx context.Context, id string) (*domain.BOM, error) {
	return s.repo.GetBOM(ctx, id)
}

func (s *MasterDataService) CreateBOM(ctx context.Context, req *dto.CreateBOMRequest) (*domain.BOM, error) {
	return s.repo.CreateBOM(ctx, req)
}

func (s *MasterDataService) UpdateBOM(ctx context.Context, id string, req *dto.UpdateBOMRequest) (*domain.BOM, error) {
	return s.repo.UpdateBOM(ctx, id, req)
}

// Customers

func (s *MasterDataService) ListCustomers(ctx context.Context) ([]*domain.Customer, error) {
	return s.repo.ListCustomers(ctx)
}

func (s *MasterDataService) GetCustomer(ctx context.Context, id string) (*domain.Customer, error) {
	return s.repo.GetCustomer(ctx, id)
}

func (s *MasterDataService) CreateCustomer(ctx context.Context, req *dto.CreateCustomerRequest) (*domain.Customer, error) {
	return s.repo.CreateCustomer(ctx, req)
}

func (s *MasterDataService) UpdateCustomer(ctx context.Context, id string, req *dto.UpdateCustomerRequest) (*domain.Customer, error) {
	return s.repo.UpdateCustomer(ctx, id, req)
}

func (s *MasterDataService) SetCustomerActive(ctx context.Context, id string, active bool) error {
	return s.repo.SetCustomerActive(ctx, id, active)
}

// Suppliers

func (s *MasterDataService) ListSuppliers(ctx context.Context) ([]*domain.Supplier, error) {
	return s.repo.ListSuppliers(ctx)
}

func (s *MasterDataService) GetSupplier(ctx context.Context, id string) (*domain.Supplier, error) {
	return s.repo.GetSupplier(ctx, id)
}

func (s *MasterDataService) CreateSupplier(ctx context.Context, req *dto.CreateSupplierRequest) (*domain.Supplier, error) {
	return s.repo.CreateSupplier(ctx, req)
}

func (s *MasterDataService) UpdateSupplier(ctx context.Context, id string, req *dto.UpdateSupplierRequest) (*domain.Supplier, error) {
	return s.repo.UpdateSupplier(ctx, id, req)
}

func (s *MasterDataService) SetSupplierActive(ctx context.Context, id string, active bool) error {
	return s.repo.SetSupplierActive(ctx, id, active)
}
