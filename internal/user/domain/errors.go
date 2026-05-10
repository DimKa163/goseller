package domain

import (
	"errors"
	"fmt"

	"github.com/DimKa163/goseller/internal/shared"
)

var ErrInvalidUserID = errors.New("invalid user ID")

type UserNotFoundError struct {
	key      any
	innerErr error
}

// Details implements [shared.SellerError].
func (u *UserNotFoundError) Details() []*shared.ErrorDetail {
	return []*shared.ErrorDetail{}
}

func (u *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found with key %d", u.key)
}

func (u *UserNotFoundError) GetCode() shared.ErrorCode {
	return shared.ErrorCodeResourceNotFound
}

func (u *UserNotFoundError) Unwrap() error {
	return u.innerErr
}

func NewUserNotFoundError(key any, err error) shared.SellerError {
	return &UserNotFoundError{key: key, innerErr: err}
}

type UserAlreadyExistError struct {
	value    any
	details  []*shared.ErrorDetail
	innerErr error
}

// Details implements [shared.SellerError].
func (u *UserAlreadyExistError) Details() []*shared.ErrorDetail {
	return u.details
}

// Unwrap implements [UserError].
func (u *UserAlreadyExistError) Unwrap() error {
	return u.innerErr
}

// Error implements [UserError].
func (u *UserAlreadyExistError) Error() string {
	return "resource already exists"
}

// GetCode implements [UserError].
func (u *UserAlreadyExistError) GetCode() shared.ErrorCode {
	return shared.ErrorCodeResourceAlreadyExists
}

func NewUserAlreadyExistError(value any, err error, details ...*shared.ErrorDetail) shared.SellerError {
	return &UserAlreadyExistError{value: value, details: details, innerErr: err}
}

type InvalidInputDataError struct {
	message string
	details []*shared.ErrorDetail
}

// Details implements [shared.SellerError].
func (i *InvalidInputDataError) Details() []*shared.ErrorDetail {
	return i.details
}

// Error implements [shared.SellerError].
func (i *InvalidInputDataError) Error() string {
	return i.message
}

// GetCode implements [shared.SellerError].
func (i *InvalidInputDataError) GetCode() shared.ErrorCode {
	return shared.ErrorCodeBadInputData
}

// Unwrap implements [shared.SellerError].
func (i *InvalidInputDataError) Unwrap() error {
	return i
}

func NewInvalidInputDataError(message string, details ...*shared.ErrorDetail) shared.SellerError {
	return &InvalidInputDataError{
		message: message,
		details: details,
	}
}
