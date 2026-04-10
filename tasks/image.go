package tasks

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// NumericMode specifies character type constraints for image recognition
type NumericMode int

const (
	// NumericModeAny applies no numeric constraint.
	NumericModeAny NumericMode = 0 // no preference
	// NumericModeNumbersOnly requires only numbers.
	NumericModeNumbersOnly NumericMode = 1 // only numbers
	// NumericModeLettersOnly requires only letters.
	NumericModeLettersOnly NumericMode = 2 // only letters
	// NumericModeEither requires numbers or letters.
	NumericModeEither NumericMode = 3 // numbers OR letters
	// NumericModeBoth requires both numbers and letters.
	NumericModeBoth NumericMode = 4 // numbers AND letters
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

// Type returns the SDK task type identifier.
func (t *ImageToTextTask) Type() unicap.TaskType {
	return unicap.TaskTypeImageToText
}

// Validate ensures required fields are present
func (t *ImageToTextTask) Validate() error {
	if t.Body == "" {
		return fmt.Errorf("body: %w", unicap.ErrInvalidTask)
	}
	return nil
}
