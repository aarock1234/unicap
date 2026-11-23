package noptcha

import (
	"encoding/json"
	"maps"
	"net/http"
	"slices"

	"github.com/aarock1234/unicap/v3"
)

// solution holds the fields of Noptcha's solution objects that unicap has a
// common field for. Which ones are present depends on the task type; the
// rest of a solution, such as hCaptcha's resp_key, stays in Solution.Raw.
type solution struct {
	// Token is the response token of the widget captchas.
	Token string `json:"token"`

	// UserAgent is the User-Agent the solver used.
	UserAgent string `json:"user_agent"`

	// Text is the answer of an image.text task.
	Text string `json:"text"`

	// Cookies are the cookies of a passed Cloudflare challenge, as name to
	// value.
	Cookies map[string]string `json:"cookies"`

	// Objects are the zero-based indexes of the matching grid tiles.
	Objects []int `json:"objects"`
}

// decodeSolution converts the solution object of a solved task and keeps the
// object itself in Raw. A field that does not decode stays empty.
func decodeSolution(raw json.RawMessage) unicap.Solution {
	var decoded solution

	_ = json.Unmarshal(raw, &decoded)

	return unicap.Solution{
		Token:     decoded.Token,
		Text:      decoded.Text,
		UserAgent: decoded.UserAgent,
		Cookies:   newCookies(decoded.Cookies),
		Indexes:   decoded.Objects,
		Raw:       raw,
	}
}

// newCookies converts a name-to-value object into cookies sorted by name.
func newCookies(values map[string]string) []*http.Cookie {
	var cookies []*http.Cookie

	for _, name := range slices.Sorted(maps.Keys(values)) {
		cookies = append(cookies, &http.Cookie{
			Name:  name,
			Value: values[name],
		})
	}

	return cookies
}
