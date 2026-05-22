package infrastructure

import (
	"context"
	"errors"

	"github.com/DimKa163/goseller/internal/category/domain"
	"github.com/DimKa163/goseller/internal/dberror"
	"github.com/DimKa163/goseller/internal/shared"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	findByIDQuery   = `SELECT id, created_at, updated_at, name, description, inactive FROM categories WHERE id = $1 and inactive = false`
	insertQuery     = `INSERT INTO categories (id, name, description) VALUES ($1, $2, $3) RETURNING id`
	updateQuery     = `UPDATE categories SET name = $1, description = $2, updated_at = now() WHERE id = $3 and inactive = false RETURNING id, created_at, updated_at, name, description, inactive`
	deactivateQuery = `UPDATE categories SET inactive = true, updated_at = now() WHERE id = $1`
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindByID(ctx context.Context, id domain.CategoryID) (*domain.Category, error) {
	var category domain.Category
	if err := r.db.QueryRow(ctx, findByIDQuery, id).Scan(&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.Name,
		&category.Description,
		&category.Inactive); err != nil {

		return nil, handleError(err)
	}
	return &category, nil

}

func (r *CategoryRepository) Insert(ctx context.Context, category *domain.Category) (domain.CategoryID, error) {
	var id domain.CategoryID
	if err := r.db.QueryRow(ctx, insertQuery, category.ID, category.Name, category.Description).Scan(&id); err != nil {
		return domain.CategoryID{}, handleError(err)
	}
	return id, nil
}

func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	var updatedCategory domain.Category
	if err := r.db.QueryRow(ctx, updateQuery, category.Name, category.Description, category.ID).Scan(&updatedCategory.ID,
		&updatedCategory.CreatedAt,
		&updatedCategory.UpdatedAt,
		&updatedCategory.Name,
		&updatedCategory.Description,
		&updatedCategory.Inactive); err != nil {

		return nil, handleError(err)
	}
	return &updatedCategory, nil
}

func (r *CategoryRepository) Deactivate(ctx context.Context, id domain.CategoryID) error {
	_, err := r.db.Exec(ctx, deactivateQuery, id)
	return handleError(err)
}

func handleError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return dberror.ErrNoRows
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "udx_categories_name":
			return dberror.NewResourceAlreadyExistError("name", pgErr, &shared.ErrorDetail{
				Field:   pgErr.ColumnName,
				Message: "name already exist",
			})
		}
	}
	return err
}
