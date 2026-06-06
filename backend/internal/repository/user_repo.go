package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/water-factory/api/internal/domain"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.email, u.password_hash, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE u.email = $1 AND u.is_active = true`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return u, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT u.id, u.email, u.password_hash, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 WHERE u.id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return u, nil
}

func (r *UserRepo) ListAll(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT u.id, u.email, u.password_hash, u.full_name, u.role_id, r.name, u.is_active, u.created_at, u.updated_at
		 FROM users u
		 JOIN roles r ON r.id = u.role_id
		 ORDER BY u.created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.RoleName, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash, fullName string, roleID int) (*domain.User, error) {
	u := &domain.User{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, full_name, role_id)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, email, password_hash, full_name, role_id, is_active, created_at, updated_at`,
		email, passwordHash, fullName, roleID,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.FullName, &u.RoleID, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (r *UserRepo) UpdateActive(ctx context.Context, id string, active bool) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET is_active=$1, updated_at=NOW() WHERE id=$2`, active, id)
	return err
}

func (r *UserRepo) ListRoles(ctx context.Context) ([]*domain.Role, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name, COALESCE(description,'') FROM roles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]*domain.Role, 0)
	for rows.Next() {
		ro := &domain.Role{}
		if err := rows.Scan(&ro.ID, &ro.Name, &ro.Description); err != nil {
			return nil, err
		}
		roles = append(roles, ro)
	}
	return roles, nil
}
