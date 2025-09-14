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

// Common service error codes
const (
	ErrCodeOrderNotFound      = "ORDER_NOT_FOUND"
	ErrCodeInvalidStatus      = "INVALID_STATUS"
	ErrCodePartNotFound       = "PART_NOT_FOUND"
	ErrCodePaymentFailed      = "PAYMENT_FAILED"
	ErrCodeInvalidUUID        = "INVALID_UUID"
	ErrCodeInternalError      = "INTERNAL_ERROR"
	ErrCodeExternalServiceErr = "EXTERNAL_SERVICE_ERROR"
)

// Error constructors
func NewOrderNotFoundError(orderUUID string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeOrderNotFound,
		Message: fmt.Sprintf("order %s not found", orderUUID),
	}
}

func NewInvalidStatusError(currentStatus OrderStatus, expectedStatus OrderStatus) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeInvalidStatus,
		Message: fmt.Sprintf("invalid order status: expected %s, got %s", expectedStatus, currentStatus),
	}
}

func NewPartNotFoundError(partUUID string) *ServiceError {
	return &ServiceError{
		Code:    ErrCodePartNotFound,
		Message: fmt.Sprintf("part %s not found", partUUID),
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

func NewExternalServiceError(service string, err error) *ServiceError {
	return &ServiceError{
		Code:    ErrCodeExternalServiceErr,
		Message: fmt.Sprintf("external service %s error", service),
		Err:     err,
	}
}
