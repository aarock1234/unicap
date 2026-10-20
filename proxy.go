package unicap

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// ProxyType is the transport protocol of a proxy.
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

// Valid reports whether t is one of the declared proxy types.
func (t ProxyType) Valid() bool {
	switch t {
	case ProxyTypeHTTP, ProxyTypeHTTPS, ProxyTypeSOCKS4, ProxyTypeSOCKS5:
		return true
	default:
		return false
	}
}

// ErrInvalidProxy reports that a proxy configuration is malformed.
var ErrInvalidProxy = errors.New("invalid proxy")

// Proxy is the proxy a provider's workers use to load the target page. A nil
// *Proxy on a task means the provider solves the task from its own network.
type Proxy struct {
	// Type is the proxy protocol.
	Type ProxyType

	// Host is the proxy's IP address or hostname. Some providers accept IP
	// addresses only.
	Host string

	// Port is the proxy's TCP port.
	Port int

	// Username is the optional proxy login.
	Username string

	// Password is the optional proxy password.
	Password string
}

// ParseProxyURL parses a proxy URL such as "socks5://user:pass@1.2.3.4:1080".
// A URL without a scheme is treated as an HTTP proxy.
func ParseProxyURL(raw string) (*Proxy, error) {
	raw = strings.TrimSpace(raw)
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}

	// The parse error is not wrapped because its message echoes the input,
	// which may contain the proxy password.
	parsed, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("parsing proxy url: %w", ErrInvalidProxy)
	}

	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		return nil, fmt.Errorf("proxy port must be an integer: %w", ErrInvalidProxy)
	}

	proxy := &Proxy{
		Type: ProxyType(strings.ToLower(parsed.Scheme)),
		Host: parsed.Hostname(),
		Port: port,
	}

	if parsed.User != nil {
		proxy.Username = parsed.User.Username()
		proxy.Password, _ = parsed.User.Password()
	}

	if err := proxy.Validate(); err != nil {
		return nil, err
	}

	return proxy, nil
}

// Validate reports whether the proxy has a known type, a host, and a port in
// the range 1-65535.
func (p *Proxy) Validate() error {
	if !p.Type.Valid() {
		return fmt.Errorf("proxy type %q: %w", p.Type, ErrInvalidProxy)
	}

	if p.Host == "" {
		return fmt.Errorf("proxy host is required: %w", ErrInvalidProxy)
	}

	if p.Port < 1 || p.Port > 65535 {
		return fmt.Errorf("proxy port %d is out of range: %w", p.Port, ErrInvalidProxy)
	}

	if p.Password != "" && p.Username == "" {
		return fmt.Errorf("proxy password requires a username: %w", ErrInvalidProxy)
	}

	return nil
}

// HasCredentials reports whether the proxy carries a username.
func (p *Proxy) HasCredentials() bool {
	return p.Username != ""
}

// HostIsIP reports whether Host is a literal IPv4 or IPv6 address.
func (p *Proxy) HostIsIP() bool {
	return net.ParseIP(p.Host) != nil
}

// URL returns the proxy as a URL, including credentials.
func (p *Proxy) URL() *url.URL {
	proxyURL := &url.URL{
		Scheme: string(p.Type),
		Host:   net.JoinHostPort(p.Host, strconv.Itoa(p.Port)),
	}

	if p.HasCredentials() {
		proxyURL.User = url.UserPassword(p.Username, p.Password)
	}

	return proxyURL
}

// String returns the proxy URL with the password redacted.
func (p *Proxy) String() string {
	return p.URL().Redacted()
}

// LogValue renders the proxy for structured logs without its password.
func (p *Proxy) LogValue() slog.Value {
	return slog.StringValue(p.String())
}
