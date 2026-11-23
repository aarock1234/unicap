package noptcha

import (
	"slices"

	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// imageToTextParams are the params of Noptcha's image.text task. A zero
// length means no limit.
type imageToTextParams struct {
	Image         string            `json:"image"`
	CaseSensitive bool              `json:"case_sensitive,omitempty"`
	Phrase        bool              `json:"phrase,omitempty"`
	Math          bool              `json:"math,omitempty"`
	Numeric       tasks.NumericMode `json:"numeric,omitempty"`
	MinLength     int               `json:"min_length,omitempty"`
	MaxLength     int               `json:"max_length,omitempty"`
}

// numericModes are the numeric modes Noptcha has a value for: 0 for no
// requirement, 1 for digits only and 2 for no digits, which are unicap's own
// codes.
var numericModes = []tasks.NumericMode{
	tasks.NumericModeAny,
	tasks.NumericModeNumbersOnly,
	tasks.NumericModeLettersOnly,
}

// mapImageToText builds the params of Noptcha's image.text task.
func mapImageToText(task *tasks.ImageToTextTask) (imageToTextParams, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
		solverapi.Setting("Module", task.Module != ""),
		solverapi.Setting("Images", len(task.Images) > 0),
		solverapi.Setting("Numeric (either or both)", !slices.Contains(numericModes, task.Numeric)),
		solverapi.Setting("Comment", task.Comment != ""),
		solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
		solverapi.Setting("LanguagePool", task.LanguagePool != ""),
	); err != nil {
		return imageToTextParams{}, err
	}

	return imageToTextParams{
		Image:         task.Body,
		CaseSensitive: task.Case,
		Phrase:        task.Phrase,
		Math:          task.Math,
		Numeric:       task.Numeric,
		MinLength:     task.MinLength,
		MaxLength:     task.MaxLength,
	}, nil
}
