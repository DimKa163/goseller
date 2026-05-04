package infrastructure

import (
	"context"
	"errors"

	"github.com/DimKa163/goseller/internal/dberror"
	"github.com/DimKa163/goseller/internal/user/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	createUserQuery      = `INSERT INTO users (name, email, phone) VALUES ($1, $2, $3) RETURNING id`
	updateUserQuery      = `UPDATE users SET name = $1, email = $2, phone = $3, updated_at = now() WHERE id = $4 RETURNING id, created_at, updated_at, name, email, phone, is_active`
	getUserByIDQuery     = `SELECT id, created_at, updated_at, name, email, phone, is_active FROM public.users WHERE id = $1 and is_active = true`
	getUserByEmailQuery  = `SELECT id, created_at, updated_at, name, email, phone, is_active FROM public.users WHERE email = $1 and is_active = true`
	getUserByPhoneQuery  = `SELECT id, created_at, updated_at, name, email, phone, is_active FROM public.users WHERE phone = $1 and is_active = true`
	deactivateUserQuery  = `UPDATE users SET is_active = false, updated_at = now() WHERE id = $1`
	getCountByEmailQuery = `SELECT COUNT(*) FROM public.users WHERE email = $1 and is_active = true`
	getCountByPhoneQuery = `SELECT COUNT(*) FROM public.users WHERE phone = $1 and is_active = true`
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (domain.UserID, error) {
	var id domain.UserID
	err := r.db.QueryRow(ctx, createUserQuery, user.Name, user.Email, user.Phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Phone, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberror.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email domain.Email) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, getUserByEmailQuery, email).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Phone, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberror.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone domain.Phone) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(ctx, getUserByPhoneQuery, phone).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Phone, &user.IsActive)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, dberror.ErrNoRows
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	if err := r.db.QueryRow(ctx, updateUserQuery, user.Name, user.Email, user.Phone, user.ID).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Name, &user.Email, &user.Phone, &user.IsActive); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id domain.UserID) error {
	_, err := r.db.Exec(ctx, deactivateUserQuery, id)
	return err
}

func (r *UserRepository) GetCountByEmail(ctx context.Context, email domain.Email) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, getCountByEmailQuery, email).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) GetCountByPhone(ctx context.Context, phone domain.Phone) (int64, error) {
	var count int64
	err := r.db.QueryRow(ctx, getCountByPhoneQuery, phone).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
