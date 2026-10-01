package utils

type AppError struct {
	Status  int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewBadRequest(message string) error {
	return &AppError{
		Status:  400,
		Message: message,
	}
}

func NewUnauthorized(message string) error {
	return &AppError{
		Status:  401,
		Message: message,
	}
}

func NewForbidden(message string) error {
	return &AppError{
		Status:  403,
		Message: message,
	}
}

func NewNotFound(message string) error {
	return &AppError{
		Status:  404,
		Message: message,
	}
}

func NewInternal(message string, err error) error {
	return &AppError{
		Status:  500,
		Message: message,
		Err:     err,
	}
}