package captchaai

import (
	"cmp"
	"encoding/json"
	"net/http"

	"github.com/aarock1234/unicap/v3"
)

// clearanceCookieName is the name of the cookie a passed Cloudflare
// challenge yields. CaptchaAI reports the cookie's value only.
const clearanceCookieName = "cf_clearance"

// solutionShape names how the solution of a task arrives in a reply to
// res.php. It is the label a unicap task ID carries.
type solutionShape string

// The solution shapes. The decoders table is the closed set of them.
const (
	// shapeToken is a token in the request field.
	shapeToken solutionShape = "token"

	// shapeAgentToken is a token in the result field, with the User-Agent
	// it has to be used with.
	shapeAgentToken solutionShape = "agent-token"

	// shapeCookie is the value of a Cloudflare clearance cookie in the
	// result field, with the User-Agent it has to be used with.
	shapeCookie solutionShape = "cookie"

	// shapeGeeTest is the three top-level values of a GeeTest solution,
	// which unicap has no common field for.
	shapeGeeTest solutionShape = "geetest"

	// shapeText is the text of an image in the request field.
	shapeText solutionShape = "text"

	// shapeGrid is the matching tiles of an image grid, as a JSON array
	// inside the request field's string.
	shapeGrid solutionShape = "grid"

	// shapeRaw is the solution of a raw task, whose shape is not known.
	shapeRaw solutionShape = "raw"
)

// decoder fills the common solution fields a shape has a value for.
type decoder func(answer reply) unicap.Solution

// decoders holds the decoder of every solution shape.
var decoders = map[solutionShape]decoder{
	shapeToken:      decodeToken,
	shapeAgentToken: decodeAgentToken,
	shapeCookie:     decodeCookie,
	shapeGeeTest:    decodeGeeTest,
	shapeText:       decodeText,
	shapeGrid:       decodeGrid,
	shapeRaw:        decodeRaw,
}

// decodeSolution converts the reply of a solved task whose solution has the
// given shape, which must be a registered one, and keeps the whole reply
// body in Raw. A field that does not decode stays empty.
func decodeSolution(shape solutionShape, answer reply, body []byte) unicap.Solution {
	solution := decoders[shape](answer)
	solution.Raw = body

	return solution
}

// decodeToken reads a token from the request field.
func decodeToken(answer reply) unicap.Solution {
	return unicap.Solution{
		Token: answer.requestText(),
	}
}

// decodeAgentToken reads a token from the result field and the User-Agent
// the solver used.
func decodeAgentToken(answer reply) unicap.Solution {
	return unicap.Solution{
		Token:     answer.Result,
		UserAgent: answer.UserAgent,
	}
}

// decodeCookie builds the Cloudflare clearance cookie from its value in the
// result field, and reads the User-Agent the solver used. A reply without a
// value yields no cookie.
func decodeCookie(answer reply) unicap.Solution {
	solution := unicap.Solution{
		UserAgent: answer.UserAgent,
	}

	if answer.Result != "" {
		solution.Cookies = []*http.Cookie{
			{
				Name:  clearanceCookieName,
				Value: answer.Result,
			},
		}
	}

	return solution
}

// decodeGeeTest fills no common field: a GeeTest solution is read from Raw,
// which decodes into tasks.GeeTestSolution.
func decodeGeeTest(reply) unicap.Solution {
	return unicap.Solution{}
}

// decodeText reads the text of an image from the request field.
func decodeText(answer reply) unicap.Solution {
	return unicap.Solution{
		Text: answer.requestText(),
	}
}

// decodeGrid reads the matching tiles of an image grid, which the request
// field holds as a JSON array inside a string, such as "[1, 3, 6, 9]". The
// indexes count from 1 in reading order and are kept as reported.
func decodeGrid(answer reply) unicap.Solution {
	var indexes []int

	_ = json.Unmarshal([]byte(answer.requestText()), &indexes)

	return unicap.Solution{
		Indexes: indexes,
	}
}

// decodeRaw reads the token of a raw task from the request field or, failing
// that, the result field. A solution of another shape leaves the token
// empty.
func decodeRaw(answer reply) unicap.Solution {
	return unicap.Solution{
		Token: cmp.Or(answer.requestText(), answer.Result),
	}
}
