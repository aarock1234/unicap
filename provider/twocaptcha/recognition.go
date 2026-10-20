package twocaptcha

import (
	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

type imageToTextTask struct {
	Type            string `json:"type"`
	Body            string `json:"body"`
	Phrase          bool   `json:"phrase,omitempty"`
	Case            bool   `json:"case,omitempty"`
	Numeric         int    `json:"numeric,omitempty"`
	Math            bool   `json:"math,omitempty"`
	MinLength       int    `json:"minLength,omitempty"`
	MaxLength       int    `json:"maxLength,omitempty"`
	Comment         string `json:"comment,omitempty"`
	ImgInstructions string `json:"imgInstructions,omitempty"`
}

type textCaptchaTask struct {
	Type    string `json:"type"`
	Comment string `json:"comment"`
}

type audioTask struct {
	Type string `json:"type"`
	Body string `json:"body"`
	Lang string `json:"lang"`
}

// clickTask is the shape shared by image tasks answered with clicks,
// outlines, or boxes.
type clickTask struct {
	Type            string `json:"type"`
	Body            string `json:"body"`
	Comment         string `json:"comment,omitempty"`
	ImgInstructions string `json:"imgInstructions,omitempty"`
	MinClicks       int    `json:"minClicks,omitempty"`
	MaxClicks       int    `json:"maxClicks,omitempty"`
	CanNoAnswer     int    `json:"canNoAnswer,omitempty"`
}

type gridTask struct {
	Type            string `json:"type"`
	Body            string `json:"body"`
	Rows            int    `json:"rows,omitempty"`
	Columns         int    `json:"columns,omitempty"`
	Comment         string `json:"comment,omitempty"`
	ImgInstructions string `json:"imgInstructions,omitempty"`
	MinClicks       int    `json:"minClicks,omitempty"`
	MaxClicks       int    `json:"maxClicks,omitempty"`
	CanNoAnswer     int    `json:"canNoAnswer,omitempty"`
	PreviousID      string `json:"previousId,omitempty"`
	ImgType         string `json:"imgType,omitempty"`
}

type rotateTask struct {
	Type            string `json:"type"`
	Body            string `json:"body"`
	Angle           int    `json:"angle,omitempty"`
	Comment         string `json:"comment,omitempty"`
	ImgInstructions string `json:"imgInstructions,omitempty"`
}

type dragAndDropTask struct {
	Type       string   `json:"type"`
	Background string   `json:"background"`
	Images     []string `json:"images"`
	Comment    string   `json:"comment,omitempty"`
}

type temuTask struct {
	Type  string   `json:"type"`
	Image string   `json:"image"`
	Parts []string `json:"parts"`
}

type vkCaptchaImageTask struct {
	Type  string `json:"type"`
	Image string `json:"image"`
	Steps []int  `json:"steps"`
}

type yandexImageTask struct {
	Type            string `json:"type"`
	Image           string `json:"image"`
	ImgInstructions string `json:"imgInstructions"`
	Comment         string `json:"comment,omitempty"`
}

type yandexPuzzleTask struct {
	Type  string `json:"type"`
	Image string `json:"image"`
	Task  string `json:"task"`
}

func mapImageToText(task *tasks.ImageToTextTask) (imageToTextTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
		solverapi.Setting("Module", task.Module != ""),
		solverapi.Setting("Images", len(task.Images) > 0),
		solverapi.Setting("LanguagePool (use twocaptcha.WithLanguagePool)", task.LanguagePool != ""),
	); err != nil {
		return imageToTextTask{}, err
	}

	return imageToTextTask{
		Type:            "ImageToTextTask",
		Body:            task.Body,
		Phrase:          task.Phrase,
		Case:            task.Case,
		Numeric:         int(task.Numeric),
		Math:            task.Math,
		MinLength:       task.MinLength,
		MaxLength:       task.MaxLength,
		Comment:         task.Comment,
		ImgInstructions: task.ImgInstructions,
	}, nil
}

func mapText(task *tasks.TextCaptchaTask) (textCaptchaTask, error) {
	return textCaptchaTask{
		Type:    "TextCaptchaTask",
		Comment: task.Question,
	}, nil
}

func mapAudio(task *tasks.AudioTask) (audioTask, error) {
	return audioTask{
		Type: "AudioTask",
		Body: task.Body,
		Lang: task.Language,
	}, nil
}

func mapCoordinates(task *tasks.CoordinatesTask) (clickTask, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("Mode", task.Mode == tasks.CoordinatesModeRectangles),
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
	); err != nil {
		return clickTask{}, err
	}

	return clickTask{
		Type:            "CoordinatesTask",
		Body:            task.Body,
		Comment:         task.Comment,
		ImgInstructions: task.ImgInstructions,
		MinClicks:       task.MinClicks,
		MaxClicks:       task.MaxClicks,
	}, nil
}

func mapGrid(task *tasks.GridTask) (gridTask, error) {
	return gridTask{
		Type:            "GridTask",
		Body:            task.Body,
		Rows:            task.Rows,
		Columns:         task.Columns,
		Comment:         task.Comment,
		ImgInstructions: task.ImgInstructions,
		MinClicks:       task.MinClicks,
		MaxClicks:       task.MaxClicks,
		CanNoAnswer:     boolFlag(task.CanSkip),
		PreviousID:      task.PreviousID,
		ImgType:         string(task.ImageType),
	}, nil
}

func mapRotate(task *tasks.RotateTask) (rotateTask, error) {
	return rotateTask{
		Type:            "RotateTask",
		Body:            task.Body,
		Angle:           task.Angle,
		Comment:         task.Comment,
		ImgInstructions: task.ImgInstructions,
	}, nil
}

func mapDrawAround(task *tasks.DrawAroundTask) (clickTask, error) {
	return clickTask{
		Type:            "DrawAroundTask",
		Body:            task.Body,
		Comment:         task.Comment,
		ImgInstructions: task.ImgInstructions,
	}, nil
}

func mapBoundingBox(task *tasks.BoundingBoxTask) (clickTask, error) {
	return clickTask{
		Type:            "BoundingBoxTask",
		Body:            task.Body,
		Comment:         task.Comment,
		ImgInstructions: task.ImgInstructions,
		CanNoAnswer:     boolFlag(task.CanSkip),
	}, nil
}

func mapDragAndDrop(task *tasks.DragAndDropTask) (dragAndDropTask, error) {
	return dragAndDropTask{
		Type:       "DragAndDropTask",
		Background: task.Background,
		Images:     task.Images,
		Comment:    task.Comment,
	}, nil
}

func mapTemu(task *tasks.TemuTask) (temuTask, error) {
	return temuTask{
		Type:  "TemuImageTask",
		Image: task.Image,
		Parts: task.Parts,
	}, nil
}

func mapVKCaptchaImage(task *tasks.VKCaptchaImageTask) (vkCaptchaImageTask, error) {
	return vkCaptchaImageTask{
		Type:  "VKCaptchaImageTask",
		Image: task.Image,
		Steps: task.Steps,
	}, nil
}

func mapYandexImage(task *tasks.YandexImageTask) (yandexImageTask, error) {
	return yandexImageTask{
		Type:            "SmartCaptchaTask",
		Image:           task.Image,
		ImgInstructions: task.ImgInstructions,
		Comment:         task.Comment,
	}, nil
}

func mapYandexPuzzle(task *tasks.YandexPuzzleTask) (yandexPuzzleTask, error) {
	return yandexPuzzleTask{
		Type:  "PazlCaptchaTask",
		Image: task.Image,
		Task:  task.Task,
	}, nil
}

// boolFlag converts a boolean to 2Captcha's 0/1 flag.
func boolFlag(value bool) int {
	if value {
		return 1
	}

	return 0
}
