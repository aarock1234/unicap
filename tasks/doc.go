// Package tasks provides the captcha task types unicap models, one per
// captcha family, together with typed solution payloads for families whose
// answers do not fit unicap.Solution's common fields.
//
// Every task validates the inputs all providers require. Providers reject
// options they cannot honour with unicap.ErrUnsupportedOption rather than
// silently dropping them.
package tasks
