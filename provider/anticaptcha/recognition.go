package anticaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// imageToTextTask always sends case because Anti-Captcha treats an absent
// value as case-sensitive.
type imageToTextTask struct {
	Type         string `json:"type"`
	Body         string `json:"body"`
	Phrase       bool   `json:"phrase,omitempty"`
	Case         bool   `json:"case"`
	Numeric      int    `json:"numeric,omitempty"`
	Math         bool   `json:"math,omitempty"`
	MinLength    int    `json:"minLength,omitempty"`
	MaxLength    int    `json:"maxLength,omitempty"`
	Comment      string `json:"comment,omitempty"`
	WebsiteURL   string `json:"websiteURL,omitempty"`
	LanguagePool string `json:"languagePool,omitempty"`
}

type imageToCoordinatesTask struct {
	Type       string `json:"type"`
	Body       string `json:"body"`
	Comment    string `json:"comment,omitempty"`
	Mode       string `json:"mode,omitempty"`
	WebsiteURL string `json:"websiteURL,omitempty"`
}

func mapImageToText(task *tasks.ImageToTextTask) (imageToTextTask, error) {
	isMixedNumeric := task.Numeric == tasks.NumericModeEither || task.Numeric == tasks.NumericModeBoth

	if err := solverapi.Unsupported(
		solverapi.Setting("Module", task.Module != ""),
		solverapi.Setting("Images", len(task.Images) > 0),
		solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
		solverapi.Setting("Numeric (either or both)", isMixedNumeric),
	); err != nil {
		return imageToTextTask{}, err
	}

	return imageToTextTask{
		Type:         "ImageToTextTask",
		Body:         task.Body,
		Phrase:       task.Phrase,
		Case:         task.Case,
		Numeric:      int(task.Numeric),
		Math:         task.Math,
		MinLength:    task.MinLength,
		MaxLength:    task.MaxLength,
		Comment:      task.Comment,
		WebsiteURL:   task.WebsiteURL,
		LanguagePool: task.LanguagePool,
	}, nil
}

func mapCoordinates(task *tasks.CoordinatesTask) (imageToCoordinatesTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
		solverapi.Setting("MinClicks", task.MinClicks != 0),
		solverapi.Setting("MaxClicks", task.MaxClicks != 0),
	); err != nil {
		return imageToCoordinatesTask{}, err
	}

	return imageToCoordinatesTask{
		Type:       "ImageToCoordinatesTask",
		Body:       task.Body,
		Comment:    task.Comment,
		Mode:       string(task.Mode),
		WebsiteURL: task.WebsiteURL,
	}, nil
}
