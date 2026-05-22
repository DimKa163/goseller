package domain

import "github.com/DimKa163/goseller/internal/shared"

type CategoryIDInvalidError struct {
	value    string
	innerErr error
}

// Details implements [shared.SellerError].
func (c CategoryIDInvalidError) Details() []*shared.ErrorDetail {
	return []*shared.ErrorDetail{}
}

// Error implements [shared.SellerError].
func (c CategoryIDInvalidError) Error() string {
	return "invalid category ID"
}

// GetCode implements [shared.SellerError].
func (c CategoryIDInvalidError) GetCode() shared.ErrorCode {
	return shared.ErrorCodeInvalidID
}

// Unwrap implements [shared.SellerError].
func (c CategoryIDInvalidError) Unwrap() error {
	return c.innerErr
}

func newCategoryIDInvalidError(value string, err error) shared.SellerError {
	return CategoryIDInvalidError{
		value:    value,
		innerErr: err,
	}
}
