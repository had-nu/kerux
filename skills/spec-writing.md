# Skill: Spec Writing

Architect's guide to producing doc/spec_projeto.md.
Template: `.kerux/templates/SPEC_TEMPLATE.md` — structure is mandatory.

## Principles

A spec is executable by the Engineer and verifiable by the Auditor. If either role has to guess, the spec failed.

- Precision over verbosity. One unambiguous sentence > three that describe intent vaguely.
- Every requirement is testable. If the Auditor cannot verify it, rewrite it.
- Pseudocode in the Blueprint is a contract, not inspiration. The Engineer follows it literally.
- Every trust boundary in the system has a threat model entry with a named control.

## Section Requirements

### Overview
- Problem statement: what exists, what's broken, why it matters.
- Proposed solution: one paragraph, concrete.
- Scope: IN (list), OUT (list). No "maybe" — explicit.

### Goals
- Measurable outcomes. "Pass all gosec checks" not "be secure."
- Non-goals stated. Prevents scope creep.

### Architecture
- System diagram: text-based, show data flow.
- Component inventory table: component, responsibility, tech, notes.
- Data flow: step-by-step for each major operation.
- External dependencies: name, version, purpose, risk if unavailable.

### Data Model
- Core structs in Go. Field types, JSON tags if relevant.
- Storage: plain text, SQLite, in-memory — say which.
- Lifecycle: when created, when purged.

### Interfaces
- CLI surface: command + flag table (type, required, default, description).
- Configuration: files, env vars — or explicitly "none."
- Output formats: stdout, exit codes, file formats.

### Performance
- Targets with units and boundary conditions.
- Bottlenecks acknowledged. Non-goals explicit.

### Security
This section is the Auditor's primary reference.
- Threat model table: ID, threat, attack vector, control.
- Every file I/O: path traversal control named.
- Every user input: validation rule named.
- Every crypto operation: algorithm + library named.
- Known limitations: documented explicitly, no hidden gaps.

### Testing
- Unit tests table: package, test name, description.
- E2E tests table: ID, test, steps, expected.
- Coverage targets if any.

### Blueprint
The executable section.

File operations marked explicitly:
- `[NEW]` file path — creates new file
- `[MODIFY]` file path — modifies existing
- `[DELETE]` file path — removes existing

For each [NEW] or [MODIFY]:
- What the file does (one sentence).
- Pseudocode for complex logic. Precise enough that the Engineer doesn't guess.
- Reference to security patterns: "Use safePath() from go-security skill."

### Guardrails
The Auditor's checklist. Each entry:
- ID (S1, S2, ...)
- Check (specific, grep-able)
- Pass condition (binary)

If the project handles files from external input, S# must include:
- Path traversal: `safePath()` with separator suffix verified
- Streaming I/O: no unbounded reads
- Error propagation: no `_` on critical paths

### CI Mirror
Exact commands that enforce spec requirements.
```bash
go vet ./...
staticcheck ./...
gosec -quiet ./...
go test -race -count=1 ./...
```

### Decision Log
Every non-obvious choice:
- ID, decision, rationale, date, status.
- "Assumed — pending user confirmation" is a valid status for elicitation-phase assumptions.

## Rules

- Never use "should" when "must" is meant. Specs are contracts.
- Never leave a trust boundary without a named control.
- Never write pseudocode that requires the Engineer to invent a security pattern — reference the go-security skill explicitly.
- If the user's intent is ambiguous and elicitation didn't fully resolve it, document the assumption in Decision Log, do not silently guess.
- A spec that the Auditor cannot audit is a failed spec. Re-read from the Auditor's perspective before handoff.

## Anti-Patterns

- Bullet points masquerading as prose. If steps are causal, write prose.
- TBD or TODO in any section. Elicit or decide — never defer.
- Pseudocode in English paragraphs instead of structured logic. Use code-like blocks.
- Security section that says "validate input" without naming HOW.
- Blueprint without file paths. Every change is anchored to a file.
- Modifying an existing file without citing the exact line anchor and current
  snippet. Brownfield changes must show what exists before showing what
  replaces it. The Engineer verifies the Before block matches the live file;
  if it doesn't, the spec is stale and the Engineer emits BLOCKED via the
  escape hatch.
