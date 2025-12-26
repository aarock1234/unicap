package tasks

import (
	"encoding/json"
	"fmt"

	"upicap/pkg/upicap"
)

// ImageToTextTask represents an image recognition task
type ImageToTextTask struct {
	Body      string
	Numeric   int
	Math      bool
	MinLength int
	MaxLength int
	Case      bool
	Comment   string
}

func (t *ImageToTextTask) Type() upicap.TaskType {
	return upicap.TaskTypeImageToText
}

func (t *ImageToTextTask) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

// Validate ensures required fields are present
func (t *ImageToTextTask) Validate() error {
	if t.Body == "" {
		return fmt.Errorf("body: %w", upicap.ErrInvalidTask)
	}
	return nil
}
