package unicap

import (
	"errors"
	"fmt"
)

// Error represents a captcha solving error returned by a provider. It carries
// the provider's raw error code and message, and may wrap a sentinel error for
// use with errors.Is.
type Error struct {
	Code      string
	Message   string
	Provider  string
	Retriable bool
	wrapped   error
}

// NewError builds a structured provider error. When wrapped is non-nil, the
// returned error unwraps to it so callers can match it with errors.Is.
func NewError(code, message, provider string, retriable bool, wrapped error) *Error {
	return &Error{
		Code:      code,
		Message:   message,
		Provider:  provider,
		Retriable: retriable,
		wrapped:   wrapped,
	}
}

// Error returns the human-readable error message.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s (code: %s)", e.Provider, e.Message, e.Code)
}

// Unwrap returns the wrapped sentinel error, if any.
func (e *Error) Unwrap() error {
	return e.wrapped
}

// Sentinel errors for common conditions.
var (
	// ErrInvalidAPIKey reports that the provider rejected the supplied API key.
	ErrInvalidAPIKey = errors.New("invalid api key")
	// ErrInsufficientFunds reports that the provider account has no balance.
	ErrInsufficientFunds = errors.New("insufficient funds")
	// ErrTaskNotFound reports that the provider no longer knows about a task ID.
	ErrTaskNotFound = errors.New("task not found")
	// ErrTimeout reports that polling exceeded the configured timeout.
	ErrTimeout = errors.New("task timeout")
	// ErrInvalidTask reports that a task is missing required fields.
	ErrInvalidTask = errors.New("invalid task parameters")
	// ErrUnsupportedTask reports that a provider cannot solve the given task type.
	ErrUnsupportedTask = errors.New("unsupported task type")
	// ErrNilProvider reports that a nil provider was passed to New.
	ErrNilProvider = errors.New("provider cannot be nil")
)
