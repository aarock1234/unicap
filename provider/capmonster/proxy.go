package capmonster

import (
	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// proxyFields converts an optional proxy to CapMonster Cloud's fields.
// CapMonster Cloud documents proxyAddress as an IPv4 or IPv6 address and
// never mentions hostnames, so a hostname is rejected. A nil proxy yields
// the zero value.
func proxyFields(proxy *unicap.Proxy) (solverapi.ProxyFields, error) {
	if err := solverapi.RequireIPProxy(proxy); err != nil {
		return solverapi.ProxyFields{}, err
	}

	return solverapi.NewProxyFields(proxy), nil
}

// requiredProxyFields converts the proxy of a task CapMonster Cloud solves
// only through the caller's proxy.
func requiredProxyFields(proxy *unicap.Proxy) (solverapi.ProxyFields, error) {
	if err := solverapi.Required(
		solverapi.Setting("Proxy", proxy != nil),
	); err != nil {
		return solverapi.ProxyFields{}, err
	}

	return proxyFields(proxy)
}
