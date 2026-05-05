package domain

import "errors"

var ErrInvalidUserID = errors.New("invalid user ID")

var ErrUserNotFound = errors.New("user not found")
