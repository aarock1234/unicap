# upicap

Universal captcha solving SDK for Go. One interface, multiple providers.

```go
provider, _ := capsolver.NewCapSolverProvider("API_KEY")
client, _ := upicap.NewClient(provider)

solution, _ := client.Solve(ctx, &tasks.ReCaptchaV2Task{
    WebsiteURL: "https://example.com",
    WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
})
```

## Features

- **Provider abstraction**: Switch providers without changing code
- **Type-safe API**: Proper structs, no `interface{}` abuse
- **Sync & async**: Auto-polling or manual control
- **Proxy support**: Optional for all task types
- **Extensible**: Create custom providers easily

## Supported Providers

| Provider    | Status       |
| ----------- | ------------ |
| CapSolver   | Full support |
| 2Captcha    | Full support |
| AntiCaptcha | Full support |

## Supported Captcha Types

| Type                 | CapSolver | 2Captcha | AntiCaptcha |
| -------------------- | --------- | -------- | ----------- |
| ReCaptcha V2         | ✓         | ✓        | ✓           |
| ReCaptcha V3         | ✓         | ✓        | ✓           |
| ReCaptcha Enterprise | ✓         | ✓        | ✓           |
| hCaptcha             | ✓         | ✓        | ✓           |
| FunCaptcha           | ✓         | ✓        | ✓           |
| Turnstile            | ✓         | ✓        | ✓           |
| GeeTest V3           | ✓         | ✓        | ✓           |
| GeeTest V4           | ✓         | ✓        | ✓           |
| Cloudflare Challenge | ✓         | -        | -           |
| DataDome             | ✓         | ✓        | -           |
| Image to Text        | ✓         | ✓        | ✓           |

## Installation

```bash
go get upicap
```

## Quick Start

### Synchronous (auto-polling)

```go
import (
    "upicap/internal/providers/capsolver"
    "upicap/pkg/upicap"
    "upicap/pkg/upicap/tasks"
)

provider, _ := capsolver.NewCapSolverProvider("API_KEY")
client, _ := upicap.NewClient(provider)

solution, err := client.Solve(ctx, &tasks.ReCaptchaV2Task{
    WebsiteURL: "https://example.com",
    WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
})

fmt.Println(solution.Token)
```

### Asynchronous (manual control)

```go
taskID, _ := client.CreateTask(ctx, task)

// Do other work...

result, _ := client.GetTaskResult(ctx, taskID)
switch result.Status {
case upicap.TaskStatusReady:
    fmt.Println(result.Solution.Token)
case upicap.TaskStatusProcessing:
    // Check again later
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
providers := []upicap.Provider{
    mustProvider(capsolver.NewCapSolverProvider("PRIMARY_KEY")),
    mustProvider(twocaptcha.NewTwoCaptchaProvider("BACKUP_KEY")),
}

for _, provider := range providers {
    client, _ := upicap.NewClient(provider)
    if solution, err := client.Solve(ctx, task); err == nil {
        fmt.Printf("solved by %s\n", provider.Name())
        break
    }
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

## Architecture

```
upicap/
├── pkg/upicap/              # Public API
│   ├── client.go            # Main client
│   ├── provider.go          # Provider interface
│   ├── task.go              # Task types
│   ├── provider_helpers.go  # Helpers for custom providers
│   └── tasks/               # Task implementations
├── internal/providers/      # Provider implementations
│   ├── capsolver/
│   ├── twocaptcha/
│   └── anticaptcha/
└── examples/                # Examples
```

## License

MIT
