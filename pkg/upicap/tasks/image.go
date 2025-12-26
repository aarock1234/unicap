package tasks

import (
	"encoding/json"
	"fmt"

	"upicap/pkg/upicap"
)

// NumericMode specifies character type constraints for image recognition
type NumericMode int

const (
	NumericModeAny         NumericMode = 0 // no preference
	NumericModeNumbersOnly NumericMode = 1 // only numbers
	NumericModeLettersOnly NumericMode = 2 // only letters
	NumericModeEither      NumericMode = 3 // numbers OR letters
	NumericModeBoth        NumericMode = 4 // numbers AND letters
)

// ImageToTextTask represents an image recognition task
type ImageToTextTask struct {
	Body            string
	WebsiteURL      string
	Module          string
	Numeric         NumericMode
	Math            bool
	MinLength       int
	MaxLength       int
	Case            bool
	Phrase          bool
	Comment         string
	ImgInstructions string
	LanguagePool    string
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
