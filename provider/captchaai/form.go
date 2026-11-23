package captchaai

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/aarock1234/unicap/v3/internal/httpapi"
)

// field is one field of a form.
type field struct {
	name  string
	value string
}

// form holds the fields of a request in the order they are sent.
type form []field

// newForm starts the form of a task that the given method solves.
func newForm(method string) form {
	return form{}.with(fieldMethod, method)
}

// with adds a field.
func (f form) with(name, value string) form {
	return append(f, field{
		name:  name,
		value: value,
	})
}

// optional adds a field unless its value is empty, which leaves CaptchaAI's
// default in place.
func (f form) optional(name, value string) form {
	if value == "" {
		return f
	}

	return f.with(name, value)
}

// flag adds a field that switches something on, and nothing when it is off.
func (f form) flag(name string, set bool) form {
	if !set {
		return f
	}

	return f.with(name, flagSet)
}

// extend adds the fields of another form.
func (f form) extend(other form) form {
	return append(f, other...)
}

// encodeForm encodes fields, in order, as a multipart/form-data body.
func encodeForm(fields form) (*httpapi.EncodedBody, error) {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	for _, entry := range fields {
		if err := writer.WriteField(entry.name, entry.value); err != nil {
			return nil, fmt.Errorf("writing form field %s: %w", entry.name, err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing form: %w", err)
	}

	return &httpapi.EncodedBody{
		ContentType: writer.FormDataContentType(),
		Data:        body.Bytes(),
	}, nil
}
