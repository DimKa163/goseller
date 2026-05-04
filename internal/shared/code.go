package shared

type ErrorCode int

const (
	ErrorCodeNone ErrorCode = iota
	ErrorCodeInvalidID
	ErrorCodeResourceNotFound
)
