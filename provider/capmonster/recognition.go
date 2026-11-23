package capmonster

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// imageToTextTask is CapMonster Cloud's image-to-text task. Its numeric
// field takes 0 and 1 only.
type imageToTextTask struct {
	Type             string `json:"type"`
	Body             string `json:"body"`
	CapMonsterModule string `json:"capMonsterModule,omitempty"`
	Case             bool   `json:"case,omitempty"`
	Numeric          int    `json:"numeric,omitempty"`
	Math             bool   `json:"math,omitempty"`
}

// mapImageToText builds CapMonster Cloud's image-to-text task. Of the
// numeric modes, only digits-only has a wire value.
func mapImageToText(task *tasks.ImageToTextTask) (imageToTextTask, error) {
	hasNumericValue := task.Numeric == tasks.NumericModeAny || task.Numeric == tasks.NumericModeNumbersOnly

	if err := solverapi.Unsupported(
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
		solverapi.Setting("Images", len(task.Images) > 0),
		solverapi.Setting("Numeric (letters only, either or both)", !hasNumericValue),
		solverapi.Setting("MinLength", task.MinLength != 0),
		solverapi.Setting("MaxLength", task.MaxLength != 0),
		solverapi.Setting("Phrase", task.Phrase),
		solverapi.Setting("Comment", task.Comment != ""),
		solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
		solverapi.Setting("LanguagePool", task.LanguagePool != ""),
	); err != nil {
		return imageToTextTask{}, err
	}

	return imageToTextTask{
		Type:             "ImageToTextTask",
		Body:             task.Body,
		CapMonsterModule: task.Module,
		Case:             task.Case,
		Numeric:          int(task.Numeric),
		Math:             task.Math,
	}, nil
}
