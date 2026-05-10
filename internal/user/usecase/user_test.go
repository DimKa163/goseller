package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/DimKa163/goseller/internal/dberror"
	"github.com/DimKa163/goseller/internal/shared"
	"github.com/DimKa163/goseller/internal/user/domain"
	"github.com/DimKa163/goseller/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCreateShouldBeSuccessful(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "Test"
	email := domain.Email("test@test.ru")
	phone := domain.Phone("00123456789")
	req := &CreateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.UserID(1), nil).Times(1)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	r, err := sut.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, domain.UserID(1), r)
}

func TestCreateShouldReturnErrorIfEmailAlreadyExists(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "Test"
	email := domain.Email("test@test.ru")
	phone := domain.Phone("00123456789")
	req := &CreateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.UserID(-1), dberror.ErrDuplicateKey).Times(1)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	id, err := sut.Create(ctx, req)

	assert.ErrorIs(t, err, dberror.ErrDuplicateKey)
	assert.Equal(t, domain.UserID(-1), id)
}

func TestCreateShouldReturnErrorIfPhoneAlreadyExists(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	name := "Test"
	email := domain.Email("test@test.ru")
	phone := domain.Phone("00123456789")
	req := &CreateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.UserID(-1), dberror.ErrDuplicateKey).Times(1)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	id, err := sut.Create(ctx, req)

	assert.ErrorIs(t, err, dberror.ErrDuplicateKey)
	assert.Equal(t, domain.UserID(-1), id)
}

func TestUpdateShouldBeSuccessful(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := domain.UserID(1)
	user := &domain.User{
		ID:        id,
		CreatedAt: time.Now().Add(time.Duration(-50 * time.Hour)),
		UpdatedAt: time.Now().Add(time.Duration(-5 * time.Hour)),
		Name:      "Test",
		Email:     domain.Email("test0@test.ru"),
		Phone:     domain.Phone("10123456789"),
	}
	name := "Test"
	email := domain.Email("test@test.ru")
	phone := domain.Phone("00123456789")
	req := &UpdateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetByID(ctx, id).Return(user, nil).Times(1)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(&domain.User{
		ID:        id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: time.Now(),
		Name:      name,
		Email:     email,
		Phone:     phone,
	}, nil).Times(1)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	r, err := sut.Update(ctx, id, req)

	assert.NoError(t, err)
	assert.Equal(t, name, r.Name)
	assert.Equal(t, email, r.Email)
	assert.Equal(t, phone, r.Phone)
}

func TestUpdateShouldReturnErrorIfUserNotFound(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := domain.UserID(1)
	name := "Test"
	email := domain.Email("test@test.ru")
	phone := domain.Phone("00123456789")
	req := &UpdateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetByID(ctx, id).Return(nil, domain.NewUserNotFoundError(id, pgx.ErrNoRows)).Times(1)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(&domain.User{
		ID:        id,
		CreatedAt: time.Now().Add(time.Duration(-50 * time.Hour)),
		UpdatedAt: time.Now(),
		Name:      name,
		Email:     email,
		Phone:     phone,
	}, nil).Times(0)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	r, err := sut.Update(ctx, id, req)
	var usErr shared.SellerError
	assert.ErrorAs(t, err, &usErr)
	assert.Nil(t, r)
}

func TestUpdateShouldReturnErrorIfEmailAlreadyExists(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := domain.UserID(1)
	user := &domain.User{
		ID:        id,
		CreatedAt: time.Now().Add(time.Duration(-50 * time.Hour)),
		UpdatedAt: time.Now().Add(time.Duration(-5 * time.Hour)),
		Name:      "Test",
		Email:     domain.Email("test0@test.ru"),
		Phone:     domain.Phone("10123456789"),
	}
	name := "Test"
	email := domain.Email("test@test.ru")
	phone := domain.Phone("00123456789")
	req := &UpdateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetByID(ctx, id).Return(user, nil).Times(1)
	var pgErr pgconn.PgError
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil, domain.NewUserAlreadyExistError("email", &pgErr, &shared.ErrorDetail{
		Field:   pgErr.ColumnName,
		Message: "email already exist",
	})).Times(1)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	r, err := sut.Update(ctx, id, req)

	var usErr shared.SellerError
	assert.ErrorAs(t, err, &usErr)
	assert.Nil(t, r)
}

func TestUpdateShouldReturnErrorIfPhoneAlreadyExists(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := domain.UserID(1)
	user := &domain.User{
		ID:        id,
		CreatedAt: time.Now().Add(time.Duration(-50 * time.Hour)),
		UpdatedAt: time.Now().Add(time.Duration(-5 * time.Hour)),
		Name:      "Test",
		Email:     domain.Email("test0@test.ru"),
		Phone:     domain.Phone("10123456789"),
	}
	name := "Test"
	email := domain.Email("test@test.ru")
	phone := domain.Phone("00123456789")
	req := &UpdateUserRequest{
		Name:  name,
		Email: email,
		Phone: phone,
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetByID(ctx, id).Return(user, nil).Times(1)
	var pgErr pgconn.PgError
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(nil, domain.NewUserAlreadyExistError("phone", &pgErr, &shared.ErrorDetail{
		Field:   pgErr.ColumnName,
		Message: "phone already exist",
	})).Times(1)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	r, err := sut.Update(ctx, id, req)

	var usErr shared.SellerError
	assert.ErrorAs(t, err, &usErr)
	assert.Nil(t, r)
}

func TestDeleteShouldBeSuccessful(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := domain.UserID(1)
	user := &domain.User{
		ID:        id,
		CreatedAt: time.Now().Add(time.Duration(-50 * time.Hour)),
		UpdatedAt: time.Now().Add(time.Duration(-5 * time.Hour)),
		Name:      "Test",
		Email:     domain.Email("test0@test.ru"),
		Phone:     domain.Phone("10123456789"),
	}
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetByID(ctx, id).Return(user, nil).Times(1)
	mockRepo.EXPECT().Delete(ctx, id).Return(nil).Times(1)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	err := sut.Delete(ctx, id)

	assert.NoError(t, err)
}

func TestDeleteShouldReturnErrorIfUserNotFound(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	id := domain.UserID(1)
	mockRepo := mocks.NewMockUserRepository(ctrl)
	mockRepo.EXPECT().GetByID(ctx, id).Return(nil, domain.NewUserNotFoundError(id, pgx.ErrNoRows)).Times(1)
	mockRepo.EXPECT().Delete(ctx, id).Return(nil).Times(0)
	logger := zap.NewNop()
	sut := NewUser(mockRepo, logger)

	err := sut.Delete(ctx, id)
	var usErr shared.SellerError
	assert.ErrorAs(t, err, &usErr)
}
