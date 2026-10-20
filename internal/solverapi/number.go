package solverapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Number is a numeric field that providers send either as a JSON number or
// as a numeric string, such as cost "0.0012" or solveCount "0". Empty
// strings and null decode as zero.
type Number float64

// UnmarshalJSON accepts a JSON number, a numeric string, an empty string, or
// null.
func (n *Number) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}

	text := string(data)
	if data[0] == '"' {
		if err := json.Unmarshal(data, &text); err != nil {
			return fmt.Errorf("decoding numeric string: %w", err)
		}
	}

	if text == "" {
		return nil
	}

	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return fmt.Errorf("decoding number %q: %w", text, err)
	}

	*n = Number(value)

	return nil
}

// Time interprets the number as Unix seconds. Zero yields the zero time.
func (n Number) Time() time.Time {
	if n <= 0 {
		return time.Time{}
	}

	return time.Unix(int64(n), 0)
}
