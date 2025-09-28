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
	ErrCodePaymentNotFound      = "PAYMENT_NOT_FOUND"
	ErrCodeInvalidPaymentMethod = "INVALID_PAYMENT_METHOD"
	ErrCodeInvalidAmount        = "INVALID_AMOUNT"
	ErrCodePaymentFailed        = "PAYMENT_FAILED"
	ErrCodeInvalidUUID          = "INVALID_UUID"
	ErrCodeInternalError        = "INTERNAL_ERROR"
	ErrCodeValidationError      = "VALIDATION_ERROR"
)

// Error constructors
func NewPaymentNotFoundError(paymentUUID string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodePaymentNotFound,
		Message: fmt.Sprintf("payment %s not found", paymentUUID),
	}
}

func NewInvalidPaymentMethodError(method PaymentMethod) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidPaymentMethod,
		Message: fmt.Sprintf("invalid payment method: %s", method),
	}
}

func NewInvalidAmountError(amount float64) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidAmount,
		Message: fmt.Sprintf("invalid amount: %f", amount),
	}
}

func NewPaymentFailedError(err error) *ServiceError {
	return &ServiceError{
		Code:    ErrCodePaymentFailed,
		Message: "payment processing failed",
		Err:     err,
	}
}

func NewInvalidUUIDError(uuid string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidUUID,
		Message: fmt.Sprintf("invalid UUID: %s", uuid),
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
