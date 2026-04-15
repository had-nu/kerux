# Skill: Go Security

Mandatory reference for the Engineer role. Read before writing Go code.
Complexity is an attack vector. Simple, idiomatic Go IS secure Go.

## Pike's Rules (foundation)

1. Don't guess where time is spent. Don't add security controls for threats you can't articulate.
2. Measure. `go test -bench`, `pprof`, `govulncheck`, `gosec`, `staticcheck`.
3. Fancy algorithms are slow when n is small — and n is almost always small. Readable > clever.
4. Simple algorithms, simple data structures. The attacker benefits from your complexity.
5. Data dominates. Design starts at structs. A well-defined input struct is the first line of defence.

## Effective Go — Rules

- Accept interfaces, return concrete structs. Define interface in the consumer package.
- Errors are values. Wrap with context at system boundaries. Never `_` on security paths.
- Goroutines have owners. `context` for cancellation, `sync.WaitGroup` for join.
- Naming: `validateBearerToken` not `chkTkn`. Clarity wins.
- Comments follow godoc. Explain WHY, not WHAT.
- `any` not `interface{}` (Go 1.18+).

## Security Patterns

### Path Validation
Every path derived from user input or external data:
```go
func safePath(base, input string) (string, error) {
    candidate := filepath.Join(base, filepath.Clean(input))
    if !strings.HasPrefix(candidate, base+string(filepath.Separator)) {
        return "", fmt.Errorf("path %q escapes base dir", input)
    }
    return candidate, nil
}
```
Without `+string(filepath.Separator)`: `/home/user/dist` matches `/home/user/distributable`.

### Crypto
```go
// Randomness: always crypto/rand
import "crypto/rand"

// Comparison: always constant-time
return subtle.ConstantTimeCompare([]byte(provided), []byte(stored)) == 1
```
Never `math/rand` for security. Never `==` for tokens/hashes.

### Context Keys
```go
type ctxKey string
const claimsKey ctxKey = "claims"

func withClaims(ctx context.Context, c *Claims) context.Context {
    return context.WithValue(ctx, claimsKey, c)
}
func claimsFromCtx(ctx context.Context) (*Claims, bool) {
    c, ok := ctx.Value(claimsKey).(*Claims)
    return c, ok
}
```
Never plain `string` keys — silent collision risk.

### Secrets
```go
func mustEnv(key string) string {
    v := os.Getenv(key)
    if v == "" {
        panic(fmt.Sprintf("required env var %q is not set", key))
    }
    return v
}
```
Fail fast at startup. No hardcode. No `.env` in prod.

### Logging
```go
import "log/slog"

slog.Info("gate decision",
    "component",  "wardex",
    "decision",   decision.Action,
    "cvss",       score.CVSS,
    "epss",       score.EPSS,
    // never: token, password, secret
)
```
Structured via stdlib. No external log deps without justification. Never log PII/secrets.

### Command Injection
```go
// Never: exec.Command("sh", "-c", userInput)
// Always: separate args + whitelist
cmd := exec.CommandContext(ctx, "ping", "-c", "1", validatedHost)
```

### Integer Overflow
```go
if size < 0 || size > maxAllowed {
    return nil, errors.New("invalid size")
}
buf := make([]byte, size)
```

## Ecosystem Patterns

### CLI (Cobra)
```go
var scanCmd = &cobra.Command{
    Use:   "scan [path]",
    Short: "Scan for exposed secrets",
    Args:  cobra.ExactArgs(1),
    RunE:  runScan, // RunE, not Run — errors propagate
}

func runScan(cmd *cobra.Command, args []string) error {
    root := args[0]
    info, err := os.Stat(root)
    if err != nil {
        return fmt.Errorf("scan: invalid path: %w", err)
    }
    if !info.IsDir() {
        return fmt.Errorf("scan: %q is not a directory", root)
    }
    return scanner.Run(cmd.Context(), root)
}
```
Validate flags before any I/O.

### File Walking
```go
func walkFiles(ctx context.Context, root string, process func(string) error) error {
    var wg sync.WaitGroup
    errCh := make(chan error, 1)

    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() || !isTargetFile(d.Name()) {
            return nil
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        wg.Add(1)
        go func(p string) {
            defer wg.Done()
            if err := process(p); err != nil {
                select {
                case errCh <- err:
                default:
                }
            }
        }(path)
        return nil
    })

    wg.Wait()
    close(errCh)

    if err != nil {
        return err
    }
    return <-errCh
}
```
Owner clear. Context for cancellation. Buffered error channel. First error wins.

### JSON Envelopes (Vexil → Wardex → Vigil)
```go
type Finding struct {
    File    string  `json:"file"`
    Line    int     `json:"line"`
    Entropy float64 `json:"entropy"`
    Snippet string  `json:"snippet"`
    RuleID  string  `json:"rule_id"`
}

type ScanEnvelope struct {
    Version   string    `json:"version"`
    Timestamp time.Time `json:"timestamp"`
    Findings  []Finding `json:"findings"`
    Summary   struct {
        Total    int `json:"total"`
        Critical int `json:"critical"`
    } `json:"summary"`
}

func ParseEnvelope(r io.Reader) (*ScanEnvelope, error) {
    dec := json.NewDecoder(io.LimitReader(r, 10<<20)) // 10MB max
    dec.DisallowUnknownFields()

    var env ScanEnvelope
    if err := dec.Decode(&env); err != nil {
        return nil, fmt.Errorf("envelope: decode: %w", err)
    }
    if env.Version == "" {
        return nil, errors.New("envelope: missing version")
    }
    return &env, nil
}
```
LimitReader on input. DisallowUnknownFields. Validate version before processing.

## Checklist (pre-handoff)

- [ ] Errors handled — no `_` on security paths
- [ ] Goroutines have owner + lifecycle
- [ ] `crypto/rand` for all secure randomness
- [ ] `crypto/subtle` for security comparisons
- [ ] Context keys are typed, not strings
- [ ] No external dep without demonstrated need
- [ ] `log/slog` — no tokens/passwords/PII logged
- [ ] Secrets via `os.Getenv` — no hardcode
- [ ] `any` not `interface{}`
- [ ] `go vet ./...` clean
- [ ] `staticcheck ./...` clean
- [ ] `go test -race ./...` no data races

## Tools

```bash
govulncheck ./...        # dep vulnerabilities
gosec -fmt=json ./...    # static security analysis
staticcheck ./...        # advanced linter
go test -race ./...      # race detection
go vet ./...             # basic static analysis
```
