package unicap

// ProxyType represents the transport type of a proxy.
type ProxyType string

const (
	// ProxyTypeHTTP routes requests through an HTTP proxy.
	ProxyTypeHTTP ProxyType = "http"
	// ProxyTypeHTTPS routes requests through an HTTPS proxy.
	ProxyTypeHTTPS ProxyType = "https"
	// ProxyTypeSOCKS4 routes requests through a SOCKS4 proxy.
	ProxyTypeSOCKS4 ProxyType = "socks4"
	// ProxyTypeSOCKS5 routes requests through a SOCKS5 proxy.
	ProxyTypeSOCKS5 ProxyType = "socks5"
)

// Proxy holds proxy configuration for captcha solving.
type Proxy struct {
	Type     ProxyType
	Address  string
	Port     int
	Login    string
	Password string
}

// IsSet reports whether the proxy is configured.
func (p *Proxy) IsSet() bool {
	return p != nil && p.Address != ""
}
