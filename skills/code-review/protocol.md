# Skill: Code Review Protocol

This is the Auditor's operational checklist. Every check must be performed.

## Security Audit (Sec)

| ID | Check | Pass Condition |
|----|-------|----------------|
| Sec-01 | Hardcoded secrets | No API keys, passwords, or tokens in source code. |
| Sec-02 | Input validation | All external input is validated before use. |
| Sec-03 | safePath | File paths from external input use `safePath` with separator suffix. |
| Sec-04 | Streaming IO | No `os.ReadFile` or `io.ReadAll` on unbounded files. |
| Sec-05 | Crypto stdlib | Use `crypto/rand` and `crypto/subtle`. No `math/rand` for security. |
| Sec-06 | Shell injection | No `exec.Command("sh", "-c", ...)`. Args are separated. |
| Sec-07 | Error handling | No `_` on critical error returns. Errors are wrapped. |
| Sec-08 | RunE not Run | Cobra commands use `RunE` to propagate errors. |

## Supply Chain Audit (SCA)

| ID | Check | Pass Condition |
|----|-------|----------------|
| SCA-01 | Deps pinned | All dependencies have specific versions in `go.mod`. |
| SCA-02 | go.sum consistent | `go.sum` matches `go.mod` and is updated. |
| SCA-03 | No mutable tags | No dependencies use `latest` or mutable branch tags. |

## Logic Verification (Dev)

| ID | Check | Pass Condition |
|----|-------|----------------|
| Dev-01 | Blueprint match | Implementation matches the spec Blueprint exactly. |
| Dev-02 | No TODO/FIXME | No "TODO" or "FIXME" markers left in the code. |
| Dev-03 | Tests exist | Unit or E2E tests exist for new logic. |
| Dev-04 | Exit codes | Appropriate exit codes for success and failure. |

## Environment & Ops (Env/Ops)

| ID | Check | Pass Condition |
|----|-------|----------------|
| Env-01 | Envelope compat | JSON envelopes follow the versioning/limit rules. |
| Ops-01 | Config safety | Configuration is via env vars or validated files. |
