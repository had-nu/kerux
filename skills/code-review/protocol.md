# Skill: Code Review Protocol

Auditor's operational checklist. Applied at REVIEWING state.
Audit implementation against spec_projeto.md + security baseline.

## Pre-Review

1. Read `spec_projeto.md` — the Blueprint and Guardrails sections are the audit scope.
2. Read the Engineer's handoff packet — note which files were created/modified.
3. Verify every file listed in the Blueprint has been addressed. Missing file = automatic REJECT.

## Security Audit (Sec)

- [ ] No hardcoded secrets (Commandment C2). Grep for: API keys, passwords, tokens in source.
- [ ] Input validation at every trust boundary (CLI flags, file input, manifest parsing, env vars).
- [ ] Path traversal protection: `safePath()` with separator suffix on all paths from external input.
      Check: `strings.HasPrefix(candidate, base+string(filepath.Separator))` — NOT just `HasPrefix(candidate, base)`.
- [ ] Streaming I/O: no `os.ReadFile` or `io.ReadAll` on files of unbounded size. Must use `io.Copy`.
- [ ] Crypto: `crypto/rand` for randomness, `crypto/subtle` for comparisons. Grep for `math/rand` and `==` on tokens.
- [ ] No shell injection: no `exec.Command("sh", "-c", userInput)`. Args must be separate.
- [ ] Error handling: no `_` on error returns in security-critical paths. Grep for `_ =` and `_ :=`.
- [ ] Cobra uses `RunE` not `Run`.

## Supply Chain Audit (SCA)

- [ ] No new unsigned or unpinned dependencies added without spec justification.
- [ ] `go.sum` consistent with `go.mod` (`go mod verify` exits 0).
- [ ] No mutable tag references in CI workflows (`@latest`, `@main`). Must pin SHA or version.
- [ ] If project has goreleaser/cosign: verify signing config intact.

## Logic Verification (Dev)

- [ ] Implementation matches spec Blueprint section exactly. Every [NEW] file exists, every [MODIFY] applied.
- [ ] No TODO/FIXME left as implementation (Commandment C3).
- [ ] Tests cover changed logic paths. Missing test for a security-relevant path = REJECT.
- [ ] Exit codes match spec (if CLI). Test or verify manually.
- [ ] Documentation preserved. Existing godoc comments not removed unless spec says so.

## Envelope Compatibility (if project produces/consumes JSON envelopes)

- [ ] Struct field additions are backward-compatible (new fields with `omitempty`).
- [ ] No field renamed or removed without version bump in envelope.
- [ ] ParseEnvelope-pattern validator updated to handle new fields.
- [ ] `io.LimitReader` + `DisallowUnknownFields` present on decode path.

## Ops Audit

- [ ] Config changes are logged and safe.
- [ ] No destructive operations without user confirmation path.
- [ ] CI Mirror commands from spec (`go vet`, `staticcheck`, `gosec`, `go test -race`) would pass.

## Verdicts

PASS — all checks pass. Emit: `→K|REV→STG|all checks pass, 0 findings|stage for commit`

COMMENT — minor improvements requested, functionality safe. Emit: `→K|REV→STG|COMMENT: {notes}|stage with notes`

REJECT — must specify:
1. Which check failed (S#, SCA#, or spec section).
2. File and line number.
3. Target state for the fix:
   - `→E|REV→IMP|REJECT: {file} L{line} {finding}|fix and re-submit` — code-level fix, spec is correct.
   - `→A|REV→DES|REJECT: {finding}|spec amendment required` — design flaw or security finding that needs spec change.

## Rules

- Never edit files directly. Auditor reads, never writes.
- Every rejection must cite a specific file, line, and check ID.
- When in doubt about a finding, REJECT with explanation. False positives are cheaper than missed vulnerabilities.
- The safePath separator check is a mandatory verification on ANY project that handles file paths from external input. If safePath exists without the separator suffix, this is an automatic REJECT.
