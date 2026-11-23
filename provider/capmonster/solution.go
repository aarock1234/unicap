package capmonster

import (
	"cmp"
	"encoding/json"
	"net/http"
	"slices"

	"github.com/aarock1234/unicap/v3"
)

// solutionExtras are the parts of a CapMonster Cloud solution the shared
// decoder does not read.
type solutionExtras struct {
	// Clearance is the cf_clearance cookie value of a Cloudflare challenge
	// solution.
	Clearance string `json:"cf_clearance"`

	// Data is the object CustomTask solutions nest their values in.
	Data solutionData `json:"data"`

	// Headers are the request headers a CustomTask solution was obtained
	// with.
	Headers map[string]string `json:"headers"`
}

// solutionData is the data object of a CustomTask solution.
type solutionData struct {
	Token string `json:"token"`
}

// placeholderCookie is the name of the entry in the {"nocookies":"true"}
// cookies object CapMonster Cloud returns with reCAPTCHA solutions that have
// no cookies.
const placeholderCookie = "nocookies"

// clearanceCookie is the name of the cookie a passed Cloudflare challenge
// sets.
const clearanceCookie = "cf_clearance"

// wafTokenCookie is the name of the cookie that carries an AWS WAF token.
const wafTokenCookie = "aws-waf-token"

// userAgentHeader is the key of the User-Agent in a solution's headers
// object.
const userAgentHeader = "User-Agent"

// refineSolution corrects and completes the solution the shared decoder
// built from a CapMonster Cloud payload.
func refineSolution(raw json.RawMessage, solution *unicap.Solution) {
	extras := decodeSolutionExtras(raw)

	dropPlaceholderCookie(solution)
	addClearanceCookie(solution, extras.Clearance)
	fillToken(solution, extras.Data.Token)
	fillUserAgent(solution, extras.Headers[userAgentHeader])
}

// decodeSolutionExtras reads the extras of a solution payload. Fields that
// are absent or do not decode stay empty.
func decodeSolutionExtras(raw json.RawMessage) solutionExtras {
	var extras solutionExtras

	_ = json.Unmarshal(raw, &extras)

	return extras
}

// dropPlaceholderCookie removes the cookie the shared decoder made of the
// {"nocookies":"true"} placeholder.
func dropPlaceholderCookie(solution *unicap.Solution) {
	solution.Cookies = slices.DeleteFunc(solution.Cookies, func(cookie *http.Cookie) bool {
		return cookie.Name == placeholderCookie
	})
}

// addClearanceCookie appends the cf_clearance cookie of a Cloudflare
// challenge solution, which carries only the cookie's value.
func addClearanceCookie(solution *unicap.Solution, clearance string) {
	if clearance == "" {
		return
	}

	solution.Cookies = append(solution.Cookies, &http.Cookie{
		Name:  clearanceCookie,
		Value: clearance,
	})
}

// fillToken sets a missing token from the token nested in the solution's
// data object, or else from the aws-waf-token cookie.
func fillToken(solution *unicap.Solution, dataToken string) {
	solution.Token = cmp.Or(
		solution.Token,
		dataToken,
		cookieValue(solution.Cookies, wafTokenCookie),
	)
}

// fillUserAgent sets a missing user agent from the solution's headers
// object.
func fillUserAgent(solution *unicap.Solution, headerUserAgent string) {
	solution.UserAgent = cmp.Or(solution.UserAgent, headerUserAgent)
}

// cookieValue returns the value of the first cookie called name, or an empty
// string when there is none.
func cookieValue(cookies []*http.Cookie, name string) string {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}

	return ""
}
