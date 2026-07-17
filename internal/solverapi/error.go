package solverapi

import "github.com/aarock1234/unicap"

// ErrorMapper maps provider-specific error codes to unicap sentinel errors.
type ErrorMapper struct {
	name     string
	mappings map[string]error
}

// NewErrorMapper creates an error mapper for the named provider.
func NewErrorMapper(name string) *ErrorMapper {
	return &ErrorMapper{
		name:     name,
		mappings: make(map[string]error),
	}
}

// Map registers a provider error code to a sentinel error. It returns the
// mapper to allow chaining.
func (m *ErrorMapper) Map(code string, sentinel error) *ErrorMapper {
	m.mappings[code] = sentinel

	return m
}

// Error builds a structured provider error. When the code maps to a known
// sentinel, the returned error wraps it (so errors.Is works) and is marked
// non-retriable; otherwise it is considered retriable.
func (m *ErrorMapper) Error(code, message string) *unicap.Error {
	if sentinel, ok := m.mappings[code]; ok {
		return unicap.NewError(code, message, m.name, false, sentinel)
	}

	return unicap.NewError(code, message, m.name, true, nil)
}

// StandardErrorMapper returns an error mapper preloaded with the common
// sentinel mappings for the given provider error codes.
func StandardErrorMapper(name string, keyInvalidCodes, balanceCodes, taskNotFoundCodes, invalidTaskCodes []string) *ErrorMapper {
	m := NewErrorMapper(name)

	for _, code := range keyInvalidCodes {
		m.Map(code, unicap.ErrInvalidAPIKey)
	}

	for _, code := range balanceCodes {
		m.Map(code, unicap.ErrInsufficientFunds)
	}

	for _, code := range taskNotFoundCodes {
		m.Map(code, unicap.ErrTaskNotFound)
	}

	for _, code := range invalidTaskCodes {
		m.Map(code, unicap.ErrInvalidTask)
	}

	return m
}
