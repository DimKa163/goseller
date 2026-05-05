package dberror

import "errors"

var ErrNoRows = errors.New("No rows in result set")

var ErrDuplicateKey = errors.New("Duplicate key value violates unique constraint")
