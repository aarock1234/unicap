# unicap

[![Go Reference](https://pkg.go.dev/badge/unicap.svg)](https://pkg.go.dev/unicap)
[![Release](https://img.shields.io/github/v/release/aarock1234/unicap)](https://unicap/releases)
[![License](https://img.shields.io/badge/license-Elastic%202.0-blue.svg)](LICENSE)

Universal captcha solving SDK for Go. One interface, multiple providers.

```go
package main

import (
    "context"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "strings"
    "time"

    "github.com/aarock1234/unicap"
    "github.com/aarock1234/unicap/provider/capsolver"
    "github.com/aarock1234/unicap/tasks"
)

func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

func run() error {
    // Initialize the provider, in this case CapSolver, with your API key
    provider, err := capsolver.New("YOUR_API_KEY")
    if err != nil {
        return err
    }

    // Create a client with the provider
    client, err := unicap.New(provider)
    if err != nil {
        return err
    }

    // Set up a context with timeout for the captcha solving operation
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // Solve the ReCaptcha V2 challenge
    // This will automatically poll until the solution is ready
    solution, err := client.Solve(ctx, &tasks.ReCaptchaV2Task{
        WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
        WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
    })
    if err != nil {
        return err
    }

    fmt.Printf("Token: %s\n", solution.Token)

    // Verify the token by submitting it to the demo page
    data := url.Values{
        "g-recaptcha-response": {solution.Token},
    }

    // Create a POST request to verify the token
    req, err := http.NewRequestWithContext(ctx, "POST", "https://www.google.com/recaptcha/api2/demo", strings.NewReader(data.Encode()))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Send the verification request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer func() { _ = resp.Body.Close() }()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    // Expected output: "Verification Success... Hooray!"
    fmt.Printf("Verification response: %s\n", string(body))

    return nil
}
```

## Features

- **Provider abstraction**: Switch providers without changing code
- **Type-safe API**: Proper structs, no `interface{}` abuse
- **Sync & async**: Auto-polling or manual control
- **Proxy support**: Optional for all task types
- **Extensible**: Create custom providers easily

## Supported Providers & Captcha Types

| Type/Provider        | [CapSolver](https://capsolver.com/) | [2Captcha](https://2captcha.com/) | [AntiCaptcha](https://anti-captcha.com/) |
| -------------------- | ----------------------------------- | --------------------------------- | ---------------------------------------- |
| Image to Text        | ✓                                   | ✓                                 | ✓                                        |
| ReCaptcha V2         | ✓                                   | ✓                                 | ✓                                        |
| ReCaptcha V3         | ✓                                   | ✓                                 | ✓                                        |
| ReCaptcha Enterprise | ✓                                   | ✓                                 | ✓                                        |
| hCaptcha             | ✓                                   | ✓                                 | ✓                                        |
| FunCaptcha           | ✓                                   | ✓                                 | ✓                                        |
| Turnstile            | ✓                                   | ✓                                 | ✓                                        |
| GeeTest V3           | ✓                                   | ✓                                 | ✓                                        |
| GeeTest V4           | ✓                                   | ✓                                 | ✓                                        |
| Cloudflare Challenge | ✓                                   | -                                 | -                                        |
| DataDome             | ✓                                   | ✓                                 | -                                        |

## Installation

```bash
go get github.com/aarock1234/unicap@latest
```

## Quick Start

The following snippets assume the same imports as the full example above.

### Synchronous (auto-polling)

```go
provider, err := capsolver.New("API_KEY")
if err != nil {
    log.Fatal(err)
}

client, err := unicap.New(provider)
if err != nil {
    log.Fatal(err)
}

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

solution, err := client.Solve(ctx, &tasks.ReCaptchaV2Task{
    WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
    WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
})
if err != nil {
    log.Fatal(err)
}

fmt.Println(solution.Token)
```

### Asynchronous (manual control)

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

task := &tasks.ReCaptchaV2Task{
    WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
    WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
}

taskID, err := client.CreateTask(ctx, task)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Task ID: %s\n", taskID)

// Do other work...
time.Sleep(10 * time.Second)

result, err := client.GetTaskResult(ctx, taskID)
if err != nil {
    log.Fatal(err)
}

switch result.Status {
case unicap.TaskStatusReady:
    fmt.Println(result.Solution.Token)
case unicap.TaskStatusProcessing:
    fmt.Println("Still processing, check again later")
case unicap.TaskStatusFailed:
    log.Printf("Task failed: %v", result.Error)
}
```

### With Proxy

```go
task := &tasks.ReCaptchaV2Task{
    WebsiteURL: "https://example.com",
    WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
    Proxy: &unicap.Proxy{
        Type:     unicap.ProxyTypeHTTP,
        Address:  "proxy.example.com",
        Port:     8080,
        Login:    "user",
        Password: "pass",
    },
}
```

### Multi-Provider Failover

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

task := &tasks.ReCaptchaV2Task{
    WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
    WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
}

providers := []unicap.Provider{}

if p, err := capsolver.New("PRIMARY_KEY"); err == nil {
    providers = append(providers, p)
}
if p, err := twocaptcha.New("BACKUP_KEY"); err == nil {
    providers = append(providers, p)
}

for _, provider := range providers {
    client, err := unicap.New(provider)
    if err != nil {
        continue
    }

    solution, err := client.Solve(ctx, task)
    if err != nil {
        log.Printf("provider %s failed: %v", provider.Name(), err)
        continue
    }

    fmt.Printf("Solved by %s: %s\n", provider.Name(), solution.Token)
    break
}
```

## Task Types

### ReCaptcha V2

```go
&tasks.ReCaptchaV2Task{
    WebsiteURL:  "https://example.com",
    WebsiteKey:  "site-key",
    IsInvisible: false,
    DataS:       "optional-data-s",
    Proxy:       proxy, // optional
}
```

### ReCaptcha V3

```go
&tasks.ReCaptchaV3Task{
    WebsiteURL: "https://example.com",
    WebsiteKey: "site-key",
    PageAction: "verify",
    MinScore:   0.7,
    Proxy:      proxy, // optional
}
```

### hCaptcha

```go
&tasks.HCaptchaTask{
    WebsiteURL: "https://example.com",
    WebsiteKey: "site-key",
    Proxy:      proxy, // optional
}
```

### FunCaptcha

```go
&tasks.FunCaptchaTask{
    WebsiteURL:       "https://example.com",
    WebsitePublicKey: "public-key",
    APIJSSubdomain:   "api.arkoselabs.com", // optional
    Proxy:            proxy, // optional
}
```

### Turnstile

```go
&tasks.TurnstileTask{
    WebsiteURL: "https://example.com",
    WebsiteKey: "site-key",
    Action:     "login", // optional
    Proxy:      proxy, // optional
}
```

### GeeTest V3

```go
&tasks.GeeTestTask{
    WebsiteURL:         "https://example.com",
    GT:                 "gt-value",
    Challenge:          "challenge-value",
    APIServerSubdomain: "api.geetest.com", // optional
    Proxy:              proxy, // optional
}
```

### GeeTest V4

```go
&tasks.GeeTestV4Task{
    WebsiteURL: "https://example.com",
    CaptchaID:  "captcha-id",
    Proxy:      proxy, // optional
}
```

### Cloudflare Challenge

```go
&tasks.CloudflareChallengeTask{
    WebsiteURL: "https://example.com",
    HTML:       "<html>...</html>", // optional
    UserAgent:  "Mozilla/5.0...",
    Proxy:      proxy, // required
}
```

### DataDome

```go
&tasks.DataDomeTask{
    WebsiteURL: "https://example.com",
    CaptchaURL: "https://geo.captcha-delivery.com/...",
    UserAgent:  "Mozilla/5.0...",
    Proxy:      proxy, // required
}
```

### Image Recognition

```go
&tasks.ImageToTextTask{
    Body:      "base64-encoded-image",
    Numeric:   tasks.NumericModeAny,
    MinLength: 4,
    MaxLength: 20,
}
```

## Custom Providers

Implement the `unicap.Provider` interface:

```go
type Provider interface {
    CreateTask(ctx context.Context, task unicap.Task) (string, error)
    GetTaskResult(ctx context.Context, taskID string) (*unicap.TaskResult, error)
    Name() string
}
```

Build your own transport and mapping logic inside the provider implementation:

```go
var _ unicap.Provider = (*customProvider)(nil)

type customProvider struct {
    apiKey string
    client *http.Client
}

func NewCustomProvider(apiKey string) (unicap.Provider, error) {
    return &customProvider{
        apiKey: apiKey,
        client: &http.Client{Timeout: 30 * time.Second},
    }, nil
}

func (p *customProvider) CreateTask(ctx context.Context, task unicap.Task) (string, error) {
    // translate the SDK task into your provider's wire format,
    // send the request, then map the response into a task ID.
    return "", nil
}

func (p *customProvider) GetTaskResult(ctx context.Context, taskID string) (*unicap.TaskResult, error) {
    // fetch the provider-specific result and map it into a unicap.TaskResult.
    return nil, nil
}

func (p *customProvider) Name() string {
    return "custom"
}
```

See `examples/custom_provider/` for complete implementation.

## Configuration

### Custom Logger

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client, err := unicap.New(
    provider,
    unicap.WithLogger(logger),
)
if err != nil {
    return err
}
```

### Custom Polling

```go
poller := unicap.NewPoller(provider, unicap.PollerConfig{
    InitialInterval: 1 * time.Second,
    MaxInterval:     10 * time.Second,
    Timeout:         3 * time.Minute,
    Multiplier:      2.0,
})

client, err := unicap.New(
    provider,
    unicap.WithPoller(poller),
)
if err != nil {
    return err
}
```

## Error Handling

```go
solution, err := client.Solve(ctx, task)
if err != nil {
    switch {
    case errors.Is(err, unicap.ErrInvalidAPIKey):
        // Handle invalid API key
    case errors.Is(err, unicap.ErrInsufficientFunds):
        // Handle low balance
    case errors.Is(err, unicap.ErrTimeout):
        // Handle timeout
    default:
        // Handle other errors
    }
}
```

## License

Elastic License 2.0

Free to use, cannot resell as a service. See [LICENSE](LICENSE) for details.
