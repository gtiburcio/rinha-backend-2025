package apperror

type AppError struct {
	Ignorable bool
	Message   string
}

func (a AppError) Error() string {
	return a.Message
}

func NewAppError(ignorable bool, message string) AppError {
	return AppError{
		Ignorable: ignorable,
		Message:   message,
	}
}

func IsIgnorableError(err error) bool {
	a, ok := err.(AppError)
	if !ok {
		return false
	}

	return a.Ignorable
}
