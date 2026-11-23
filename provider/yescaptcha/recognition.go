package yescaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// imageToTextTask is YesCaptcha's asynchronous image-to-text task, the only
// image-to-text type that accepts answer constraints.
type imageToTextTask struct {
	Type      string `json:"type"`
	Body      string `json:"body"`
	Phrase    bool   `json:"phrase,omitempty"`
	Case      bool   `json:"case,omitempty"`
	Numeric   int    `json:"numeric,omitempty"`
	Math      bool   `json:"math,omitempty"`
	MinLength int    `json:"minLength,omitempty"`
	MaxLength int    `json:"maxLength,omitempty"`
}

// mapImageToText builds YesCaptcha's image-to-text task. YesCaptcha's numeric
// field has no value for the either and both modes.
func mapImageToText(task *tasks.ImageToTextTask) (imageToTextTask, error) {
	isMixedNumeric := task.Numeric == tasks.NumericModeEither || task.Numeric == tasks.NumericModeBoth

	if err := solverapi.Unsupported(
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
		solverapi.Setting("Module", task.Module != ""),
		solverapi.Setting("Images", len(task.Images) > 0),
		solverapi.Setting("Numeric (either or both)", isMixedNumeric),
		solverapi.Setting("Comment", task.Comment != ""),
		solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
		solverapi.Setting("LanguagePool", task.LanguagePool != ""),
	); err != nil {
		return imageToTextTask{}, err
	}

	return imageToTextTask{
		Type:      "ImageToTextTask",
		Body:      task.Body,
		Phrase:    task.Phrase,
		Case:      task.Case,
		Numeric:   int(task.Numeric),
		Math:      task.Math,
		MinLength: task.MinLength,
		MaxLength: task.MaxLength,
	}, nil
}
