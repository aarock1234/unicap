package nopecha

import (
	"slices"

	"github.com/aarock1234/unicap/v3"
	"github.com/aarock1234/unicap/v3/internal/solverapi"
)

// credentialSchemes are the proxy types NopeCHA takes a username and a
// password for. It ignores the credentials of a SOCKS proxy.
var credentialSchemes = []unicap.ProxyType{
	unicap.ProxyTypeHTTP,
	unicap.ProxyTypeHTTPS,
}

// proxy is NopeCHA's proxy object, the same on every token endpoint. Its
// schemes are unicap's proxy types, and the username and password are
// optional.
type proxy struct {
	Scheme   unicap.ProxyType `json:"scheme"`
	Host     string           `json:"host"`
	Port     int              `json:"port"`
	Username string           `json:"username,omitempty"`
	Password string           `json:"password,omitempty"`
}

// checkProxy reports the credentials of a proxy whose type NopeCHA takes
// none for as unsupported, since NopeCHA would ignore them. A nil proxy
// passes.
func checkProxy(source *unicap.Proxy) error {
	if source == nil {
		return nil
	}

	ignored := source.HasCredentials() && !slices.Contains(credentialSchemes, source.Type)

	return solverapi.Unsupported(
		solverapi.Setting("Proxy credentials with type "+string(source.Type), ignored),
	)
}

// newProxy converts an optional proxy to NopeCHA's proxy object. A nil proxy
// yields nil, which leaves the proxy out of the body.
func newProxy(source *unicap.Proxy) *proxy {
	if source == nil {
		return nil
	}

	return &proxy{
		Scheme:   source.Type,
		Host:     source.Host,
		Port:     source.Port,
		Username: source.Username,
		Password: source.Password,
	}
}
