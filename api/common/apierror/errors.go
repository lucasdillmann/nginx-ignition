package apierror

type ApiError struct {
	StatusCode int
	Message    string
}

func New(statusCode int, message string) *ApiError {
	return &ApiError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e ApiError) Error() string {
	return e.Message
}
