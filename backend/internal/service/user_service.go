package service

import (
	"context"
	"fmt"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.ListAll(ctx)
}

func (s *UserService) Create(ctx context.Context, email, password, fullName string, roleID int) (*domain.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	return s.userRepo.Create(ctx, email, string(hash), fullName, roleID)
}

func (s *UserService) SetActive(ctx context.Context, id string, active bool) error {
	return s.userRepo.UpdateActive(ctx, id, active)
}

func (s *UserService) ListRoles(ctx context.Context) ([]*domain.Role, error) {
	return s.userRepo.ListRoles(ctx)
}
