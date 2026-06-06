package service

import (
	"context"
	"fmt"

	"github.com/water-factory/api/internal/domain"
	"github.com/water-factory/api/internal/dto"
	"github.com/water-factory/api/internal/repository"
	jwtpkg "github.com/water-factory/api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo *repository.UserRepo
	jwtMgr   *jwtpkg.Manager
}

func NewAuthService(userRepo *repository.UserRepo, jwtMgr *jwtpkg.Manager) *AuthService {
	return &AuthService{userRepo: userRepo, jwtMgr: jwtMgr}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokens, err := s.jwtMgr.Generate(user.ID, user.Email, user.RoleName)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User: dto.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.RoleName,
		},
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	claims, err := s.jwtMgr.Validate(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	tokens, err := s.jwtMgr.Generate(user.ID, user.Email, user.RoleName)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.ExpiresIn,
		User: dto.UserInfo{
			ID:       user.ID,
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.RoleName,
		},
	}, nil
}

func (s *AuthService) Me(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
