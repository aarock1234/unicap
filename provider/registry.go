// Package provider selects built-in or custom unicap providers by name, for
// applications that choose a provider from configuration at runtime.
package provider

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/provider/anticaptcha"
	"github.com/aarock1234/unicap/v3/provider/capsolver"
	"github.com/aarock1234/unicap/v3/provider/twocaptcha"
)

// ErrUnknownProvider reports that no factory is registered for a provider
// name.
var ErrUnknownProvider = errors.New("unknown provider")

// Settings are the provider-independent options a Factory receives. Zero
// fields keep the provider's defaults.
type Settings struct {
	// HTTPClient is the HTTP client for API requests.
	HTTPClient *http.Client

	// Logger receives request diagnostics.
	Logger *slog.Logger

	// BaseURL overrides the provider's API root.
	BaseURL string

	// CallbackURL is the URL the provider notifies when a task finishes.
	CallbackURL string
}

// Option configures the Settings passed to a Factory.
type Option func(*Settings)

// WithHTTPClient sets the HTTP client for API requests.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(s *Settings) {
		s.HTTPClient = httpClient
	}
}

// WithLogger sets the logger for request diagnostics.
func WithLogger(logger *slog.Logger) Option {
	return func(s *Settings) {
		s.Logger = logger
	}
}

// WithBaseURL overrides the provider's API root.
func WithBaseURL(baseURL string) Option {
	return func(s *Settings) {
		s.BaseURL = baseURL
	}
}

// WithCallbackURL sets the URL the provider notifies when a task finishes.
func WithCallbackURL(callbackURL string) Option {
	return func(s *Settings) {
		s.CallbackURL = callbackURL
	}
}

// Factory creates a provider from an API key and settings.
type Factory func(apiKey string, settings Settings) (unicap.Provider, error)

// Registry maps provider names to factories. Its zero value is an empty
// registry ready for use, and it is safe for concurrent use.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// NewRegistry creates a registry preloaded with the built-in providers.
func NewRegistry() *Registry {
	r := &Registry{}

	r.Register(capsolver.Name, newCapSolver)
	r.Register(twocaptcha.Name, newTwoCaptcha)
	r.Register(anticaptcha.Name, newAntiCaptcha)

	return r
}

// Register adds or replaces the factory for name.
func (r *Registry) Register(name string, factory Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.factories == nil {
		r.factories = make(map[string]Factory)
	}

	r.factories[name] = factory
}

// Get returns the factory registered for name.
func (r *Registry) Get(name string) (Factory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, ok := r.factories[name]

	return factory, ok
}

// Names returns the registered provider names in sorted order.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.factories))
	for name := range r.factories {
		names = append(names, name)
	}
	slices.Sort(names)

	return names
}

// New creates the provider registered as name.
func (r *Registry) New(name, apiKey string, opts ...Option) (unicap.Provider, error) {
	factory, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("%q: %w", name, ErrUnknownProvider)
	}

	var settings Settings
	for _, opt := range opts {
		opt(&settings)
	}

	return factory(apiKey, settings)
}

func newCapSolver(apiKey string, settings Settings) (unicap.Provider, error) {
	client, err := capsolver.New(apiKey,
		capsolver.WithHTTPClient(settings.HTTPClient),
		capsolver.WithLogger(settings.Logger),
		capsolver.WithBaseURL(settings.BaseURL),
		capsolver.WithCallbackURL(settings.CallbackURL),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newTwoCaptcha(apiKey string, settings Settings) (unicap.Provider, error) {
	client, err := twocaptcha.New(apiKey,
		twocaptcha.WithHTTPClient(settings.HTTPClient),
		twocaptcha.WithLogger(settings.Logger),
		twocaptcha.WithBaseURL(settings.BaseURL),
		twocaptcha.WithCallbackURL(settings.CallbackURL),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newAntiCaptcha(apiKey string, settings Settings) (unicap.Provider, error) {
	client, err := anticaptcha.New(apiKey,
		anticaptcha.WithHTTPClient(settings.HTTPClient),
		anticaptcha.WithLogger(settings.Logger),
		anticaptcha.WithBaseURL(settings.BaseURL),
		anticaptcha.WithCallbackURL(settings.CallbackURL),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
