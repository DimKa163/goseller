package dberror

import (
	"errors"

	"github.com/DimKa163/goseller/internal/shared"
)

var ErrNoRows = errors.New("No rows in result set")

var ErrDuplicateKey = errors.New("Duplicate key value violates unique constraint")

type ResourceAlreadyExistError struct {
	value    any
	details  []*shared.ErrorDetail
	innerErr error
}

func (u *ResourceAlreadyExistError) Unwrap() error {
	return u.innerErr
}

// Error implements [UserError].
func (u *ResourceAlreadyExistError) Error() string {
	return "resource already exists"
}

func NewResourceAlreadyExistError(value any, err error, details ...*shared.ErrorDetail) *ResourceAlreadyExistError {
	return &ResourceAlreadyExistError{value: value, innerErr: err, details: details}
}
