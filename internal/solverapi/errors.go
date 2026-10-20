package solverapi

import (
	"strconv"

	"github.com/aarock1234/unicap/v3"
)

// ErrorRules maps provider error codes to unicap sentinel errors. Codes that
// are not listed produce a *unicap.Error without a sentinel, which is
// neither retryable nor transient.
type ErrorRules map[string]error

// classify builds the unicap error for a failed response envelope.
func (r ErrorRules) classify(provider string, envelope *Envelope) *unicap.Error {
	code := envelope.ErrorCode
	if code == "" {
		code = "ERROR_ID_" + strconv.Itoa(envelope.ErrorID)
	}

	return &unicap.Error{
		Provider: provider,
		Code:     code,
		Message:  envelope.ErrorDescription,
		Err:      r[code],
	}
}
