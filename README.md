# upicap

[![Go Reference](https://pkg.go.dev/badge/github.com/aarock1234/unicap.svg)](https://pkg.go.dev/github.com/aarock1234/unicap)
[![Release](https://img.shields.io/github/v/release/aarock1234/unicap)](https://github.com/aarock1234/unicap/releases)
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

    "github.com/aarock1234/unicap/internal/providers/capsolver"
    "github.com/aarock1234/unicap/pkg/upicap"
    "github.com/aarock1234/unicap/pkg/upicap/tasks"
)

func main() {
    // Initialize the provider, in this case CapSolver, with your API key
    provider, err := capsolver.NewCapSolverProvider("YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }

    // Create a client with the provider
    client, err := upicap.NewClient(provider)
    if err != nil {
        log.Fatal(err)
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
        log.Fatal(err)
    }

    fmt.Printf("Token: %s\n", solution.Token)

    // Verify the token by submitting it to the demo page
    data := url.Values{
        "g-recaptcha-response": {solution.Token},
    }

    // Create a POST request to verify the token
    req, err := http.NewRequestWithContext(ctx, "POST", "https://www.google.com/recaptcha/api2/demo", strings.NewReader(data.Encode()))
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    // Send the verification request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    // Expected output: "Verification Success... Hooray!"
    fmt.Printf("Verification response: %s\n", string(body))
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
go get -u github.com/aarock1234/unicap
```

## Quick Start

### Synchronous (auto-polling)

```go
provider, err := capsolver.NewCapSolverProvider("API_KEY")
if err != nil {
    log.Fatal(err)
}

client, err := upicap.NewClient(provider)
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
case upicap.TaskStatusReady:
    fmt.Println(result.Solution.Token)
case upicap.TaskStatusProcessing:
    fmt.Println("Still processing, check again later")
case upicap.TaskStatusFailed:
    log.Printf("Task failed: %v", result.Error)
}
```

### With Proxy

```go
task := &tasks.ReCaptchaV2Task{
    WebsiteURL: "https://example.com",
    WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
    Proxy: &upicap.Proxy{
        Type:     upicap.ProxyTypeHTTP,
        Address:  "proxy.example.com",
        Port:     8080,
        Login:    "user",
        Password: "pass",
    },
}
```

### Multi-Provider Failover

```go
providers := []upicap.Provider{}

if p, err := capsolver.NewCapSolverProvider("PRIMARY_KEY"); err == nil {
    providers = append(providers, p)
}
if p, err := twocaptcha.NewTwoCaptchaProvider("BACKUP_KEY"); err == nil {
    providers = append(providers, p)
}

for _, provider := range providers {
    client, err := upicap.NewClient(provider)
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
    Numeric:   0, // 0=any, 1=numbers only
    MinLength: 4,
    MaxLength: 20,
}
```

## Custom Providers

Implement the `Provider` interface:

```go
type Provider interface {
    CreateTask(ctx context.Context, task Task) (string, error)
    GetTaskResult(ctx context.Context, taskID string) (*TaskResult, error)
    Name() string
}
```

Use the built-in helpers:

```go
type customProvider struct {
    apiKey string
    client *upicap.BaseHTTPClient
    errors *upicap.ErrorMapper
}

func NewCustomProvider(apiKey string) (upicap.Provider, error) {
    return &customProvider{
        apiKey: apiKey,
        client: &upicap.BaseHTTPClient{
            HTTPClient: &http.Client{Timeout: 30 * time.Second},
            Logger:     slog.New(slog.NewTextHandler(io.Discard, nil)),
            BaseURL:    "https://api.yourservice.com",
        },
        errors: upicap.StandardErrorMapper(
            "yourservice",
            []string{"INVALID_KEY"},
            []string{"NO_FUNDS"},
            []string{"NOT_FOUND"},
            []string{"BAD_REQUEST"},
        ),
    }, nil
}

func (p *customProvider) CreateTask(ctx context.Context, task upicap.Task) (string, error) {
    req := createTaskRequest{APIKey: p.apiKey, Task: mapTask(task)}
    var resp createTaskResponse

    if err := p.client.DoJSON(ctx, "/create", req, &resp); err != nil {
        return "", err
    }

    if resp.Error != "" {
        return "", p.errors.MapError(resp.ErrorCode, resp.Error)
    }

    return resp.TaskID, nil
}
```

See `examples/custom_provider/` for complete implementation.

## Configuration

### Custom Logger

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

client, _ := upicap.NewClient(
    provider,
    upicap.WithLogger(logger),
)
```

### Custom Polling

```go
poller := upicap.NewPoller(provider, upicap.PollerConfig{
    InitialInterval: 1 * time.Second,
    MaxInterval:     10 * time.Second,
    Timeout:         3 * time.Minute,
    Multiplier:      2.0,
})

client, _ := upicap.NewClient(
    provider,
    upicap.WithPoller(poller),
)
```

## Error Handling

```go
solution, err := client.Solve(ctx, task)
if err != nil {
    switch {
    case errors.Is(err, upicap.ErrInvalidAPIKey):
        // Handle invalid API key
    case errors.Is(err, upicap.ErrInsufficientFunds):
        // Handle low balance
    case errors.Is(err, upicap.ErrTimeout):
        // Handle timeout
    default:
        // Handle other errors
    }
}
```

## License

Elastic License 2.0

Free to use, cannot resell as a service. See [LICENSE](LICENSE) for details.
