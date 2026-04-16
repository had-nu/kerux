# Edicts

Authoritative guidance. Deviations require justification in the Decision Log of the active spec.

## E1: Architecture Style

- Prefer composition over inheritance.
- Keep components focused. Single responsibility per package.
- Define interfaces in the consumer, not the producer.
- Structs first, functions follow (Pike's Rule 5).
- When modifying existing code, preserve the existing style unless the spec directs a refactor.

## E2: Commit Discipline

- Use Conventional Commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`, `perf:`, `security:`.
- Scope optional: `feat(wardex): add policy file loader`.
- Body references spec version and related tasks: `Implements spec v1.2 §3.1`.
- Breaking changes marked with `!` and explained in body: `refactor!: change envelope schema`.

## E3: File Management

- Group related files in domain directories (`internal/scanner/`, `cmd/verify/`).
- Avoid god files handling multiple unrelated responsibilities.
- Split when a file exceeds ~400 lines or covers more than one domain concept.
- Test files colocated with implementation: `foo.go` + `foo_test.go` in same directory.

## E4: Review Loop

- Non-trivial code changes always go through the Auditor.
- Trivial changes (typo, doc comment): still passes through REVIEWING, but the Auditor can PASS quickly.
- Never route around the Auditor to "save time" — the loop is the quality mechanism.
- Reviewer REJECT triggers mandatory redesign or fix phase, never silent override.

## E5: Dependency Stance

- Stdlib first. External deps require Architect justification in spec's §3.4 External Dependencies.
- Pin versions explicitly in `go.mod`. No `latest`, no floating tags.
- New dep = supply chain risk. Evaluate: maintainer, age, license, vulnerability history.
- Prefer well-established libs (cobra, slog, testify) over novel ones.

## E6: Error Handling

- Errors are values. Return them. Don't panic except at startup for invariant violations (mustEnv pattern).
- Wrap with context at system boundaries: `fmt.Errorf("scan: invalid path: %w", err)`.
- Never `_` on error in a security-critical path.
- Errors bubble up to the caller that has context to handle them. Not before.

## E7: Logging

- Use `log/slog` structured. Never `fmt.Println` for operational logs.
- Log decisions with context: component, action, inputs, outcome.
- Never log secrets, tokens, passwords, PII.
- Log levels: DEBUG for traces, INFO for decisions, WARN for recoverable anomalies, ERROR for failures requiring attention.

## E8: Testing Discipline

- Unit tests for logic paths. E2E for user-observable behaviour.
- Tests colocated with code. `go test ./...` runs all.
- `go test -race` required before handoff from Engineer.
- Table-driven tests preferred for multiple input variants.
- Mocks only when stdlib + real implementation is impractical (filesystem, network). Prefer interface injection.
