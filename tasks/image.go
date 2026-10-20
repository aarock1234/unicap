package tasks

import (
	"github.com/aarock1234/unicap/v3"
)

// NumericMode constrains the characters in an image-to-text answer. Its
// values are the wire codes the provider APIs expect.
type NumericMode int

const (
	// NumericModeAny applies no character constraint.
	NumericModeAny NumericMode = 0

	// NumericModeNumbersOnly allows digits only.
	NumericModeNumbersOnly NumericMode = 1

	// NumericModeLettersOnly allows letters only.
	NumericModeLettersOnly NumericMode = 2

	// NumericModeEither allows digits only or letters only.
	NumericModeEither NumericMode = 3

	// NumericModeBoth requires both digits and letters.
	NumericModeBoth NumericMode = 4
)

// ImageToTextTask reads the text in a captcha image. The answer is in
// unicap.Solution.Text.
type ImageToTextTask struct {
	// Body is the base64-encoded image.
	Body string

	// WebsiteURL is the optional page the image came from.
	WebsiteURL string

	// Module selects a provider recognition model, such as CapSolver's
	// "common" or "number".
	Module string

	// Images are extra base64-encoded images for multi-image modules.
	Images []string

	// Numeric constrains the answer's characters.
	Numeric NumericMode

	// Math marks an arithmetic captcha whose answer is the result.
	Math bool

	// MinLength is the minimum answer length. Zero means no limit.
	MinLength int

	// MaxLength is the maximum answer length. Zero means no limit.
	MaxLength int

	// Case marks a case-sensitive answer.
	Case bool

	// Phrase marks an answer that contains spaces.
	Phrase bool

	// Comment is an instruction for the worker.
	Comment string

	// ImgInstructions is a base64-encoded image with instructions for the
	// worker.
	ImgInstructions string

	// LanguagePool selects the worker language pool, such as "en" or "rn".
	LanguagePool string
}

// Type returns unicap.TaskTypeImageToText.
func (t *ImageToTextTask) Type() unicap.TaskType {
	return unicap.TaskTypeImageToText
}

// Validate reports a missing image.
func (t *ImageToTextTask) Validate() error {
	return validate(
		required("Body", t.Body),
	)
}

// TextCaptchaTask answers a natural-language question. The answer is in
// unicap.Solution.Text.
type TextCaptchaTask struct {
	// Question is the question to answer.
	Question string
}

// Type returns unicap.TaskTypeText.
func (t *TextCaptchaTask) Type() unicap.TaskType {
	return unicap.TaskTypeText
}

// Validate reports a missing question.
func (t *TextCaptchaTask) Validate() error {
	return validate(
		required("Question", t.Question),
	)
}

// AudioTask transcribes an audio captcha. The answer is in
// unicap.Solution.Text.
type AudioTask struct {
	// Body is the base64-encoded MP3 audio.
	Body string

	// Language is the spoken language, such as "en", "fr", or "de".
	Language string
}

// Type returns unicap.TaskTypeAudio.
func (t *AudioTask) Type() unicap.TaskType {
	return unicap.TaskTypeAudio
}

// Validate reports missing required fields.
func (t *AudioTask) Validate() error {
	return validate(
		required("Body", t.Body),
		required("Language", t.Language),
	)
}

// CoordinatesMode selects whether a coordinates task returns points or
// rectangles.
type CoordinatesMode string

const (
	// CoordinatesModePoints returns clicked points in
	// unicap.Solution.Coordinates.
	CoordinatesModePoints CoordinatesMode = "points"

	// CoordinatesModeRectangles returns selected rectangles in
	// unicap.Solution.Boxes.
	CoordinatesModeRectangles CoordinatesMode = "rectangles"
)

// CoordinatesTask asks a worker to click points, or select rectangles, on an
// image.
type CoordinatesTask struct {
	// Body is the base64-encoded image.
	Body string

	// Mode selects points or rectangles. Empty means points.
	Mode CoordinatesMode

	// Comment is an instruction for the worker.
	Comment string

	// ImgInstructions is a base64-encoded image with instructions for the
	// worker.
	ImgInstructions string

	// MinClicks is the minimum number of points. Zero means no limit.
	MinClicks int

	// MaxClicks is the maximum number of points. Zero means no limit.
	MaxClicks int

	// WebsiteURL is the optional page the image came from.
	WebsiteURL string
}

// Type returns unicap.TaskTypeCoordinates.
func (t *CoordinatesTask) Type() unicap.TaskType {
	return unicap.TaskTypeCoordinates
}

// Validate reports a missing image.
func (t *CoordinatesTask) Validate() error {
	return validate(
		required("Body", t.Body),
	)
}

// GridImageType tells the provider which widget a grid image came from.
type GridImageType string

const (
	// GridImageTypeReCaptcha marks a reCAPTCHA image grid.
	GridImageTypeReCaptcha GridImageType = "recaptcha"

	// GridImageTypeFunCaptcha marks a FunCaptcha image grid.
	GridImageTypeFunCaptcha GridImageType = "funcaptcha"
)

// GridTask asks a worker to select the matching tiles of an image grid. The
// selected tiles are in unicap.Solution.Indexes.
type GridTask struct {
	// Body is the base64-encoded grid image.
	Body string

	// Rows is the number of grid rows.
	Rows int

	// Columns is the number of grid columns.
	Columns int

	// Comment is an instruction for the worker. Set Comment or
	// ImgInstructions.
	Comment string

	// ImgInstructions is a base64-encoded image with instructions for the
	// worker.
	ImgInstructions string

	// MinClicks is the minimum number of tiles. Zero means no limit.
	MinClicks int

	// MaxClicks is the maximum number of tiles. Zero means no limit.
	MaxClicks int

	// CanSkip allows the worker to answer that no tile matches.
	CanSkip bool

	// PreviousID is the task ID of the previous round for multi-round grids.
	PreviousID string

	// ImageType names the widget the grid came from, enabling specialised
	// workers.
	ImageType GridImageType
}

// Type returns unicap.TaskTypeGrid.
func (t *GridTask) Type() unicap.TaskType {
	return unicap.TaskTypeGrid
}

// Validate reports a missing image or instruction.
func (t *GridTask) Validate() error {
	return validate(
		required("Body", t.Body),
		requiredInstruction(t.Comment, t.ImgInstructions),
	)
}

// RotateTask asks a worker how far to rotate an image upright. The answer is
// in unicap.Solution.Angle.
type RotateTask struct {
	// Body is the base64-encoded image.
	Body string

	// Angle is the rotation step in degrees.
	Angle int

	// Comment is an instruction for the worker.
	Comment string

	// ImgInstructions is a base64-encoded image with instructions for the
	// worker.
	ImgInstructions string
}

// Type returns unicap.TaskTypeRotate.
func (t *RotateTask) Type() unicap.TaskType {
	return unicap.TaskTypeRotate
}

// Validate reports a missing image.
func (t *RotateTask) Validate() error {
	return validate(
		required("Body", t.Body),
	)
}

// DrawAroundTask asks a worker to outline objects on an image. Decode the
// outlines with unicap.DecodeSolution[DrawAroundSolution].
type DrawAroundTask struct {
	// Body is the base64-encoded image.
	Body string

	// Comment is an instruction for the worker. Set Comment or
	// ImgInstructions.
	Comment string

	// ImgInstructions is a base64-encoded image with instructions for the
	// worker.
	ImgInstructions string
}

// Type returns unicap.TaskTypeDrawAround.
func (t *DrawAroundTask) Type() unicap.TaskType {
	return unicap.TaskTypeDrawAround
}

// Validate reports a missing image or instruction.
func (t *DrawAroundTask) Validate() error {
	return validate(
		required("Body", t.Body),
		requiredInstruction(t.Comment, t.ImgInstructions),
	)
}

// DrawAroundSolution is the decoded payload of a DrawAroundTask.
type DrawAroundSolution struct {
	// Canvas holds one outline per object.
	Canvas [][]unicap.Coordinate `json:"canvas"`
}

// BoundingBoxTask asks a worker to draw boxes around objects on an image.
// The boxes are in unicap.Solution.Boxes.
type BoundingBoxTask struct {
	// Body is the base64-encoded image.
	Body string

	// Comment is an instruction for the worker. Set Comment or
	// ImgInstructions.
	Comment string

	// ImgInstructions is a base64-encoded image with instructions for the
	// worker.
	ImgInstructions string

	// CanSkip allows the worker to answer that no object is present.
	CanSkip bool
}

// Type returns unicap.TaskTypeBoundingBox.
func (t *BoundingBoxTask) Type() unicap.TaskType {
	return unicap.TaskTypeBoundingBox
}

// Validate reports a missing image or instruction.
func (t *BoundingBoxTask) Validate() error {
	return validate(
		required("Body", t.Body),
		requiredInstruction(t.Comment, t.ImgInstructions),
	)
}

// DragAndDropTask asks a worker where to drop each piece onto a background.
// The drop points are in unicap.Solution.Coordinates, in Images order.
type DragAndDropTask struct {
	// Background is the base64-encoded background image.
	Background string

	// Images are the base64-encoded pieces to place.
	Images []string

	// Comment is an instruction for the worker.
	Comment string
}

// Type returns unicap.TaskTypeDragAndDrop.
func (t *DragAndDropTask) Type() unicap.TaskType {
	return unicap.TaskTypeDragAndDrop
}

// Validate reports missing required fields.
func (t *DragAndDropTask) Validate() error {
	return validate(
		required("Background", t.Background),
		requiredItems("Images", t.Images),
	)
}

// TemuTask solves a Temu puzzle. The piece positions are in
// unicap.Solution.Coordinates.
type TemuTask struct {
	// Image is the base64-encoded background image.
	Image string

	// Parts are the base64-encoded puzzle pieces.
	Parts []string
}

// Type returns unicap.TaskTypeTemu.
func (t *TemuTask) Type() unicap.TaskType {
	return unicap.TaskTypeTemu
}

// Validate reports missing required fields.
func (t *TemuTask) Validate() error {
	return validate(
		required("Image", t.Image),
		requiredItems("Parts", t.Parts),
	)
}

// VKCaptchaImageTask solves a VK slider image captcha. Decode the result
// with unicap.DecodeSolution[VKCaptchaImageSolution].
type VKCaptchaImageTask struct {
	// Image is the base64-encoded captcha image.
	Image string

	// Steps are the slider steps from the captcha configuration.
	Steps []int
}

// Type returns unicap.TaskTypeVKCaptchaImage.
func (t *VKCaptchaImageTask) Type() unicap.TaskType {
	return unicap.TaskTypeVKCaptchaImage
}

// Validate reports missing required fields.
func (t *VKCaptchaImageTask) Validate() error {
	return validate(
		required("Image", t.Image),
		requiredItems("Steps", t.Steps),
	)
}

// VKCaptchaImageSolution is the decoded payload of a VKCaptchaImageTask.
type VKCaptchaImageSolution struct {
	BestStep int    `json:"best_step"`
	Preview  string `json:"preview"`
	Solution []int  `json:"solution"`
	Answer   string `json:"answer"`
}

// YandexImageTask solves a Yandex SmartCaptcha image challenge. The clicks
// are in unicap.Solution.Coordinates.
type YandexImageTask struct {
	// Image is the base64-encoded challenge image.
	Image string

	// ImgInstructions is the base64-encoded instruction image.
	ImgInstructions string

	// Comment is an instruction for the worker.
	Comment string
}

// Type returns unicap.TaskTypeYandexImage.
func (t *YandexImageTask) Type() unicap.TaskType {
	return unicap.TaskTypeYandexImage
}

// Validate reports missing required fields.
func (t *YandexImageTask) Validate() error {
	return validate(
		required("Image", t.Image),
		required("ImgInstructions", t.ImgInstructions),
	)
}

// YandexPuzzleTask solves a Yandex SmartCaptcha puzzle. The answer is in
// unicap.Solution.Text.
type YandexPuzzleTask struct {
	// Image is the base64-encoded puzzle image.
	Image string

	// Task is the puzzle description as a JSON-encoded array.
	Task string
}

// Type returns unicap.TaskTypeYandexPuzzle.
func (t *YandexPuzzleTask) Type() unicap.TaskType {
	return unicap.TaskTypeYandexPuzzle
}

// Validate reports missing required fields.
func (t *YandexPuzzleTask) Validate() error {
	return validate(
		required("Image", t.Image),
		required("Task", t.Task),
	)
}

// VisionEngineModule selects a CapSolver VisionEngine model.
type VisionEngineModule string

const (
	// VisionEngineSlider finds a slider distance.
	VisionEngineSlider VisionEngineModule = "slider_1"

	// VisionEngineRotate finds a rotation angle.
	VisionEngineRotate VisionEngineModule = "rotate_1"

	// VisionEngineRotateSingle finds a rotation angle from a single image.
	VisionEngineRotateSingle VisionEngineModule = "rotate_2"

	// VisionEngineShein finds the Shein targets.
	VisionEngineShein VisionEngineModule = "shein"

	// VisionEngineOCRGIF reads text from an animated GIF.
	VisionEngineOCRGIF VisionEngineModule = "ocr_gif"

	// VisionEngineBotDeflector finds the BotDeflector click points.
	VisionEngineBotDeflector VisionEngineModule = "botdeflector"
)

// VisionEngineTask runs a CapSolver VisionEngine recognition model. Decode
// module-specific values with unicap.DecodeSolution[VisionEngineSolution].
type VisionEngineTask struct {
	// Module is the recognition model.
	Module VisionEngineModule

	// Image is the base64-encoded foreground image.
	Image string

	// ImageBackground is the base64-encoded background image, for modules
	// that need one.
	ImageBackground string

	// Question is the prompt, for modules that need one.
	Question string

	// WebsiteURL is the optional page the image came from.
	WebsiteURL string
}

// Type returns unicap.TaskTypeVisionEngine.
func (t *VisionEngineTask) Type() unicap.TaskType {
	return unicap.TaskTypeVisionEngine
}

// Validate reports missing required fields.
func (t *VisionEngineTask) Validate() error {
	return validate(
		required("Module", string(t.Module)),
		required("Image", t.Image),
	)
}

// VisionEngineSolution is the decoded payload of a VisionEngineTask.
type VisionEngineSolution struct {
	Distance float64      `json:"distance"`
	Angle    float64      `json:"angle"`
	Text     string       `json:"text"`
	Points   [][2]float64 `json:"points"`
	Rects    []unicap.Box `json:"rects"`
}
