package upicap

// ProxyType represents the type of proxy
type ProxyType string

const (
	ProxyTypeHTTP   ProxyType = "http"
	ProxyTypeHTTPS  ProxyType = "https"
	ProxyTypeSOCKS4 ProxyType = "socks4"
	ProxyTypeSOCKS5 ProxyType = "socks5"
)

// Proxy contains proxy configuration for captcha solving
type Proxy struct {
	Type     ProxyType
	Address  string
	Port     int
	Login    string
	Password string
}

// IsSet returns true if the proxy is configured
func (p *Proxy) IsSet() bool {
	return p != nil && p.Address != ""
}
