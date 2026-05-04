package transport

type ErrorResponse struct {
	Error *Error `json:"error"`
}

type Error struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details []*ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
