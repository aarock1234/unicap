package noptcha

import (
	"github.com/aarock1234/unicap/v3"
)

// proxy is Noptcha's proxy object, which goes inside a task's params. Its
// types are unicap's proxy types, and the login and password are optional.
type proxy struct {
	Type     unicap.ProxyType `json:"type"`
	Host     string           `json:"host"`
	Port     int              `json:"port"`
	Login    string           `json:"login,omitempty"`
	Password string           `json:"password,omitempty"`
}

// newProxy converts an optional proxy to Noptcha's proxy object. A nil proxy
// yields nil, which leaves the proxy out of the params.
func newProxy(source *unicap.Proxy) *proxy {
	if source == nil {
		return nil
	}

	return &proxy{
		Type:     source.Type,
		Host:     source.Host,
		Port:     source.Port,
		Login:    source.Username,
		Password: source.Password,
	}
}
