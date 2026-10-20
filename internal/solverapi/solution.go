package solverapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/aarock1234/unicap/v3"
)

// solutionFields are the raw values of a solution object, keyed by field
// name. Each known field is decoded on its own, so one field with an
// unexpected shape never hides the others.
type solutionFields map[string]json.RawMessage

// decodeSolution normalises a provider solution payload into unicap's common
// fields and keeps the payload itself in Raw.
func decodeSolution(raw json.RawMessage) unicap.Solution {
	solution := unicap.Solution{
		Raw: raw,
	}

	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return solution
	}

	switch trimmed[0] {
	case '"':
		_ = json.Unmarshal(trimmed, &solution.Text)

		return solution
	case '{':
		var fields solutionFields
		if err := json.Unmarshal(trimmed, &fields); err != nil {
			return solution
		}

		fields.fill(&solution)

		return solution
	default:
		return solution
	}
}

// fill copies every recognised field into solution.
func (f solutionFields) fill(solution *unicap.Solution) {
	solution.Token = f.firstString("gRecaptchaResponse", "token")
	solution.Text = f.firstString("text")
	solution.UserAgent = f.firstString("userAgent")
	solution.Cookies = f.cookies()
	solution.Coordinates = f.coordinates()
	solution.Boxes = f.boxes()
	solution.Indexes = f.indexes()
	solution.Angle = f.angle()

	// CapSolver returns AWS WAF tokens as a bare "cookie" value.
	if solution.Token == "" && len(solution.Cookies) == 0 {
		if cookie := f.firstString("cookie"); cookie != "" && !strings.Contains(cookie, "=") {
			solution.Token = cookie
		}
	}
}

// decode unmarshals a field into target and reports whether it succeeded.
func (f solutionFields) decode(name string, target any) bool {
	value, ok := f[name]
	if !ok {
		return false
	}

	return json.Unmarshal(value, target) == nil
}

// firstString returns the first non-empty string field among names.
func (f solutionFields) firstString(names ...string) string {
	for _, name := range names {
		var value string
		if f.decode(name, &value) && value != "" {
			return value
		}
	}

	return ""
}

// cookies gathers cookies from "cookie" strings and "cookies" strings, maps,
// or arrays.
func (f solutionFields) cookies() []*http.Cookie {
	var cookies []*http.Cookie

	if text := f.firstString("cookie"); strings.Contains(text, "=") {
		cookies = append(cookies, parseCookieString(text)...)
	}

	var cookieText string
	if f.decode("cookies", &cookieText) {
		return append(cookies, parseCookieString(cookieText)...)
	}

	var cookieMap map[string]string
	if f.decode("cookies", &cookieMap) {
		return append(cookies, cookiesFromMap(cookieMap)...)
	}

	var cookieList []namedCookie
	if f.decode("cookies", &cookieList) {
		for _, cookie := range cookieList {
			cookies = append(cookies, &http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			})
		}
	}

	return cookies
}

// namedCookie is the {name, value} cookie object some providers return.
type namedCookie struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// parseCookieString parses either a Set-Cookie header value, such as
// DataDome's "datadome=...; Max-Age=...; Path=/", or a Cookie header value.
func parseCookieString(text string) []*http.Cookie {
	lower := strings.ToLower(text)
	isSetCookie := strings.Contains(lower, "max-age=") ||
		strings.Contains(lower, "path=") ||
		strings.Contains(lower, "expires=") ||
		strings.Contains(lower, "domain=")

	if isSetCookie {
		cookie, err := http.ParseSetCookie(text)
		if err != nil {
			return nil
		}

		return []*http.Cookie{
			cookie,
		}
	}

	cookies, err := http.ParseCookie(text)
	if err != nil {
		return nil
	}

	return cookies
}

// cookiesFromMap converts a name-to-value map into cookies sorted by name.
func cookiesFromMap(values map[string]string) []*http.Cookie {
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	slices.Sort(names)

	cookies := make([]*http.Cookie, 0, len(names))
	for _, name := range names {
		cookies = append(cookies, &http.Cookie{
			Name:  name,
			Value: values[name],
		})
	}

	return cookies
}

// coordinates reads points from "coordinates" ({x,y} objects or [x,y]
// pairs), "points" ([x,y] pairs), or "box" (a single [x,y] pair).
func (f solutionFields) coordinates() []unicap.Coordinate {
	var objects []unicap.Coordinate
	if f.decode("coordinates", &objects) {
		return objects
	}

	for _, name := range []string{"coordinates", "points"} {
		var pairs [][]float64
		if f.decode(name, &pairs) {
			return pointsFromPairs(pairs)
		}
	}

	var box []float64
	if f.decode("box", &box) && len(box) == 2 {
		return pointsFromPairs([][]float64{
			box,
		})
	}

	return nil
}

// pointsFromPairs converts [x, y] pairs, skipping entries of any other
// length.
func pointsFromPairs(pairs [][]float64) []unicap.Coordinate {
	var points []unicap.Coordinate
	for _, pair := range pairs {
		if len(pair) != 2 {
			continue
		}

		points = append(points, unicap.Coordinate{
			X: int(pair[0]),
			Y: int(pair[1]),
		})
	}

	return points
}

// boxes reads rectangles from "bounding_boxes" (min/max objects), "rects"
// (corner objects), or "coordinates" ([x1,y1,x2,y2] arrays).
func (f solutionFields) boxes() []unicap.Box {
	var extents []boxExtent
	if f.decode("bounding_boxes", &extents) {
		boxes := make([]unicap.Box, 0, len(extents))
		for _, extent := range extents {
			boxes = append(boxes, unicap.Box{
				X1: extent.XMin,
				Y1: extent.YMin,
				X2: extent.XMax,
				Y2: extent.YMax,
			})
		}

		return boxes
	}

	var corners []unicap.Box
	if f.decode("rects", &corners) {
		return corners
	}

	var quads [][]float64
	if !f.decode("coordinates", &quads) {
		return nil
	}

	var boxes []unicap.Box
	for _, quad := range quads {
		if len(quad) != 4 {
			continue
		}

		boxes = append(boxes, unicap.Box{
			X1: int(quad[0]),
			Y1: int(quad[1]),
			X2: int(quad[2]),
			Y2: int(quad[3]),
		})
	}

	return boxes
}

// boxExtent is 2Captcha's bounding-box shape.
type boxExtent struct {
	XMin int `json:"xMin"`
	XMax int `json:"xMax"`
	YMin int `json:"yMin"`
	YMax int `json:"yMax"`
}

// indexes reads selected tiles from "click" or "objects".
func (f solutionFields) indexes() []int {
	for _, name := range []string{"click", "objects"} {
		var values []int
		if f.decode(name, &values) {
			return values
		}
	}

	return nil
}

// angle reads a rotation from "rotate" or "angle".
func (f solutionFields) angle() float64 {
	for _, name := range []string{"rotate", "angle"} {
		var value Number
		if f.decode(name, &value) && value != 0 {
			return float64(value)
		}
	}

	return 0
}
