# {Project Name} — Technical Specification
<!-- Version: X.Y | Status: {Draft|Active|Superseded} | Author: {name} | Date: YYYY-MM-DD -->

> **RFC 2119 Convention**: The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHOULD", "SHOULD NOT", "MAY", and "OPTIONAL" are interpreted as in [RFC 2119](https://www.ietf.org/rfc/rfc2119.txt).

> **Kerux context**: This spec is the single source of truth for the Engineer's implementation and the Auditor's verification. Neither role improvises beyond it.

---

## 1. Overview

### 1.1 Problem Statement
{What exists today. What is broken or missing. Why it matters. One paragraph, concrete.}

### 1.2 Proposed Solution
{One paragraph. What this project does. Who uses it. How it resolves §1.1.}

### 1.3 Scope

**IN:**
- {Capability 1}
- {Capability 2}

**OUT:**
- {Explicit non-goal 1}
- {Explicit non-goal 2}

---

## 2. Goals & Non-Goals

### 2.1 Goals
1. {Measurable outcome with a verification method}
2. {Measurable outcome}
3. Pass all automated checks: `go vet`, `staticcheck`, `gosec`, `go test -race`.

### 2.2 Non-Goals
1. {Explicit exclusion}
2. {Explicit exclusion}

---

## 3. Architecture

### 3.1 System Diagram

```
{text-based component diagram showing data flow}
```

### 3.2 Component Inventory

| Component | Responsibility | Technology | Notes |
|-----------|---------------|------------|-------|
| `{path/to/component}` | {one-line responsibility} | {lib or stdlib package} | {constraint or integration note} |

### 3.3 Data Flow

**{Operation name} flow:**
1. {Step}
2. {Step}
3. {Step}

### 3.4 External Dependencies

| Dependency | Version | Purpose | Risk if Unavailable |
|-----------|---------|---------|---------------------|
| `{import path}` | {version} | {why needed} | {impact} |
| Go stdlib ({packages}) | go{version}+ | {core use} | N/A |

---

## 4. Data Model

### 4.1 Core Entities

```go
// {package path}

type {Name} struct {
    {Field} {Type} `{tags}` // {comment}
}
```

### 4.2 Storage
{Plain text / SQLite / in-memory / N/A. Include format if text.}

### 4.3 Data Lifecycle
{When created. When modified. When purged. Retention policy or "user manages."}

---

## 5. Interfaces

### 5.1 CLI / API Surface

```
{command} {subcommand} [flags]
```

| Command | Flag | Type | Required | Default | Description |
|---------|------|------|----------|---------|-------------|
| `{sub}` | `--{flag}` | {type} | {yes/no} | {default or —} | {description} |

### 5.2 Configuration
{Config files, env vars — or "None. All inputs via CLI flags."}

### 5.3 Output Formats

**{Command}**: {stdout format, file format, or "No stdout on success."}

**Exit codes**: 0 ({success}), 1 ({expected failure}), 2 ({operational error}).

---

## 6. Performance & Capacity

### 6.1 Targets

| Metric | Target | Boundary Condition |
|--------|--------|---------------------|
| {Metric} | {value with unit} | {measurement context} |

### 6.2 Bottlenecks & Limits
{Known constraints. Explicit non-goals re: concurrency, scaling.}

### 6.3 Scaling Strategy
{Not applicable / horizontal / vertical — with rationale.}

---

## 7. Security & Compliance

### 7.1 Threat Model

Every trust boundary has an entry. No exceptions.

| ID | Threat | Attack Vector | Control |
|----|--------|--------------|---------|
| T1 | {threat name} | {how it's exploited} | {named control — reference go-security skill pattern if applicable} |

Canonical controls for common surfaces (reference the go-security skill):
- File path from external input → `safePath()` with separator suffix
- File reads of unbounded size → streaming via `io.Copy` into hash
- Symlinks in walk → skip via `d.Type()` check
- User input to shell → `exec.CommandContext` with separated args
- Token comparison → `crypto/subtle.ConstantTimeCompare`
- Randomness → `crypto/rand`
- Secrets → `os.Getenv` with startup validation

### 7.2 Data Handling
{Secrets: none / encrypted / env var. PII: none / how handled. Network: none / TLS / mTLS.}

### 7.3 Compliance Requirements
{None / ISO 27001 / DORA / internal-only / specify framework and scope.}

---

## 8. Deployment & Operations

### 8.1 Infrastructure
{Local binary / container / library. Build command.}

### 8.2 Build & Release

```bash
{build commands}
```

{CI pipeline reference or "Manual build for this scope."}

### 8.3 Monitoring & Observability
{Logs via slog / metrics exported / none for CLI.}

### 8.4 Disaster Recovery
{Not applicable / backup strategy / rollback procedure.}

---

## 9. Testing Strategy

### 9.1 Unit Tests

| Package | Test | Description |
|---------|------|-------------|
| `{pkg}` | `Test{Name}` | {what it verifies} |

### 9.2 Integration Tests
{Table with same shape, or "Not applicable — unit + E2E cover the surface."}

### 9.3 End-to-End / Acceptance Tests

| ID | Test | Steps | Expected |
|----|------|-------|----------|
| F1 | {test name} | {1. ... 2. ... 3. ...} | {observable outcome + exit code} |

### 9.4 Performance / Load Tests
{Not applicable / benchmark targets.}

---

## 10. Milestones & Deliverables

| Phase | Deliverable | Success Criteria | Target Date |
|-------|------------|------------------|-------------|
| 1 | {artifact} | {verifiable condition} | {date or —} |

---

## 11. Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| {risk} | {impact description} | {Low/Med/High} | {mitigation or "Accepted as known limitation"} |

---

## 12. Decision Log

Every non-obvious choice is logged. Assumptions from elicitation phase live here with status "Assumed — pending user confirmation."

| ID | Decision | Rationale | Date | Status |
|----|----------|-----------|------|--------|
| D-001 | {decision} | {why} | YYYY-MM-DD | {Accepted/Rejected/Superseded/Assumed} |

---

## 13. Open Questions

- [ ] {Question that remains after elicitation, documented for user follow-up.}
- [x] {Resolved question with resolution noted.}

---

## Appendices

### A. Glossary

| Term | Definition |
|------|-----------|
| {Term} | {Definition} |

### B. References

- {External spec, RFC, or standard}
- {Internal spec dependency}

---

## Blueprint

The executable section. The Engineer implements this literally. The Auditor verifies against this.

### Files to create / modify / delete

#### `[NEW]` {relative/path/to/file.go}
{One sentence describing what the file does.}

```pseudocode
{Precise pseudocode. Not English paragraphs. Structure matches the intended code.

Reference security patterns by name, do not reinvent:
- "Use safePath() from go-security skill" — not "validate the path"
- "Use io.Copy streaming hash" — not "hash the file"

The Engineer follows this literally. If it's ambiguous here, it's ambiguous in production.}
```

#### `[MODIFY]` {relative/path/to/file.go}
{One sentence describing the change.}

**Anchor**: `{file}:L{start}-L{end}`

**Before** (current code at anchor, verified by Analyst to match the live file):
~~~go
{exact snippet of current code, verbatim from file}
~~~

**After** (target code after the change):
~~~go
{exact snippet of proposed code, verbatim}
~~~

**Rationale**: {why this change, referenced to spec section or threat model ID.
If the change is security-motivated, name the control pattern from go-security.}

#### `[DELETE]` {relative/path/to/file.go}
{Reason for deletion.}

---

## Guardrails

The Auditor's checklist. Each entry is a binary check with a pass condition.

| ID | Check | Pass condition |
|----|-------|----------------|
| S1 | {specific check} | {exact pass criterion — grep-able or testable} |
| S2 | {check} | {criterion} |

Mandatory guardrails when applicable:
- Projects touching files: S# for `safePath` with separator suffix verified
- Projects doing crypto: S# for `crypto/rand` and `crypto/subtle` usage
- Projects with JSON envelopes: S# for `LimitReader` + `DisallowUnknownFields`
- All projects: S# for no `_` on error returns in security-critical paths
- All projects: S# for no TODO/FIXME left as implementation

---

## CI Mirror

Automated checks that enforce the spec's requirements. Exact commands.

```bash
go vet ./...
staticcheck ./...
gosec -quiet ./...
govulncheck ./...
go test -race -count=1 ./...
```

Each command maps to spec requirements it enforces. If a requirement has no corresponding CI check, flag it as a gap.

---

<!--
  Template notes for the Architect (remove before finalizing):
  - Every {placeholder} must be replaced or the section removed explicitly.
  - Never leave "TBD" or "TODO" in a final spec. Elicit or decide.
  - If a section does not apply to the project, replace its body with "Not applicable" plus one sentence of rationale. Do not delete the section header — structural uniformity helps the Auditor.
  - The Blueprint, Guardrails, and CI Mirror sections are mandatory. The spec fails without them.
-->
