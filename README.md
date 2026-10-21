# unicap

[![Go Reference](https://pkg.go.dev/badge/github.com/aarock1234/unicap/v3.svg)](https://pkg.go.dev/github.com/aarock1234/unicap/v3)
[![Release](https://img.shields.io/github/v/release/aarock1234/unicap)](https://github.com/aarock1234/unicap/releases)
[![License](https://img.shields.io/badge/license-Elastic%202.0-blue.svg)](LICENSE)

A universal captcha-solving SDK for Go. Write a task once and solve it with CapSolver, 2Captcha, Anti-Captcha, or your own provider. Results and errors look the same whichever provider you use.

```go
provider, err := capsolver.New(os.Getenv("CAPSOLVER_API_KEY"))
if err != nil {
	return err
}

client, err := unicap.New(provider)
if err != nil {
	return err
}

result, err := client.Solve(ctx, &tasks.ReCaptchaV2Task{
	WebsiteURL: "https://www.google.com/recaptcha/api2/demo",
	WebsiteKey: "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-",
})
if err != nil {
	return err
}

fmt.Println(result.Solution.Token)
```

## Installation

```bash
go get github.com/aarock1234/unicap/v3@latest
```

Requires Go 1.27 or newer. If you're upgrading from v1 or v2, see [MIGRATION.md](MIGRATION.md).

## Features

- **Tasks**: 51 typed captcha tasks that work across providers, plus `RawTask` for anything unicap doesn't model yet
- **Strict**: options a provider can't honour return `ErrUnsupportedOption` instead of being silently dropped
- **Results**: task ID, cost, and timing, with typed solutions for tokens, cookies, user agents, and coordinates
- **Errors**: provider codes map to sentinels for `errors.Is`, and `unicap.IsRetryable` tells you when to resubmit
- **Polling**: each provider's recommended pace by default, and instant results for tasks solved synchronously
- **Account**: balance checks, answer reports, affiliate IDs, callback URLs, and language pools

## Usage

### Solve and wait

`Solve` validates the task, submits it, and polls until the solution is ready, the context is canceled, or polling times out.

```go
result, err := client.Solve(ctx, &tasks.TurnstileTask{
	WebsiteURL: "https://example.com/login",
	WebsiteKey: "0x4AAAAAAA...",
})
```

### Submit now, collect later

```go
created, err := client.CreateTask(ctx, task)
if err != nil {
	return err
}

// Recognition tasks some providers solve synchronously are already ready.
if !created.Ready() {
	result, err := client.Wait(ctx, created.TaskID)
	// ...
}
```

`client.GetTaskResult(ctx, taskID)` checks the task once without waiting.

### Cookies, user agents, and family-specific solutions

```go
result, err := client.Solve(ctx, &tasks.DataDomeTask{
	WebsiteURL: "https://example.com",
	CaptchaURL: captchaURL,
	UserAgent:  userAgent,
	Proxy:      proxy,
})
if err != nil {
	return err
}

cookie, ok := result.Solution.Cookie("datadome")
```

```go
result, err := client.Solve(ctx, &tasks.GeeTestV4Task{
	WebsiteURL: "https://example.com",
	CaptchaID:  captchaID,
})
if err != nil {
	return err
}

geetest, err := unicap.DecodeSolution[tasks.GeeTestV4Solution](result.Solution)
```

### Proxies

```go
proxy, err := unicap.ParseProxyURL("socks5://user:pass@203.0.113.10:1080")
if err != nil {
	return err
}

task := &tasks.ReCaptchaV2Task{
	WebsiteURL: "https://example.com",
	WebsiteKey: siteKey,
	UserAgent:  userAgent,
	Proxy:      proxy,
}
```

A nil `Proxy` means the provider solves the task from its own network. When you log a `*unicap.Proxy` or print it with `String`, the password is redacted.

### Errors and retries

```go
result, err := client.Solve(ctx, task)

switch {
case errors.Is(err, unicap.ErrInsufficientFunds):
	// Top up the account.
case errors.Is(err, unicap.ErrUnsupportedOption):
	// The provider cannot honour one of the task's options.
case unicap.IsRetryable(err):
	// ErrUnsolvable or ErrUnavailable: submitting again may succeed.
}

if providerErr, ok := errors.AsType[*unicap.Error](err); ok {
	fmt.Println(providerErr.Provider, providerErr.Code, providerErr.Message)
}
```

To have `Solve` resubmit after retryable failures, use `unicap.WithSolveRetries(n)`. Providers may bill for every attempt.

### Balance and reports

```go
balance, err := client.Balance(ctx)

// After the target site accepts or rejects the token:
err = client.Report(ctx, result, accepted)
```

Both calls return `unicap.ErrUnsupportedOperation` when the provider doesn't offer that feature.

### Choosing a provider at runtime

```go
solver, err := provider.NewRegistry().New("2captcha", apiKey,
	provider.WithLogger(logger),
)
```

### Client options

| Option                          | Effect                                                                       |
| ------------------------------- | ---------------------------------------------------------------------------- |
| `unicap.WithLogger(logger)`     | Logs task lifecycle events. Logs are discarded by default.                   |
| `unicap.WithPolling(polling)`   | Overrides the provider's recommended polling. Start from `unicap.DefaultPolling()`. |
| `unicap.WithSolveRetries(n)`    | Resubmits up to `n` times after retryable failures.                          |

### Provider options

| Option                | CapSolver | 2Captcha | Anti-Captcha |
| --------------------- | --------- | -------- | ------------ |
| `WithHTTPClient`      | ✓         | ✓        | ✓            |
| `WithLogger`          | ✓         | ✓        | ✓            |
| `WithBaseURL`         | ✓         | ✓        | ✓            |
| `WithCallbackURL`     | ✓         | ✓        | ✓            |
| `WithAppID`           | ✓         | –        | –            |
| `WithSoftID`          | –         | ✓        | ✓            |
| `WithLanguagePool`    | –         | ✓        | per task     |

The Anti-Captcha client also has `PushVariable` for AntiGate templates.

## Supported captchas

| Task                                     | CapSolver | 2Captcha | Anti-Captcha |
| ---------------------------------------- | :-------: | :------: | :----------: |
| `ReCaptchaV2Task`                        | ✓         | ✓        | ✓            |
| `ReCaptchaV2EnterpriseTask`              | ✓         | ✓        | ✓            |
| `ReCaptchaV3Task`                        | ✓         | ✓        | ✓            |
| `ReCaptchaV3EnterpriseTask`              | ✓         | ✓        | ✓            |
| `ReCaptchaClassificationTask`            | ✓         | –        | –            |
| `HCaptchaTask`¹                          | –         | –        | –            |
| `FunCaptchaTask`                         | –         | ✓        | ✓            |
| `TurnstileTask`                          | ✓         | ✓        | ✓            |
| `CloudflareChallengeTask`                | ✓         | –        | –            |
| `DataDomeTask`                           | ✓         | ✓        | –            |
| `GeeTestTask`, `GeeTestV4Task`           | ✓         | ✓        | ✓            |
| `AWSWAFTask`                             | ✓         | ✓        | ✓            |
| `AWSWAFClassificationTask`               | ✓         | –        | –            |
| `MTCaptchaTask`                          | ✓         | ✓        | –            |
| `FriendlyCaptchaTask`                    | –         | ✓        | ✓            |
| `ProsopoTask`                            | –         | ✓        | ✓            |
| `AltchaTask`                             | –         | ✓        | ✓            |
| `BotDeflectorTask`                       | ✓         | –        | –            |
| `AntiGateTask`                           | –         | –        | ✓            |
| `LeminTask`, `CutCaptchaTask`, `CapyTask`, `KeyCaptchaTask` | – | ✓ | – |
| `CyberSiAraTask`, `AtbCaptchaTask`, `TencentTask`, `CaptchaFoxTask` | – | ✓ | – |
| `YandexSmartCaptchaTask`, `VKCaptchaTask`, `YidunTask`, `BinanceTask` | – | ✓ | – |
| `HuntTask`, `TSPDTask`, `BasiliskTask`, `ImpervaTask`, `AlibabaTask` | – | ✓ | – |
| `ImageToTextTask`                        | ✓         | ✓        | ✓            |
| `CoordinatesTask`                        | –         | ✓        | ✓            |
| `VisionEngineTask`                       | ✓         | –        | –            |
| `TextCaptchaTask`, `AudioTask`, `GridTask`, `RotateTask` | – | ✓ | – |
| `DrawAroundTask`, `BoundingBoxTask`, `DragAndDropTask` | – | ✓ | – |
| `TemuTask`, `VKCaptchaImageTask`, `YandexImageTask`, `YandexPuzzleTask` | – | ✓ | – |
| `RawTask` (passthrough)                  | ✓         | ✓        | ✓            |

¹ All three providers have withdrawn hCaptcha support, so the built-in providers reject it. The task type remains available for custom providers.

`client.Supports(taskType)` reports support at runtime.

## Custom providers

Implement `unicap.Provider`:

```go
type Provider interface {
	Name() string
	CreateTask(ctx context.Context, task Task) (*Result, error)
	GetTaskResult(ctx context.Context, taskID string) (*Result, error)
}
```

Return provider failures as `*unicap.Error`, and wrap a sentinel in its `Err` field so callers can use `errors.Is`. You can also implement `BalanceChecker`, `Reporter`, `TaskSupporter`, and `PollingAdvisor`. [`examples/custom_provider`](examples/custom_provider) has a complete implementation.

## Examples

- [`examples/solve`](examples/solve): choose a provider at runtime, check the balance, solve, and report. Run it with `make example provider=2captcha` after setting `CAPSOLVER_API_KEY`, `TWOCAPTCHA_API_KEY`, or `ANTICAPTCHA_API_KEY`.
- [`examples/custom_provider`](examples/custom_provider): plug a custom provider into unicap.

## Development

```bash
make help     # list targets
make lint     # golangci-lint
make format   # gofmt and goimports
make vet
make build
```

## License

[Elastic License 2.0](LICENSE)
