package unicap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Solution is a solved captcha. Providers fill the typed fields that apply to
// the captcha family and always keep the provider's full payload in Raw, so
// family-specific values remain available through DecodeSolution.
type Solution struct {
	// Token is the response token for token-based captchas such as
	// reCAPTCHA, Turnstile, and most widget captchas.
	Token string

	// Text is the answer for image-to-text, text, and audio captchas.
	Text string

	// UserAgent is the User-Agent the provider's worker solved the captcha
	// with. Send it with the token when the provider reports one.
	UserAgent string

	// Cookies are cookies to send with the follow-up request, such as
	// DataDome's datadome cookie or Cloudflare's cf_clearance.
	Cookies []*http.Cookie

	// Coordinates are the selected points for click and coordinate captchas.
	Coordinates []Coordinate

	// Boxes are the selected rectangles for bounding-box captchas.
	Boxes []Box

	// Indexes are the selected tiles for grid and classification captchas.
	Indexes []int

	// Angle is the rotation in degrees for rotate captchas.
	Angle float64

	// Raw is the provider's solution payload exactly as received.
	Raw json.RawMessage
}

// Coordinate is a point within an image, measured from the top-left corner.
type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Box is a rectangle within an image, measured from the top-left corner.
type Box struct {
	X1 int `json:"x1"`
	Y1 int `json:"y1"`
	X2 int `json:"x2"`
	Y2 int `json:"y2"`
}

// ErrNoSolutionPayload reports that a solution has no raw payload to decode.
var ErrNoSolutionPayload = errors.New("solution has no raw payload")

// Cookie returns the named cookie from the solution.
func (s *Solution) Cookie(name string) (*http.Cookie, bool) {
	for _, cookie := range s.Cookies {
		if cookie.Name == name {
			return cookie, true
		}
	}

	return nil, false
}

// DecodeSolution decodes the provider's raw solution payload into T. Use it
// with the solution types in the tasks package, such as
// tasks.GeeTestV4Solution, to read family-specific fields.
func DecodeSolution[T any](solution Solution) (T, error) {
	var decoded T
	if len(solution.Raw) == 0 {
		return decoded, ErrNoSolutionPayload
	}

	if err := json.Unmarshal(solution.Raw, &decoded); err != nil {
		return decoded, fmt.Errorf("decoding solution: %w", err)
	}

	return decoded, nil
}
