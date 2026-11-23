package captchaai

import (
	"fmt"
	"slices"
	"strconv"
	"unicode/utf8"

	"github.com/aarock1234/unicap/v3/internal/solverapi"
	"github.com/aarock1234/unicap/v3/tasks"
)

// imageTypeReCaptcha is the only widget CaptchaAI solves image grids of.
const imageTypeReCaptcha = "recaptcha"

// maxInstructionsLength is the longest instruction, in characters, that
// CaptchaAI takes with an image grid.
const maxInstructionsLength = 140

// gridSides are the sizes of the square grids CaptchaAI solves: 3x3 and 4x4.
var gridSides = []int{
	3,
	4,
}

// mapImageToText builds the form of CaptchaAI's Normal Captcha task for a
// base64-encoded image. CaptchaAI documents fields for a module, a numeric
// mode and case sensitivity only.
func mapImageToText(task *tasks.ImageToTextTask) (form, error) {
	if err := solverapi.Unsupported(
		solverapi.Setting("WebsiteURL", task.WebsiteURL != ""),
		solverapi.Setting("Images", len(task.Images) > 0),
		solverapi.Setting("Math", task.Math),
		solverapi.Setting("MinLength", task.MinLength != 0),
		solverapi.Setting("MaxLength", task.MaxLength != 0),
		solverapi.Setting("Phrase", task.Phrase),
		solverapi.Setting("Comment", task.Comment != ""),
		solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
		solverapi.Setting("LanguagePool", task.LanguagePool != ""),
	); err != nil {
		return nil, err
	}

	return newForm(methodBase64).
		with("body", task.Body).
		optional("module", task.Module).
		optional("numeric", numericCode(task.Numeric)).
		flag("regsense", task.Case), nil
}

// mapGrid builds the form of CaptchaAI's Grid Image task for a
// base64-encoded image. CaptchaAI solves reCAPTCHA grids of 3x3 and 4x4
// tiles only, and needs a written instruction of limited length.
func mapGrid(task *tasks.GridTask) (form, error) {
	if err := solverapi.Check(
		solverapi.Required(
			solverapi.Setting("Comment", task.Comment != ""),
		),
		solverapi.Unsupported(
			solverapi.Setting("a grid other than 3x3 or 4x4", !isSolvableGrid(task.Rows, task.Columns)),
			solverapi.Setting(fmt.Sprintf("Comment over %d characters", maxInstructionsLength), isTooLong(task.Comment)),
			solverapi.Setting("ImgInstructions", task.ImgInstructions != ""),
			solverapi.Setting("MinClicks", task.MinClicks != 0),
			solverapi.Setting("MaxClicks", task.MaxClicks != 0),
			solverapi.Setting("CanSkip", task.CanSkip),
			solverapi.Setting("PreviousID", task.PreviousID != ""),
			solverapi.Setting("ImageType "+string(task.ImageType), !isSolvableGridImage(task.ImageType)),
		),
	); err != nil {
		return nil, err
	}

	return newForm(methodBase64).
		with("body", task.Body).
		with("img_type", imageTypeReCaptcha).
		with("instructions", task.Comment).
		with("grid_size", gridSize(task.Rows, task.Columns)), nil
}

// numericCode returns the value of the numeric field for a numeric mode,
// whose codes are CaptchaAI's own. No constraint yields an empty value,
// which leaves the field out.
func numericCode(mode tasks.NumericMode) string {
	if mode == tasks.NumericModeAny {
		return ""
	}

	return strconv.Itoa(int(mode))
}

// isSolvableGrid reports whether CaptchaAI solves a grid of the given
// dimensions.
func isSolvableGrid(rows, columns int) bool {
	return rows == columns && slices.Contains(gridSides, rows)
}

// isTooLong reports whether an instruction exceeds the length CaptchaAI
// takes.
func isTooLong(instructions string) bool {
	return utf8.RuneCountInString(instructions) > maxInstructionsLength
}

// isSolvableGridImage reports whether CaptchaAI solves grids of the given
// widget: reCAPTCHA's, which an unspecified widget is taken to be.
func isSolvableGridImage(imageType tasks.GridImageType) bool {
	return imageType == "" || imageType == tasks.GridImageTypeReCaptcha
}

// gridSize returns the value of the grid_size field, such as "3x3".
func gridSize(rows, columns int) string {
	return strconv.Itoa(rows) + "x" + strconv.Itoa(columns)
}
