package domain

import "errors"

var ErrInvalidUserID = errors.New("invalid user ID")

var ErrUserNotFound = errors.New("user not found")

var ErrPhoneAlreadyExists = errors.New("phone already exists")

var ErrEmailAlreadyExists = errors.New("email already exists")
