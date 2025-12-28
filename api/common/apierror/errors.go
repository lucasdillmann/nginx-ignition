package apierror

type APIError struct {
	Message    string
	StatusCode int
}

func New(statusCode int, message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e APIError) Error() string {
	return e.Message
}
