package providers

import (
	"fmt"

	"github.com/aarock1234/unicap/providers/anticaptcha"
	"github.com/aarock1234/unicap/providers/capsolver"
	"github.com/aarock1234/unicap/providers/twocaptcha"
	"github.com/aarock1234/unicap/unicap"
)

var providerRegistry = make(map[string]ProviderFactory)

// ProviderFactory creates a provider instance
type ProviderFactory func(apiKey string) (unicap.Provider, error)

// RegisterProvider adds a provider to the registry
func RegisterProvider(name string, factory ProviderFactory) {
	providerRegistry[name] = factory
}

// GetProvider retrieves a provider factory by name
func GetProvider(name string) (ProviderFactory, bool) {
	factory, ok := providerRegistry[name]
	return factory, ok
}

// GetProviderNames returns all registered provider names
func GetProviderNames() []string {
	names := make([]string, 0, len(providerRegistry))
	for name := range providerRegistry {
		names = append(names, name)
	}
	return names
}

// NewProvider creates a provider by name and API key
func NewProvider(name, apiKey string) (unicap.Provider, error) {
	factory, ok := GetProvider(name)
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
	return factory(apiKey)
}

func init() {
	RegisterProvider("capsolver", func(apiKey string) (unicap.Provider, error) {
		return capsolver.NewCapSolverProvider(apiKey)
	})

	RegisterProvider("2captcha", func(apiKey string) (unicap.Provider, error) {
		return twocaptcha.NewTwoCaptchaProvider(apiKey)
	})

	RegisterProvider("anticaptcha", func(apiKey string) (unicap.Provider, error) {
		return anticaptcha.NewAntiCaptchaProvider(apiKey)
	})
}
