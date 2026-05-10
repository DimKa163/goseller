package shared

import "net/http"

type ErrorCode int

const (
	ErrorCodeNone ErrorCode = iota
	ErrorCodeInvalidID
	ErrorCodeBadInputData
	ErrorCodeResourceNotFound
	ErrorCodeResourceAlreadyExists
	ErrorCodeInternalServerError
)

type SellerError interface {
	error
	GetCode() ErrorCode
	Unwrap() error
	Details() []*ErrorDetail
}
type ErrorDetail struct {
	Field   string
	Message string
}

func ErrorCodeToHttpStatusCode(code ErrorCode) int {
	switch code {
	case ErrorCodeInvalidID:
		return http.StatusBadRequest
	case ErrorCodeBadInputData:
		return http.StatusBadRequest
	case ErrorCodeResourceAlreadyExists:
		return http.StatusConflict
	case ErrorCodeResourceNotFound:
		return http.StatusNotFound
	default:
		return http.StatusOK
	}
}
