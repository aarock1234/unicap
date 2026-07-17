package solverapi

import (
	"errors"
	"testing"

	"github.com/aarock1234/unicap"
)

func TestErrorMapperError(t *testing.T) {
	mapper := StandardErrorMapper(
		"testprovider",
		[]string{"ERROR_KEY_INVALID"},
		[]string{"ERROR_ZERO_BALANCE"},
		[]string{"ERROR_TASK_NOT_FOUND"},
		[]string{"ERROR_BAD_TASK"},
	)

	tests := []struct {
		name          string
		code          string
		wantSentinel  error
		wantRetriable bool
	}{
		{
			name:          "invalid key maps to sentinel",
			code:          "ERROR_KEY_INVALID",
			wantSentinel:  unicap.ErrInvalidAPIKey,
			wantRetriable: false,
		},
		{
			name:          "zero balance maps to sentinel",
			code:          "ERROR_ZERO_BALANCE",
			wantSentinel:  unicap.ErrInsufficientFunds,
			wantRetriable: false,
		},
		{
			name:          "unknown code is retriable",
			code:          "ERROR_SOMETHING_ELSE",
			wantSentinel:  nil,
			wantRetriable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapper.Error(tt.code, "message")

			if err.Retriable != tt.wantRetriable {
				t.Errorf("Retriable = %v, want %v", err.Retriable, tt.wantRetriable)
			}

			if err.Provider != "testprovider" {
				t.Errorf("Provider = %q, want testprovider", err.Provider)
			}

			if err.Code != tt.code {
				t.Errorf("Code = %q, want %q", err.Code, tt.code)
			}

			if tt.wantSentinel != nil && !errors.Is(err, tt.wantSentinel) {
				t.Errorf("errors.Is(%v, %v) = false, want true", err, tt.wantSentinel)
			}
		})
	}
}
