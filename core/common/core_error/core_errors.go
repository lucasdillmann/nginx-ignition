package core_error

type CoreError struct {
	Message     string
	UserRelated bool
}

func New(message string, blameUser bool) *CoreError {
	return &CoreError{Message: message, UserRelated: blameUser}
}

func (e CoreError) Error() string {
	return e.Message
}
