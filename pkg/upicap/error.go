package upicap

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
	ErrInvalidAPIKey     = errors.New("invalid api key")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrTaskNotFound      = errors.New("task not found")
	ErrTimeout           = errors.New("task timeout")
	ErrInvalidTask       = errors.New("invalid task parameters")
)
