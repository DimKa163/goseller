package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/DimKa163/goseller/internal/dberror"
	"github.com/DimKa163/goseller/internal/user/domain"
	"go.uber.org/zap"
)

type CreateUserRequest struct {
	Name  string       `json:"name" binding:"required"`
	Email domain.Email `json:"email" binding:"required,email"`
	Phone domain.Phone `json:"phone" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string       `json:"name" binding:"required"`
	Email domain.Email `json:"email" binding:"required,email"`
	Phone domain.Phone `json:"phone" binding:"required"`
}

type User struct {
	userRepository domain.UserRepository
	logger         *zap.Logger
}

func NewUser(userRepository domain.UserRepository, logger *zap.Logger) *User {
	return &User{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (u *User) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	user, err := u.userRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, dberror.ErrNoRows) {
			return nil, fmt.Errorf("%w; id: %d", domain.ErrUserNotFound, id)
		}
		return nil, err
	}
	return user, nil
}

func (u *User) Create(ctx context.Context, req *CreateUserRequest) (domain.UserID, error) {
	var usID domain.UserID
	var err error
	user := domain.CreateNewUser(req.Name, req.Email, req.Phone)
	usID, err = u.userRepository.Create(ctx, user)
	if err != nil {
		return usID, err
	}

	return usID, nil
}

func (u *User) Update(ctx context.Context, id domain.UserID, req *UpdateUserRequest) (*domain.User, error) {
	var user *domain.User
	var err error
	user, err = u.userRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, dberror.ErrNoRows) {
			return nil, fmt.Errorf("%w; id: %d", domain.ErrUserNotFound, id)
		}
		return nil, err
	}
	user.Update(req.Name, req.Email, req.Phone)
	if user, err = u.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Delete(ctx context.Context, id domain.UserID) error {
	var user *domain.User
	var err error
	user, err = u.userRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, dberror.ErrNoRows) {
			return fmt.Errorf("%w; id: %d", domain.ErrUserNotFound, id)
		}
		return err
	}
	user.Deactivate()
	if err = u.userRepository.Delete(ctx, user.ID); err != nil {
		return err
	}
	return nil
}
