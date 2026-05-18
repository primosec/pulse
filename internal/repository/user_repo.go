package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/primosec/pulse/internal/domain"
	"github.com/primosec/pulse/internal/repository/db"
)

type UserRepository struct {
	q *db.Queries
}

func NewUserRepository(sqlDB *sql.DB) *UserRepository {
	return &UserRepository{q: db.New(sqlDB)}
}

func (r *UserRepository) Create(ctx context.Context, email, password, name string) (*domain.User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Email:    email,
		Password: password,
		Name:     name,
	})
	if err != nil {
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	u, err := r.q.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	return toDomainUser(u), nil
}

func toDomainUser(u db.User) *domain.User {
	return &domain.User{
		ID:        u.ID.String(),
		Email:     u.Email,
		Password:  u.Password,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}
}
