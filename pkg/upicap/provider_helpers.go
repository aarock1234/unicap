package unicap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// BaseHTTPClient provides common HTTP functionality for provider implementations
type BaseHTTPClient struct {
	HTTPClient *http.Client
	Logger     *slog.Logger
	BaseURL    string
}

// DoJSON sends a JSON request and decodes the JSON response
func (c *BaseHTTPClient) DoJSON(ctx context.Context, endpoint string, reqBody, respBody any) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	c.Logger.DebugContext(ctx, "sending request",
		slog.String("endpoint", endpoint),
		slog.String("body", string(jsonData)),
	)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	c.Logger.DebugContext(ctx, "received response",
		slog.Int("status_code", resp.StatusCode),
		slog.String("body", string(body)),
	)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if err := json.Unmarshal(body, respBody); err != nil {
		return fmt.Errorf("unmarshaling response: %w", err)
	}

	return nil
}

// ErrorMapper provides a helper for mapping provider error codes to sentinel errors
type ErrorMapper struct {
	providerName string
	mappings     map[string]error
}

// NewErrorMapper creates a new error mapper for a provider
func NewErrorMapper(providerName string) *ErrorMapper {
	return &ErrorMapper{
		providerName: providerName,
		mappings:     make(map[string]error),
	}
}

// Map registers an error code to sentinel error mapping
func (m *ErrorMapper) Map(code string, sentinelErr error) *ErrorMapper {
	m.mappings[code] = sentinelErr
	return m
}

// MapError maps a provider error code to a wrapped sentinel error
func (m *ErrorMapper) MapError(code, message string) error {
	err := &Error{
		Code:     code,
		Message:  message,
		Provider: m.providerName,
	}

	if sentinelErr, ok := m.mappings[code]; ok {
		err.Retriable = false
		return fmt.Errorf("%w: %s", sentinelErr, message)
	}

	err.Retriable = true
	return err
}

// StandardErrorMapper returns a pre-configured error mapper with common mappings
func StandardErrorMapper(providerName string, keyInvalidCodes, balanceCodes, taskNotFoundCodes, invalidTaskCodes []string) *ErrorMapper {
	m := NewErrorMapper(providerName)

	for _, code := range keyInvalidCodes {
		m.Map(code, ErrInvalidAPIKey)
	}
	for _, code := range balanceCodes {
		m.Map(code, ErrInsufficientFunds)
	}
	for _, code := range taskNotFoundCodes {
		m.Map(code, ErrTaskNotFound)
	}
	for _, code := range invalidTaskCodes {
		m.Map(code, ErrInvalidTask)
	}

	return m
}
