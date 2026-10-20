package capsolver

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type imageToTextTask struct {
	Type       string   `json:"type"`
	Body       string   `json:"body"`
	WebsiteURL string   `json:"websiteURL,omitempty"`
	Module     string   `json:"module,omitempty"`
	Images     []string `json:"images,omitempty"`
}

type awsWAFClassificationTask struct {
	Type       string   `json:"type"`
	Images     []string `json:"images"`
	Question   string   `json:"question"`
	WebsiteURL string   `json:"websiteURL,omitempty"`
}

type visionEngineTask struct {
	Type            string `json:"type"`
	Module          string `json:"module"`
	Image           string `json:"image"`
	ImageBackground string `json:"imageBackground,omitempty"`
	Question        string `json:"question,omitempty"`
	WebsiteURL      string `json:"websiteURL,omitempty"`
}

func mapImageToText(task *tasks.ImageToTextTask) (imageToTextTask, error) {
	if err := solverapi.Unsupported(
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
		Type:       "ImageToTextTask",
		Body:       task.Body,
		WebsiteURL: task.WebsiteURL,
		Module:     task.Module,
		Images:     task.Images,
	}, nil
}

func mapAWSWAFClassification(task *tasks.AWSWAFClassificationTask) (awsWAFClassificationTask, error) {
	return awsWAFClassificationTask{
		Type:       "AwsWafClassification",
		Images:     task.Images,
		Question:   task.Question,
		WebsiteURL: task.WebsiteURL,
	}, nil
}

func mapVisionEngine(task *tasks.VisionEngineTask) (visionEngineTask, error) {
	return visionEngineTask{
		Type:            "VisionEngine",
		Module:          string(task.Module),
		Image:           task.Image,
		ImageBackground: task.ImageBackground,
		Question:        task.Question,
		WebsiteURL:      task.WebsiteURL,
	}, nil
}
