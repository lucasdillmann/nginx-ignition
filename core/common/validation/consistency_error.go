package validation

type ConsistencyError struct {
	Violations []ConsistencyViolation
}

type ConsistencyViolation struct {
	Path    string
	Message string
}

func (err ConsistencyError) Error() string {
	return "One or more problems where found in the provided value"
}

func NewError(violations []ConsistencyViolation) *ConsistencyError {
	return &ConsistencyError{Violations: violations}
}

func SingleFieldError(
	path string,
	message string,
) *ConsistencyError {
	violation := ConsistencyViolation{
		Path:    path,
		Message: message,
	}

	return NewError([]ConsistencyViolation{violation})
}
