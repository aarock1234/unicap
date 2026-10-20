package solverapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aarock1234/unicap/v3"
)

// TaskID is a provider task identifier. It decodes from either a JSON string
// (CapSolver) or a JSON number (2Captcha, Anti-Captcha) and always renders
// as a string.
type TaskID string

// UnmarshalJSON accepts both a JSON string and a JSON number.
func (t *TaskID) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}

	if data[0] == '"' {
		var text string
		if err := json.Unmarshal(data, &text); err != nil {
			return fmt.Errorf("unmarshaling task id: %w", err)
		}

		*t = TaskID(text)

		return nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err != nil {
		return fmt.Errorf("unmarshaling task id: %w", err)
	}

	*t = TaskID(number.String())

	return nil
}

// String returns the task ID as a string.
func (t TaskID) String() string {
	return string(t)
}

// TaskIDFormat is the JSON type a provider expects for taskId in requests.
type TaskIDFormat int

const (
	// TaskIDString sends task IDs as JSON strings.
	TaskIDString TaskIDFormat = iota

	// TaskIDNumber sends task IDs as JSON integers.
	TaskIDNumber
)

// encode converts taskID to the JSON value the format requires.
func (f TaskIDFormat) encode(taskID string) (any, error) {
	if taskID == "" {
		return nil, fmt.Errorf("task id is empty: %w", unicap.ErrTaskNotFound)
	}

	if f == TaskIDString {
		return taskID, nil
	}

	if _, err := strconv.ParseUint(taskID, 10, 64); err != nil {
		return nil, fmt.Errorf("task id %q is not numeric: %w", taskID, unicap.ErrTaskNotFound)
	}

	return json.Number(taskID), nil
}
