# Migrating to v3

v3 changes the module path and the core API. The v2 tags were never installable, because `go.mod` lacked the `/v2` suffix. v3 is therefore the first major release after v1.

## Module path

```bash
go get github.com/aarock1234/unicap/v3@latest
```

Update your imports from `github.com/aarock1234/unicap/...` to `github.com/aarock1234/unicap/v3/...`. Go 1.27 or newer is required.

## Solve returns a Result

```go
// v1
solution, err := client.Solve(ctx, task)
token := solution.Token

// v3
result, err := client.Solve(ctx, task)
token := result.Solution.Token
```

`Result` carries `TaskID`, `Cost`, `CreatedAt`, `SolvedAt`, and the submitted `Task`. You need `TaskID` and `Task` for `client.Report`.

`CreateTask` now returns `*Result` instead of a task ID string. Use `result.TaskID`, and check `result.Ready()` for providers that solve some tasks synchronously. To poll an existing task, call `client.Wait(ctx, taskID)`.

## Solutions

- `Solution.Extra map[string]any` is gone. The provider payload is now `Solution.Raw`, as `json.RawMessage`.
- Use `unicap.DecodeSolution[T]` with types such as `tasks.GeeTestV4Solution`, `tasks.AWSWAFSolution`, or `tasks.AntiGateSolution`.
- `Solution.Cookie string` became `Solution.Cookies []*http.Cookie`. Look one up with `Solution.Cookie(name)`.
- New fields: `UserAgent`, `Boxes`, `Indexes`, and `Angle`. `Coordinates` is now actually filled in.

## Errors

- Provider failures are always returned as errors. `TaskStatusFailed` and `TaskResult.Error` are gone.
- `unicap.NewError(code, message, provider, retriable, wrapped)` is replaced by a struct literal:

  ```go
  &unicap.Error{
  	Provider: "myprovider",
  	Code:     code,
  	Message:  message,
  	Err:      unicap.ErrUnsolvable,
  }
  ```

- `Error.Retriable` is now the `Retryable()` method, and there is a package-level `unicap.IsRetryable(err)`. Retryability comes from the sentinel, so unknown codes are no longer treated as retryable.
- New sentinels:
  - `ErrUnsupportedOption`
  - `ErrUnsupportedOperation`
  - `ErrUnsolvable`
  - `ErrUnavailable`
  - `ErrProxy`
  - `ErrIPBlocked`
  - `ErrTaskFailed`
  - `ErrInvalidProxy`
  - `ErrInvalidPolling`

## Polling

`Poller`, `PollerConfig`, `NewPoller`, `WithPoller`, and `WithPollerLogger` are replaced by `unicap.Polling` and `unicap.WithPolling`:

```go
polling := unicap.DefaultPolling()
polling.Timeout = 2 * time.Minute

client, err := unicap.New(provider, unicap.WithPolling(polling))
```

Without `WithPolling`, the client uses the provider's recommended pace. A `Timeout` of zero now means "rely on the context" rather than "expire immediately".

## Proxies

- `Proxy.Address` was renamed to `Host`, and `Proxy.Login` to `Username`.
- `IsSet` is gone. A nil `*Proxy` means proxyless.
- New: `Validate`, `URL`, and `ParseProxyURL`. When you log or print a proxy, the password is redacted.

## Tasks

- **Strict options.** Options a provider can't honour now fail with `ErrUnsupportedOption` instead of being dropped. For example, a proxy on 2Captcha reCAPTCHA v3, or `MinScore` on CapSolver.
- **reCAPTCHA v3.** `ReCaptchaV3Task.IsEnterprise` was removed; use `ReCaptchaV3EnterpriseTask`. `HighScore` selects CapSolver's high-score (M1) mode.
- **reCAPTCHA v3 on 2Captcha and Anti-Captcha.** These providers now require `MinScore`.
- **Anti-Captcha proxies.** Tasks that use a proxy need `UserAgent` for reCAPTCHA v2, FunCaptcha, and GeeTest, and the proxy host must be an IP address.
- **hCaptcha.** `HCaptchaTask.EnterpriseData` was renamed to `EnterprisePayload`. The built-in providers no longer accept hCaptcha, because every provider has withdrawn it.
- **AWS WAF.** `AWSWAFTask` gained `JSAPIScript`, `ProblemURL`, `APIKey`, and `ExistingToken`. CapSolver now maps `awsApiJs` from `JSAPIScript`.

## Providers

- `capsolver.New`, `twocaptcha.New`, and `anticaptcha.New` return concrete `*Client` types, which add `Balance`, `Report`, and `Supports`.
- Each package exports a `Name` constant.
- Registry factories now receive `provider.Settings`. `Registry.New` accepts `provider.With*` options.
