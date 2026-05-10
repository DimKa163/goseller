package usecase

import (
	"context"
	"errors"

	"github.com/DimKa163/goseller/internal/dberror"
	"github.com/DimKa163/goseller/internal/shared"
	"github.com/DimKa163/goseller/internal/shared/sellerlog"
	"github.com/DimKa163/goseller/internal/user/domain"
	"go.uber.org/zap"
)

type CreateUserRequest struct {
	Name  string       `json:"name" binding:"required"`
	Email domain.Email `json:"email" binding:"required,seller_email"`
	Phone domain.Phone `json:"phone" binding:"required,phone"`
}

func (r *CreateUserRequest) Validate() error {
	details := make([]*shared.ErrorDetail, 0)
	if err := r.Email.Validate(); err != nil {
		details = append(details, &shared.ErrorDetail{
			Field:   "Email",
			Message: err.Error(),
		})
	}
	if err := r.Phone.Validate(); err != nil {
		details = append(details, &shared.ErrorDetail{
			Field:   "Email",
			Message: err.Error(),
		})
	}
	if len(details) == 0 {
		return nil
	}
	return domain.NewInvalidInputDataError("request is not valid", details...)
}

type UpdateUserRequest struct {
	Name  string       `json:"name" binding:"required"`
	Email domain.Email `json:"email" binding:"required,seller_email"`
	Phone domain.Phone `json:"phone" binding:"required,phone"`
}

func (r *UpdateUserRequest) Validate() error {
	details := make([]*shared.ErrorDetail, 0)
	if err := r.Email.Validate(); err != nil {
		details = append(details, &shared.ErrorDetail{
			Field:   "Email",
			Message: err.Error(),
		})
	}
	if err := r.Phone.Validate(); err != nil {
		details = append(details, &shared.ErrorDetail{
			Field:   "Email",
			Message: err.Error(),
		})
	}
	if len(details) == 0 {
		return nil
	}
	return domain.NewInvalidInputDataError("request is not valid", details...)
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
		log := sellerlog.FromContext(ctx, u.logger)
		log.Sugar().Errorf("error occured: %w", err)
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
		log := sellerlog.FromContext(ctx, u.logger)
		log.Sugar().Errorf("error occured: %w", err)
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
			return nil, domain.NewUserNotFoundError(id, err)
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
			return domain.NewUserNotFoundError(id, err)
		}
		return err
	}
	user.Deactivate()
	if err = u.userRepository.Delete(ctx, user.ID); err != nil {
		return err
	}
	return nil
}
