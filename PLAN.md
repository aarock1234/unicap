# unicap: audit and implementation plan

Status: implemented in v3 (October 2026). See [MIGRATION.md](MIGRATION.md) for the resulting API changes.

Final decisions:

- **Go version:** Go 1.27.2, matching the template.
- **Module path:** `github.com/aarock1234/unicap/v3`. No retract directive was added: the old v2 tags sit under a different module path and the toolchain never resolves them.
- **hCaptcha:** `HCaptchaTask` is kept. The built-in providers reject it with `ErrUnsupportedTask`.
- **Unsupported options:** return `ErrUnsupportedOption`.
- **Task coverage:** P1, P2 and P3 are all typed.

Sources:

- The current tree at `0302839`.
- The Go style guide (`~/.agents/skills/go-style`).
- The reference project `~/Documents/Coding/go-template`.
- The live provider docs, crawled October 2026:
  - [CapSolver](https://docs.capsolver.com)
  - [2Captcha API v2](https://2captcha.com/api-docs)
  - [Anti-Captcha](https://anti-captcha.com/apidoc)

Baseline: `go build`, `go vet`, `gofmt -l` and `golangci-lint` v2.14.0 report no issues today. Every finding below is about design, correctness against the real provider APIs, or missing features. None of them is a lint failure.

---

## 1. Critical findings

These need to be fixed first. Each one silently breaks real users today.

| #   | Finding | Impact |
| --- | ------- | ------ |
| C1  | **Module path and tags disagree.** The tags go up to `v2.1.0`, but `go.mod` declares `github.com/aarock1234/unicap` with no `/v2` suffix. | `go get …@v2.1.0` is rejected by the Go toolchain, so `@latest` resolves to `v1.0.3`. Nobody can install anything from the v2 series. |
| C2  | **CapSolver returns errors as HTTP 400/401 with a JSON body.** `solverapi.doJSON` gives up on any non-200 status before decoding. | Every CapSolver error, including bad key, zero balance and invalid task, becomes `unexpected status code: 400` and loses its code. |
| C3  | **The CapSolver error table uses codes that don't exist** (`ERROR_KEY_INVALID`, `ERROR_KEY_DOES_NOT_EXIST`, `ERROR_NO_SLOT_AVAILABLE`, `ERROR_TASK_NOT_FOUND`). | No CapSolver error ever matches a sentinel. The real codes are `ERROR_KEY_DENIED_ACCESS`, `ERROR_TASKID_INVALID`, `ERROR_TASK_NOT_SUPPORTED` and so on. |
| C4  | **The 2Captcha and Anti-Captcha error tables are wrong.** `ERROR_TASK_ABSENT` means "the `task` field is missing", not "task not found" (that is `ERROR_NO_SUCH_CAPCHA_ID`). `ERROR_NO_SLOT_AVAILABLE` is a retryable capacity error, not insufficient funds. `ERROR_WRONG_TASK_DATA` and `ERROR_WRONG_USER_KEY` are not documented. | Callers get the wrong sentinels, and a retryable condition is reported as fatal. |
| C5  | **`getTaskResult` sends `taskId` as a JSON string.** 2Captcha and Anti-Captcha document it as an integer. | This only works if the providers coerce the type. That behaviour is undocumented and could stop at any time. |
| C6  | **CapSolver recognition tasks finish synchronously.** `createTask` returns `status: "ready"` with the `solution`, and the client throws both away and starts polling. | `ImageToTextTask` and every classification task wastes requests and may fail, because polling a task that already finished is not guaranteed to work. |
| C7  | **Some type strings point at tasks that no longer exist.** CapSolver removed hCaptcha and FunCaptcha. CapSolver Turnstile and GeeTest have no proxy variant (`AntiTurnstileTask` and `GeeTestTask` are invented). 2Captcha and Anti-Captcha removed hCaptcha. 2Captcha and Anti-Captcha have no `RecaptchaV3EnterpriseTaskProxyless`: v3 Enterprise is `RecaptchaV3TaskProxyless` with `isEnterprise: true`. Neither has a v3 proxy variant. | Users get confusing provider errors (or charges) instead of a clear `ErrUnsupportedTask` up front. |
| C8  | **Anti-Captcha proxy tasks require `userAgent`** (reCAPTCHA v2, v2 Enterprise, FunCaptcha and GeeTest). The mapper never sends it, and `ReCaptchaV2EnterpriseTask` has no `UserAgent` field at all. | Every Anti-Captcha task that uses a proxy is rejected. |
| C9  | **`solverapi.ProxyString` always emits the login and password segments.** With no credentials it produces `http:1.2.3.4:8080::`. | IP-whitelisted proxies break on CapSolver. |
| C10 | **Options are silently dropped.** Anti-Captcha v3 ignores `Proxy` and `IsEnterprise`. CapSolver v3 ignores `IsEnterprise` and `MinScore`. 2Captcha v2 ignores `PageAction`. The 2Captcha v2 Enterprise wire struct has `UserAgent` and `Cookies` fields, but nothing ever fills them. | Requests behave differently from what the caller asked for, and nothing reports it. |
| C11 | **`Solution.Coordinates` is never filled in** by any provider. GeeTest, AWS WAF, Cloudflare `cf_clearance`, the `userAgent` field and every cookie map are only reachable through `Extra map[string]any`. | The typed API promises data it never delivers. |

---

## 2. Style and structure sweep (compared with the template and style guide)

Each finding links to the guide section it comes from.

### Module, tooling and repository

| Area | Current | Target (template) |
| ---- | ------- | ----------------- |
| Go baseline | `go 1.26` | The template uses `1.27.2`. **Decision D1.** |
| Makefile | No `help`, unpinned `golangci-lint@latest`, `format` runs `go fmt`, `check` also runs tests | Copy the template: `help` with `##` docs, a pinned `GOLANGCI_LINT_VERSION`, `format` running `golangci-lint fmt`, and separate `lint`, `fix` and `build` targets |
| `.golangci.yml` | Has errorlint, gocritic, misspell and revive | Add `nilerr` and `sloglint` (`attr-only`, `context: scope`, `key-naming-case: snake`, `msg-style: lowercased`) as the template does |
| CI | `checkout@v6` and `setup-go@v6`, lint `v2.11.1`, no `permissions`, no tidy check | `@v7` actions, `permissions: contents: read`, `go mod tidy -diff`, lint `v2.14.0` pinned in step with the Makefile, and `workflow_call` so the release workflow reuses CI instead of copying it |
| Dependencies | `github.com/joho/godotenv` is a required dependency of the **library** module, only because one example uses it | Remove it. The Makefile already includes `.env` |
| Leftover directories | Empty untracked `providers/{anticaptcha,capsolver,twocaptcha}`, `unicap/tasks` and `internal/logger` | Delete them |
| Examples | `cmd/example` and `examples/capsolver` overlap. `cmd/example` uses `http.DefaultClient` (no timeout), `time.Sleep` and `slog.Info` without a context | One `examples/` tree. Each example shows one feature, uses `run()`, a bounded client and context-aware logging |
| README | Broken badges (`pkg.go.dev/unicap`, `https://unicap/releases`), install instructions that resolve to v1, and a support matrix that no longer matches the providers | Regenerate it from the capability tables (see §4.6) |

### Core package (`unicap`)

1. **Long positional constructors.** `NewError(code, message, provider, retriable, wrapped)` takes five positional arguments, two of them strings that are easy to swap. Replace it with exported fields and a keyed struct literal, as the template's `apperr.Error` does (go-style `structs.md`: record layout).
2. **Inverted retry logic.** `ErrorMapper.Error` marks every *unknown* code as retriable and every *known* code as not retriable. That gets `ERROR_NO_SLOT_AVAILABLE` and `ERROR_CAPTCHA_UNSOLVABLE` wrong. The poller also ignores `Retriable` completely: it retries `ErrInvalidAPIKey` three times. Classify each code explicitly (see §4.2).
3. **`Client.Solve` duplicates `Client.CreateTask`**, including validation, creation and logging. `Solve` should call `CreateTask`.
4. **`Solve` returns only `*Solution`.** The task ID, cost and solve time are lost, so callers can't report a wrong answer or account for spend. Return a `*Result` instead (§4.1).
5. **The poller duplicates its dependency.** `NewPoller(provider, config, opts...)` holds its own copy of the provider, so `WithPoller` means building the provider twice. The config isn't validated: `Timeout: 0` expires immediately, and `Multiplier < 1` or `InitialInterval <= 0` causes a busy-loop. Compare with the template's `retry` package, which validates its config (go-style `structs.md`: functional options).
6. **Generic failure sentinel.** A failed status with no error becomes `ErrInvalidTask`, which is misleading. Add `ErrTaskFailed`.
7. **Discard logger.** `slog.New(slog.NewTextHandler(io.Discard, nil))` is written out twice. Go 1.24+ has `slog.DiscardHandler`.
8. **`Solution.Extra map[string]any`** breaks the rule about not decoding known shapes into untyped maps (go-style `types.md`: JSON boundaries). Keep a raw payload as `json.RawMessage` and decode the known shapes into typed fields.
9. **Proxy type gaps.** `Proxy` has no validation (port range, `ProxyType` membership), no URL parsing, and no redaction. Following the template's `client/proxy.go`, add `Validate`, `URL()`, `ParseProxyURL` and a `LogValue` that hides the password.
10. **Undocumented exported fields.** Most exported struct fields, especially in `tasks/`, have no doc comment (go-style `modules.md`: doc comments).

### `internal/solverapi`

1. **Positional constructor.** `New(name, baseURL, apiKey, mapper, errs, opts...)` takes six positional arguments. Replace it with a `Config` struct.
2. **Opaque error table builder.** `StandardErrorMapper(name, []string, []string, []string, []string)` takes four anonymous slices. Replace it with a keyed `map[string]ErrorRule` literal that each provider can read top to bottom.
3. **HTTP handling problems:**
   - The response body is read without a size limit.
   - The full response body is logged at debug level, including tokens and cookies.
   - A non-200 status throws away a decodable error body (C2).
   - A 5xx or 429 on `getTaskResult` isn't marked transient.
4. **No task ID check.** An empty `taskId` with `errorId: 0` is accepted. Reject it at the boundary.
5. **Lenient decoding only.** `TaskID` decodes from a string or a number, but there is no symmetric encoding back to the provider's wire type (C5). The same mismatch exists for `cost` and `solveCount`, which arrive as strings or numbers.

### Provider packages

1. **Constructors return the interface.** `New` returns `unicap.Provider`. Returning the concrete `*capsolver.Client` (accept interfaces, return structs) lets each provider expose `Balance` and `Report` without type assertions (go-style `interfaces.md`: consumer interfaces).
2. **Mapper files are 325 to 548 lines**, with wire structs and mapping mixed together, and the same 6-line proxyless/proxy switch repeated about 15 times per provider. Split them by family (`recaptcha.go`, `turnstile.go`, …) and replace the giant `switch` with a registered table. That table also answers "is this task supported?" (§4.6).
3. **Identical option wrappers in all three packages.** That is acceptable for a public API, but `WithSoftID`/`WithAppID`, `WithCallbackURL` and `WithPollDefaults` are missing.
4. **Registry factories only take an API key**, so options can't be passed through. Provider names are bare strings; export `capsolver.Name` and the equivalents.

---

## 3. Feature gap map

Legend: **✓** = supported today. **✗** = the provider supports it but unicap does not. **⚠** = unicap maps it, but the mapping is wrong or obsolete. **–** = the provider doesn't offer it.

### 3.1 Platform features

| Feature | CapSolver | 2Captcha | Anti-Captcha | unicap today |
| ------- | --------- | -------- | ------------ | ------------ |
| Synchronous result from `createTask` | ✓ recognition tasks | – | – | ✗ ignored (C6) |
| `getBalance` | ✓ | ✓ | ✓ | ✗ |
| Report a correct or incorrect answer | `feedbackTask` (needs the original task and result) | `reportCorrect` / `reportIncorrect` | `reportCorrectRecaptcha`, `reportIncorrectRecaptcha`, `reportIncorrectImageCaptcha` | ✗ |
| Affiliate ID | `appId` (string) | `softId` (int) | `softId` (int) | ✗ |
| `callbackUrl` webhooks | ✓ (format not documented) | ✓ (form-encoded `id`, `code`) | ✓ (same body as `getTaskResult`) | ✗ |
| `languagePool` | – | ✓ (top-level) | ✓ (`ImageToTextTask` field) | Anti-Captcha only, on the wrong side |
| Cost and timing in the result | – | `cost`, `createTime`, `endTime`, `solveCount` | same | ✗ |
| Recommended poll pacing | about 3 s, 120 polls max, 5 min TTL | 1–2 s first wait, then 5 s | about 5 s first wait, then 2 s, result TTL 60 s | One global default |
| Proxy formats | Separate fields **or** a `proxy` string / URL | Separate fields | Separate fields (IP addresses only, no hostnames) | ⚠ C9 |

### 3.2 Task coverage (current plus new)

**Existing tasks: required fixes**

| Task | CapSolver | 2Captcha | Anti-Captcha |
| ---- | --------- | -------- | ------------ |
| reCAPTCHA v2 / v2 Enterprise | ✓. Add `cookies`. The solution also carries `userAgent`, `secChUa` and `recaptcha-ca-*` | ✓. Add `pageAction`; send `userAgent`/`cookies` for Enterprise | ⚠ C8. Send `userAgent`/`cookies` |
| reCAPTCHA v3 / v3 Enterprise | ⚠ `IsEnterprise` is ignored. **New:** M1 high-score types | ⚠ C7. Enterprise means `isEnterprise`; there is no proxy variant | ⚠ C7 and C10 |
| hCaptcha | ⚠ removed | ⚠ removed | ⚠ removed. **Deprecate the task (D3)** |
| FunCaptcha | ⚠ removed | ✓ | ✓ (send `userAgent` with a proxy) |
| Turnstile | ⚠ no proxy variant | ✓ | ✓ |
| Cloudflare Challenge | ✓. Decode `cookies.cf_clearance` into a typed field | – (handled by Turnstile with `pagedata`) | – |
| GeeTest v3/v4 | ⚠ proxyless only. **New:** v4 `riskType` | ✓. Add `riskType` | ✓ |
| DataDome | ✓. Add optional `html` | ✓ | – |
| AWS WAF | ✓. **New:** `awsProblemUrl`, `awsApiKey`, `awsExistingToken` | ✓. **New:** `jsapiScript` form | ✗ **New:** `AmazonTask(Proxyless)` |
| MTCaptcha | ✓ | ✓ | – |
| Friendly Captcha | – | ✓. **New:** `version` v2, module scripts | ✓ |
| Prosopo | – | ✓ | ✓ |
| Altcha | – | ✓ | ✗ **New** |
| Lemin, Cutcaptcha, Text | – | ✓ | – |
| ImageToText | ⚠ C6. **New:** `images` / `number` module | ✓ | ✓. Note: `numeric: 2` means "no digits" here, and `case` defaults to true |

**New task types, in priority order**

The order is driven by how many providers offer each task, then by demand.

| Priority | Task (proposed unicap name) | Providers and wire types |
| -------- | --------------------------- | ------------------------ |
| P1 | `CoordinatesTask` (click points or rectangles on an image) | 2Captcha `CoordinatesTask`, Anti-Captcha `ImageToCoordinatesTask` |
| P1 | `GridTask` (select matching tiles) | 2Captcha `GridTask` (`imgType: recaptcha/funcaptcha`), CapSolver `ReCaptchaV2Classification` |
| P1 | reCAPTCHA v3 high score | CapSolver `ReCaptchaV3M1Task*` (as an option on the v3 task) |
| P1 | `AntiGateTask` (typed, with `PushVariable`) | Anti-Captcha |
| P2 | `RotateTask`, `AudioTask`, `BoundingBoxTask`, `DrawAroundTask`, `DragAndDropTask` | 2Captcha |
| P2 | `CapyTask`, `KeyCaptchaTask`, `TencentTask`, `CaptchaFoxTask`, `YandexSmartCaptchaTask`, `CyberSiAraTask`, `AtbCaptchaTask` | 2Captcha |
| P2 | `VisionEngineTask`, `AWSWAFClassificationTask`, `BotDeflectorTask` | CapSolver |
| P3 | `VKCaptchaTask`, `TemuTask`, `YidunTask`, `BinanceTask`, `HuntTask`, `BasiliskTask`, `AlibabaTask`, `TSPDTask`, `IncapsulaTask` | 2Captcha (all added in 2025–2026; several are documented inconsistently) |

P3 tasks stay reachable through `RawTask` until someone needs a typed version.

---

## 4. Target design

### 4.1 Results instead of bare solutions

```go
// Result describes a finished task.
type Result struct {
	TaskID   string
	Solution Solution
	Cost     Cost          // zero when the provider does not report it
	Elapsed  time.Duration
}
```

`Client.Solve` returns `*Result`. `Solution` keeps the common typed fields:

- `Token`
- `Text`
- `UserAgent`
- `Cookies []*http.Cookie` (decoded from both cookie strings and cookie maps)
- `Coordinates`

It also gets `Raw json.RawMessage`. Family-specific shapes get typed decoders, for example `solution.GeeTestV4()`, `solution.AWSWAF()` and `solution.Grid()`. Each decoder returns a struct and an error, instead of a `map[string]any` the caller has to pick apart.

### 4.2 Errors that say whether to retry

```go
type Error struct {
	Provider  string
	Code      string
	Message   string
	Retryable bool
	Err       error // sentinel, if any
}
```

Each provider writes one table that reads top to bottom:

```go
var errorRules = map[string]solverapi.ErrorRule{
	"ERROR_KEY_DENIED_ACCESS": {
		Err: unicap.ErrInvalidAPIKey,
	},
	"ERROR_CAPTCHA_UNSOLVABLE": {
		Err:       unicap.ErrUnsolvable,
		Retryable: true,
	},
	// ...
}
```

Unknown codes default to **not retryable**. New sentinels:

- `ErrUnsolvable`
- `ErrNoCapacity` (no slot available or rate limited; retryable)
- `ErrProxy` (bad, banned or unreachable proxy)
- `ErrIPBlocked`
- `ErrTaskFailed`

### 4.3 Polling and retries

- The poller becomes a value-only configuration. The client passes the provider at call time, so the provider is no longer stored twice.
- The config is validated once at `New`.
- Each provider supplies default pacing: an initial delay, an interval and a TTL (see §3.1). Callers override it with `WithPolling(...)`.
- Transient fetch errors still use a budget, but only errors that are actually retryable count against it. Errors that are not retryable return immediately.
- A new opt-in `WithSolveRetries(n)` creates a fresh task when the error is `ErrUnsolvable` or `ErrNoCapacity`. A failed attempt can still cost money, so the docs will say so clearly.

### 4.4 Optional capabilities as small consumer interfaces

```go
type BalanceChecker interface {
	Balance(ctx context.Context) (Balance, error)
}

type Reporter interface {
	Report(ctx context.Context, result *Result, correct bool) error
}
```

The client exposes `client.Balance(ctx)` and `client.Report(ctx, result, correct)`. When the provider doesn't implement a capability, they return `ErrUnsupported`. CapSolver's `feedbackTask` needs the original task, so `Result` keeps an unexported reference to it.

### 4.5 `solverapi` HTTP boundary

- Use a `Config` struct instead of positional arguments.
- Decode the JSON envelope for **any** status code when the body is JSON. Use the HTTP status only when there is no envelope.
- Limit the response size with `io.LimitReader`.
- Log only the endpoint, status and duration, never the request or response body.
- Treat 5xx and 429 responses on `getTaskResult` as transient.
- Encode `taskId` in the provider's wire type: a number for 2Captcha and Anti-Captcha, a string for CapSolver.
- When `createTask` returns `status: "ready"` with a solution, short-circuit and return it immediately.
- Decode `cost` and `solveCount` from either a string or a number, with a small typed decoder like `TaskID`.

### 4.6 Provider task tables

Each provider registers its supported tasks in a table instead of a type switch:

```go
var mappers = solverapi.Mappers{
	unicap.TaskTypeReCaptchaV2: solverapi.Map(mapReCaptchaV2),
	unicap.TaskTypeTurnstile:   solverapi.Map(mapTurnstile),
	// ...
}
```

`solverapi.Map[T]` is a small generic adapter. The table has several uses:

- It answers `Supports(unicap.TaskType) bool`.
- It powers a `Capabilities()` method.
- It is the single source for the README support matrix.

The choice between the proxy and proxyless type string goes through one named helper, `taskType(base, proxy)`, instead of 15 hand-written copies.

**Validation:** when a provider can't honour a set option (a proxy on 2Captcha v3, `MinScore` on CapSolver), the mapper returns `ErrUnsupportedOption` wrapped with the field name instead of dropping it silently. **Decision D4.**

### 4.7 Proxy

- `Proxy.Validate()` checks the port range and the `ProxyType` value. Anti-Captcha's mapper additionally rejects hostnames, since that API only accepts IP addresses.
- `ParseProxyURL("socks5://user:pass@host:port")` parses a proxy URL.
- `URL()` returns the proxy as a `*url.URL`.
- `LogValue` redacts the password.
- `ProxyString` leaves out empty credentials (C9).

---

## 5. Phased implementation

Each phase leaves the module building, vetted, formatted and lint-clean.

### Phase 0: Decisions (needs your input)

- **D1. Go baseline.** Keep `1.26` so library consumers can use older toolchains, or move to the template's `1.27`, which allows generic methods such as a typed `doJSON[T]`. *Recommendation: 1.26, because a library should not force a toolchain upgrade. Nothing below requires 1.27.*
- **D2. Module path.** The v2 tags are unusable (C1). Release the redesign as `github.com/aarock1234/unicap/v3`, skipping the broken v2 line, and retract `v2.0.0`–`v2.1.0` in `go.mod`. *Recommendation: yes.*
- **D3. hCaptcha.** All three providers have removed it. Either delete `HCaptchaTask`, or keep it marked `Deprecated` while every mapper returns `ErrUnsupportedTask`. *Recommendation: delete it in v3.*
- **D4. Unsupported options.** Return an error (strict) or silently drop the option (lenient)? *Recommendation: strict by default.*
- **D5. P3 tasks.** Type them now, or leave them to `RawTask`? *Recommendation: leave them to `RawTask`.*

### Phase 1: Repository and tooling hygiene

- Delete the empty leftover directories.
- Remove `godotenv`, then run `go mod tidy`.
- Update the Makefile, `.golangci.yml` (`nilerr`, `sloglint`) and the CI and release workflows to the template's conventions.
- Merge `cmd/example` into `examples/`.

### Phase 2: Correctness fixes (C2–C11)

These are wire-level fixes inside the current API surface, so they can also ship as a `v1.1.0` patch if wanted.

- Decode the error envelope on non-200 responses.
- Correct the three error tables.
- Encode `taskId` as an integer for 2Captcha and Anti-Captcha.
- Short-circuit synchronous `createTask` results.
- Fix the type strings: remove dead types, fix v3 Enterprise, remove invented proxy variants.
- Send `userAgent` and `cookies` to Anti-Captcha.
- Fix `ProxyString`.
- Stop dropping fields silently.

### Phase 3: Core API redesign (v3)

- `Result`, a typed `Solution` with `Raw` and family decoders, the `Error` struct with retry rules, and the new sentinels.
- Poller rewrite: validated config, provider pacing defaults, retry-aware error budget.
- `solverapi.Config`, the hardened HTTP boundary and the `Mappers` table.
- Concrete provider return types and `Name` constants.
- Proxy validation, parsing and redaction.
- Doc comments on every exported field.

### Phase 4: Platform features

- `Balance` on all three providers.
- `Report` on all three (CapSolver through `feedbackTask`).
- `WithSoftID`/`WithAppID`, `WithCallbackURL` and `WithLanguagePool`.
- Opt-in `WithSolveRetries`.
- Registry factories that accept options.

### Phase 5: Task coverage

- Bring the existing tasks up to date with the field additions in §3.2.
- Add the P1 tasks, then P2.
- Regenerate the README matrix from the provider tables.

### Phase 6: Documentation

- Rewrite the README for v3: install path, quick start, `Result`, errors and retries, capabilities, and a custom-provider guide.
- Write a short migration note from v1/v2 to v3.

---

## 6. Open questions in the provider docs

These are things the docs contradict or don't say. The implementation should handle them safely.

- **CapSolver casing.** CapSolver spells type names inconsistently (`GeeTestTaskProxyLess` / `GeetestTaskProxyless`, `DatadomeSliderTask` / `DataDomeSliderTask`). Use the spelling from each page's type table.
- **CapSolver callback format.** The `callbackUrl` payload is not documented. Expose the option, but don't parse callbacks yet.
- **2Captcha hCaptcha.** 2Captcha's hCaptcha removal is inferred: the page now returns 404 and the maintainer says so in [Furry/2captcha#55](https://github.com/Furry/2captcha/issues/55). There is no official notice.
- **2Captcha Temu type name.** The Temu page contradicts itself (`TemuImageTask` vs `TemuCaptchaTask`). Leave it to `RawTask` until this is settled.
- **Anti-Captcha numeric fields.** Anti-Captcha returns `cost` and `solveCount` as strings even though the docs call them numbers. Decode both forms.
- **Anti-Captcha proxy credentials.** Whether `proxyLogin` and `proxyPassword` are required on the newer proxy task types is unclear. Treat them as optional.
