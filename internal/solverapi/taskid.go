package solverapi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// TaskID is a captcha task identifier. It decodes from either a JSON string
// (2Captcha, CapSolver) or a JSON number (Anti-Captcha) and always renders as
// a string.
type TaskID string

// UnmarshalJSON accepts both a JSON string and a JSON number.
func (t *TaskID) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}

	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return fmt.Errorf("unmarshaling task id: %w", err)
		}

		*t = TaskID(s)

		return nil
	}

	var n json.Number
	if err := json.Unmarshal(data, &n); err != nil {
		return fmt.Errorf("unmarshaling task id: %w", err)
	}

	*t = TaskID(n.String())

	return nil
}

// String returns the task ID as a string.
func (t TaskID) String() string {
	return string(t)
}
