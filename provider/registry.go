// Package provider contains runtime registry helpers for selecting built-in
// unicap providers by name.
package provider

import (
	"fmt"
	"sort"
	"sync"

	"github.com/aarock1234/unicap"
	"github.com/aarock1234/unicap/provider/anticaptcha"
	"github.com/aarock1234/unicap/provider/capsolver"
	"github.com/aarock1234/unicap/provider/twocaptcha"
)

// Factory creates a provider instance.
type Factory func(apiKey string) (unicap.Provider, error)

// Registry maps provider names to provider factories. It is safe for
// concurrent use.
type Registry struct {
	mu        sync.RWMutex
	factories map[string]Factory
}

// NewRegistry creates a registry preloaded with the built-in providers.
func NewRegistry() *Registry {
	r := &Registry{
		factories: make(map[string]Factory),
	}

	r.Register("capsolver", func(apiKey string) (unicap.Provider, error) {
		return capsolver.New(apiKey)
	})
	r.Register("2captcha", func(apiKey string) (unicap.Provider, error) {
		return twocaptcha.New(apiKey)
	})
	r.Register("anticaptcha", func(apiKey string) (unicap.Provider, error) {
		return anticaptcha.New(apiKey)
	})

	return r
}

// Register adds or replaces a provider factory in the registry.
func (r *Registry) Register(name string, factory Factory) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.factories[name] = factory
}

// Get retrieves a provider factory by name.
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
	sort.Strings(names)

	return names
}

// New creates a provider by name and API key.
func (r *Registry) New(name, apiKey string) (unicap.Provider, error) {
	factory, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}

	return factory(apiKey)
}
