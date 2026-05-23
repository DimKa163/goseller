package usecase

import (
	"context"
	"errors"

	"github.com/DimKa163/goseller/internal/category/domain"
	"github.com/DimKa163/goseller/internal/dberror"
	"go.uber.org/zap"
)

type (
	CategoryRequest struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
)

type CategoryAppService struct {
	logger *zap.Logger
	rep    domain.CategoryRepository
}

func NewCategoryAppService(logger *zap.Logger, rep domain.CategoryRepository) *CategoryAppService {
	return &CategoryAppService{
		logger: logger,
		rep:    rep,
	}
}

func (s *CategoryAppService) GetByID(ctx context.Context, id domain.CategoryID) (*domain.Category, error) {
	cat, err := s.rep.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, dberror.ErrNoRows) {
			return nil, newCategoryNotFoundError(id.String(), err)
		}
		return nil, err
	}
	return cat, nil
}

func (s *CategoryAppService) Create(ctx context.Context, request *CategoryRequest) (domain.CategoryID, error) {
	category := domain.NewCategory(request.Name, request.Description)
	id, err := s.rep.Insert(ctx, category)
	if err != nil {
		var dbErr *dberror.ResourceAlreadyExistError
		if errors.As(err, &dbErr) {
			return domain.CategoryID{}, newCategoryAlreadyExistError(request.Name, err)
		}
		return domain.CategoryID{}, err
	}
	return id, nil
}

func (s *CategoryAppService) Update(ctx context.Context, id domain.CategoryID, request *CategoryRequest) (*domain.Category, error) {
	cat, err := s.rep.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, dberror.ErrNoRows) {
			return nil, newCategoryNotFoundError(id, err)
		}
		return nil, err
	}
	cat.Name = request.Name
	cat.Description = request.Description
	updatedCat, err := s.rep.Update(ctx, cat)
	if err != nil {
		var dbErr *dberror.ResourceAlreadyExistError
		if errors.As(err, &dbErr) {
			return nil, newCategoryAlreadyExistError(request.Name, dbErr)
		}
		return nil, err
	}
	return updatedCat, nil
}

func (s *CategoryAppService) Delete(ctx context.Context, id domain.CategoryID) error {
	err := s.rep.Deactivate(ctx, id)
	if err != nil {
		if errors.Is(err, dberror.ErrNoRows) {
			return newCategoryNotFoundError(id, err)
		}
		return err
	}
	return nil
}
