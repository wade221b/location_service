package errors

import (
	"fmt"

	"github.com/your-username/your-project/src/constants"
)

// ServiceError is a custom error type that includes a code and a message.
type ServiceError struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewServiceError creates a new ServiceError.
// Optionally, you can wrap an underlying error by passing it as err.
func NewServiceError(code string, err error) error {
	// Look up the default message from the constants package.
	msg, ok := constants.ErrorMessages[code]
	if !ok {
		msg = "Unknown error"
	}
	// If an underlying error is provided, append its message.
	if err != nil {
		msg = fmt.Sprintf("%s: %v", msg, err)
	}
	return &ServiceError{
		Code:    code,
		Message: msg,
	}
}
