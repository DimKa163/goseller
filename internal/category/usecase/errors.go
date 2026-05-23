package usecase

import (
	"fmt"

	"github.com/DimKa163/goseller/internal/shared"
)

type CategoryNotFoundError struct {
	key      any
	innerErr error
}

func (u *CategoryNotFoundError) Details() []*shared.ErrorDetail {
	return []*shared.ErrorDetail{}
}

func (u *CategoryNotFoundError) Error() string {
	return fmt.Sprintf("category not found with key %s", u.key)
}

func (u *CategoryNotFoundError) GetCode() shared.ErrorCode {
	return shared.ErrorCodeResourceNotFound
}

func (u *CategoryNotFoundError) Unwrap() error {
	return u.innerErr
}

func newCategoryNotFoundError(key any, err error) shared.SellerError {
	return &CategoryNotFoundError{key: key, innerErr: err}
}

type CategoryAlreadyExistError struct {
	value    any
	details  []*shared.ErrorDetail
	innerErr error
}

func (u *CategoryAlreadyExistError) Details() []*shared.ErrorDetail {
	return u.details
}

// Unwrap implements [UserError].
func (u *CategoryAlreadyExistError) Unwrap() error {
	return u.innerErr
}

// Error implements [UserError].
func (u *CategoryAlreadyExistError) Error() string {
	return "resource already exists"
}

// GetCode implements [UserError].
func (u *CategoryAlreadyExistError) GetCode() shared.ErrorCode {
	return shared.ErrorCodeResourceAlreadyExists
}

func newCategoryAlreadyExistError(value any, err error, details ...*shared.ErrorDetail) shared.SellerError {
	return &CategoryAlreadyExistError{value: value, innerErr: err, details: details}
}
