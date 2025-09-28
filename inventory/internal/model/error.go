package model

import "fmt"

// ServiceError represents a service layer error
type ServiceError struct {
	Code    string
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

const (
	ErrCodePartNotFound    = "PART_NOT_FOUND"
	ErrCodeInvalidUUID     = "INVALID_UUID"
	ErrCodeInvalidFilter   = "INVALID_FILTER"
	ErrCodeInternalError   = "INTERNAL_ERROR"
	ErrCodeValidationError = "VALIDATION_ERROR"
)

// Error constructors
func NewPartNotFoundError(partUUID string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodePartNotFound,
		Message: fmt.Sprintf("part %s not found", partUUID),
	}
}

func NewInvalidUUIDError(uuid string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidUUID,
		Message: fmt.Sprintf("invalid UUID: %s", uuid),
	}
}

func NewInvalidFilterError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidFilter,
		Message: message,
	}
}

func NewInternalError(err error) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInternalError,
		Message: "internal service error",
		Err:     err,
	}
}

func NewValidationError(message string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeValidationError,
		Message: message,
	}
}
