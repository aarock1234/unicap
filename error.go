package unicap

import (
	"errors"
	"fmt"
)

// Sentinel errors for conditions callers commonly handle. Provider errors wrap
// one of these when the provider's error code maps to a known condition, so
// callers can match them with errors.Is.
var (
	// ErrInvalidAPIKey reports that the provider rejected the supplied API key.
	ErrInvalidAPIKey = errors.New("invalid api key")

	// ErrInsufficientFunds reports that the provider account has no balance.
	ErrInsufficientFunds = errors.New("insufficient funds")

	// ErrTaskNotFound reports that the provider does not know the task ID, or
	// its result has already expired.
	ErrTaskNotFound = errors.New("task not found")

	// ErrInvalidTask reports that a task is missing required fields or carries
	// values the provider rejected.
	ErrInvalidTask = errors.New("invalid task parameters")

	// ErrUnsupportedTask reports that a provider cannot solve the given task
	// type.
	ErrUnsupportedTask = errors.New("unsupported task type")

	// ErrUnsupportedOption reports that a task sets an option the selected
	// provider cannot honour. unicap rejects such tasks instead of silently
	// dropping the option.
	ErrUnsupportedOption = errors.New("unsupported task option")

	// ErrUnsupportedOperation reports that a provider does not offer an
	// operation such as balance lookups or answer reports.
	ErrUnsupportedOperation = errors.New("unsupported operation")

	// ErrUnsolvable reports that the provider could not solve the captcha.
	// Submitting a fresh task may succeed.
	ErrUnsolvable = errors.New("captcha unsolvable")

	// ErrUnavailable reports that the provider is temporarily unable to accept
	// work, for example because no worker slot is free or a rate limit was
	// hit. Repeating the request later may succeed.
	ErrUnavailable = errors.New("provider temporarily unavailable")

	// ErrProxy reports that the provider could not use the supplied proxy, for
	// example because it refused connections or is banned by the target site.
	ErrProxy = errors.New("proxy rejected")

	// ErrIPBlocked reports that the provider blocked the caller's IP address
	// or account.
	ErrIPBlocked = errors.New("ip or account blocked")

	// ErrTaskFailed reports that the provider marked the task as failed
	// without a more specific reason.
	ErrTaskFailed = errors.New("task failed")

	// ErrTimeout reports that polling exceeded the configured timeout.
	ErrTimeout = errors.New("task timeout")

	// ErrNilProvider reports that a nil provider was passed to New.
	ErrNilProvider = errors.New("provider cannot be nil")
)

// Error is a failure reported by a provider. It carries the provider's raw
// error code and message and wraps a sentinel error when the code maps to a
// known condition.
type Error struct {
	// Provider is the name of the provider that reported the error.
	Provider string

	// Code is the provider's raw error code, such as "ERROR_ZERO_BALANCE".
	Code string

	// Message is the provider's human-readable error description.
	Message string

	// Err is the sentinel error the code maps to, or nil for codes unicap
	// does not classify.
	Err error
}

// Error returns the provider, message, and code.
func (e *Error) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("%s: %s", e.Provider, e.Code)
	}

	return fmt.Sprintf("%s: %s (code: %s)", e.Provider, e.Message, e.Code)
}

// Unwrap returns the sentinel error the code maps to, if any.
func (e *Error) Unwrap() error {
	return e.Err
}

// Retryable reports whether submitting the task again may succeed without
// changing it.
func (e *Error) Retryable() bool {
	return errors.Is(e.Err, ErrUnsolvable) || errors.Is(e.Err, ErrUnavailable)
}

// IsRetryable reports whether err indicates that submitting the task again
// may succeed without changing it.
func IsRetryable(err error) bool {
	return errors.Is(err, ErrUnsolvable) || errors.Is(err, ErrUnavailable)
}

// isTransient reports whether a polling error is worth repeating the same
// request for. Transport failures and temporary provider unavailability are
// transient; any other provider error is final for the task.
func isTransient(err error) bool {
	if errors.Is(err, ErrUnavailable) {
		return true
	}

	_, isProviderError := errors.AsType[*Error](err)

	return !isProviderError
}
