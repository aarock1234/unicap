package unicap

import (
	"errors"
	"fmt"
)

// Error represents a captcha solving error
type Error struct {
	Code      string
	Message   string
	Provider  string
	Retriable bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s (code: %s)", e.Provider, e.Message, e.Code)
}

// Sentinel errors for common conditions
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
)
