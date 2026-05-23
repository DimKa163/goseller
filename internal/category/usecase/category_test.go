package usecase

import (
	"context"
	"testing"

	"github.com/DimKa163/goseller/internal/category/domain"
	"github.com/DimKa163/goseller/internal/category/mocks"
	"github.com/DimKa163/goseller/internal/dberror"
	"github.com/beevik/guid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestGetByIDShouldBeSuccessful(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)
	uid := guid.New()
	category := domain.NewCategory("Test Category", "Test Description")
	id, _ := domain.NewCategoryID(uid.String())
	sut := fixture.withFindByID(ctx, id, category, nil).build()

	result, err := sut.GetByID(ctx, id)

	assert.NoError(t, err)
	assert.Equal(t, category, result)
}

func TestGetByIDShouldReturnNotFoundError(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)
	uid := guid.New()
	id, _ := domain.NewCategoryID(uid.String())
	sut := fixture.withFindByID(ctx, id, nil, dberror.ErrNoRows).build()
	result, err := sut.GetByID(ctx, id)

	var categoryNotFoundErr *CategoryNotFoundError
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorAs(t, err, &categoryNotFoundErr)
}

func TestCreateShouldBeSuccessful(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)

	request := &CategoryRequest{
		Name:        "Test Category",
		Description: "Test Description",
	}

	categoryID, _ := domain.NewCategoryID(guid.New().String())
	sut := fixture.withInsert(ctx, categoryID, nil).build()

	result, err := sut.Create(ctx, request)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, categoryID, result)
}

func TestCreateShouldReturnAlreadyExistError(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)

	request := &CategoryRequest{
		Name:        "Test Category",
		Description: "Test Description",
	}
	categoryID := domain.CategoryID{}
	sut := fixture.withInsert(ctx, categoryID, dberror.NewResourceAlreadyExistError("Category already exists", nil)).build()

	result, err := sut.Create(ctx, request)

	assert.Error(t, err)
	var categoryAlreadyExistErr *CategoryAlreadyExistError
	assert.ErrorAs(t, err, &categoryAlreadyExistErr)
	assert.Empty(t, result)
	assert.Equal(t, categoryID, result)
}

func TestUpdateShouldBeSuccessful(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)
	uid := guid.New()
	category := domain.NewCategory("Test Category", "Test Description")
	id, _ := domain.NewCategoryID(uid.String())

	request := &CategoryRequest{
		Name:        "Updated Category",
		Description: "Updated Description",
	}
	sut := fixture.withFindByID(ctx, id, category, nil).
		withUpdate(ctx, request.Name, request.Description, nil).
		build()
	result, err := sut.Update(ctx, id, request)

	assert.NoError(t, err)
	assert.Equal(t, request.Name, result.Name)
	assert.Equal(t, request.Description, result.Description)
}

func TestUpdateShouldReturnNotFoundError(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)
	uid := guid.New()
	id, _ := domain.NewCategoryID(uid.String())
	request := &CategoryRequest{
		Name:        "Updated Category",
		Description: "Updated Description",
	}
	sut := fixture.withFindByID(ctx, id, &domain.Category{}, dberror.ErrNoRows).
		withoutUpdate(ctx, request.Name, request.Description, nil).
		build()

	result, err := sut.Update(ctx, id, request)

	var categoryNotFoundErr *CategoryNotFoundError
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorAs(t, err, &categoryNotFoundErr)
}

func TestUpdateShouldReturnResourceAlreadyExistError(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)
	uid := guid.New()
	category := domain.NewCategory("Test Category", "Test Description")
	id, _ := domain.NewCategoryID(uid.String())

	request := &CategoryRequest{
		Name:        "Updated Category",
		Description: "Updated Description",
	}
	sut := fixture.withFindByID(ctx, id, category, nil).
		withUpdate(ctx, request.Name, request.Description,
			dberror.NewResourceAlreadyExistError("Category already exists", nil)).
		build()
	result, err := sut.Update(ctx, id, request)

	assert.Error(t, err)
	var categoryAlreadyExistErr *CategoryAlreadyExistError
	assert.ErrorAs(t, err, &categoryAlreadyExistErr)
	assert.Nil(t, result)
}

func TestDeleteShouldBeSuccessful(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)
	uid := guid.New()
	id, _ := domain.NewCategoryID(uid.String())
	sut := fixture.withDeactivate(ctx, id, nil).build()
	err := sut.Delete(ctx, id)
	assert.NoError(t, err)
}

func TestDeleteShouldReturnNotFoundError(t *testing.T) {
	ctx := context.Background()
	fixture := newCategoryAppFixture(t)
	uid := guid.New()
	id, _ := domain.NewCategoryID(uid.String())
	sut := fixture.withDeactivate(ctx, id, dberror.ErrNoRows).build()
	err := sut.Delete(ctx, id)
	var categoryNotFoundErr *CategoryNotFoundError
	assert.Error(t, err)
	assert.ErrorAs(t, err, &categoryNotFoundErr)
}

type categoryAppFixture struct {
	t            *testing.T
	categoryRepo *mocks.MockCategoryRepository
	logger       *zap.Logger
}

func newCategoryAppFixture(t *testing.T) *categoryAppFixture {
	ctrl := gomock.NewController(t)
	return &categoryAppFixture{
		t:            t,
		categoryRepo: mocks.NewMockCategoryRepository(ctrl),
		logger:       zap.NewNop(),
	}
}

func (f *categoryAppFixture) build() *CategoryAppService {
	return NewCategoryAppService(f.logger, f.categoryRepo)
}

func (f *categoryAppFixture) withFindByID(ctx context.Context, categoryId domain.CategoryID, category *domain.Category, err error) *categoryAppFixture {
	f.categoryRepo.EXPECT().FindByID(ctx, categoryId).Return(category, err).Times(1)
	return f
}

func (f *categoryAppFixture) withInsert(ctx context.Context, categoryID domain.CategoryID, err error) *categoryAppFixture {
	f.categoryRepo.EXPECT().Insert(ctx, gomock.Any()).Return(categoryID, err).Times(1)
	return f
}

func (f *categoryAppFixture) withUpdate(ctx context.Context, name, description string, err error) *categoryAppFixture {
	f.categoryRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, cat *domain.Category) (*domain.Category, error) {
		assert.Equal(f.t, name, cat.Name)
		assert.Equal(f.t, description, cat.Description)
		return cat, err
	}).Times(1)
	return f
}

func (f *categoryAppFixture) withoutUpdate(ctx context.Context, name, description string, err error) *categoryAppFixture {
	f.categoryRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, cat *domain.Category) (*domain.Category, error) {
		assert.Equal(f.t, name, cat.Name)
		assert.Equal(f.t, description, cat.Description)
		return cat, err
	}).Times(0)
	return f
}

func (f *categoryAppFixture) withDeactivate(ctx context.Context, categoryId domain.CategoryID, err error) *categoryAppFixture {
	f.categoryRepo.EXPECT().Deactivate(ctx, categoryId).Return(err).Times(1)
	return f
}

func (f *categoryAppFixture) withoutDeactivate(ctx context.Context, categoryId domain.CategoryID, err error) *categoryAppFixture {
	f.categoryRepo.EXPECT().Deactivate(ctx, categoryId).Return(err).Times(0)
	return f
}
