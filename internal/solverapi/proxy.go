package solverapi

import (
	"fmt"

	"github.com/aarock1234/unicap"
)

// ProxyFields holds the proxy parameters shared by all solver-API task types.
// Embed it in a provider task struct to gain the standard proxy JSON fields.
type ProxyFields struct {
	ProxyType     string `json:"proxyType,omitempty"`
	ProxyAddress  string `json:"proxyAddress,omitempty"`
	ProxyPort     int    `json:"proxyPort,omitempty"`
	ProxyLogin    string `json:"proxyLogin,omitempty"`
	ProxyPassword string `json:"proxyPassword,omitempty"`
}

// ProxyFieldsFrom builds ProxyFields from a unicap proxy. A nil or unset proxy
// yields the zero value, whose fields are omitted from JSON.
func ProxyFieldsFrom(p *unicap.Proxy) ProxyFields {
	if !p.IsSet() {
		return ProxyFields{}
	}

	return ProxyFields{
		ProxyType:     string(p.Type),
		ProxyAddress:  p.Address,
		ProxyPort:     p.Port,
		ProxyLogin:    p.Login,
		ProxyPassword: p.Password,
	}
}

// ProxyString formats a proxy as "type:address:port:login:password" for
// providers that accept a single proxy string instead of discrete fields. A nil
// or unset proxy yields an empty string.
func ProxyString(p *unicap.Proxy) string {
	if !p.IsSet() {
		return ""
	}

	return fmt.Sprintf("%s:%s:%d:%s:%s", p.Type, p.Address, p.Port, p.Login, p.Password)
}
