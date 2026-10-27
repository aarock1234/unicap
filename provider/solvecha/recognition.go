package solvecha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type imageToTextTask struct {
	Type   string `json:"type"`
	Body   string `json:"body"`
	Module string `json:"module,omitempty"`
}

type reCaptchaClassificationTask struct {
	Type     string `json:"type"`
	Image    string `json:"image"`
	Question string `json:"question"`
}

func mapImageToText(task *tasks.ImageToTextTask) (imageToTextTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Images", len(task.Images) > 0),
		solverapi.Setting("Numeric", task.Numeric != tasks.NumericModeAny),
		solverapi.Setting("Math", task.Math),
		solverapi.Setting("MinLength", task.MinLength != 0),
		solverapi.Setting("MaxLength", task.MaxLength != 0),
		solverapi.Setting("Case", task.Case),
		solverapi.Setting("Phrase", task.Phrase),
		solverapi.Setting("Comment", task.Comment != ""),
		solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
		solverapi.Setting("LanguagePool", task.LanguagePool != ""),
	); err != nil {
		return imageToTextTask{}, err
	}

	return imageToTextTask{
		Type:   "ImageToTextTask",
		Body:   task.Body,
		Module: task.Module,
	}, nil
}

func mapReCaptchaClassification(task *tasks.ReCaptchaClassificationTask) (reCaptchaClassificationTask, error) {
	return reCaptchaClassificationTask{
		Type:     "ReCaptchaV2Classification",
		Image:    task.Image,
		Question: task.Question,
	}, nil
}
