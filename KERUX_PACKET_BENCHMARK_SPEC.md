# KERUX_PACKET_BENCHMARK_SPEC v1.0.0

**Project**: Kerux Orchestrator — Packet Format Evaluation  
**Author**: hadnu  
**Date**: 2026-04-15  
**Status**: DRAFT — Pending owner approval  
**Precondition**: `KERUX_CONSOLIDATION_SPEC v1.0.0` defines the architectural context for both formats.

---

## 0. Glossary

| Term | Definition |
|------|-----------|
| **Arm** | One complete run of the Kerux flow from IDLE to COMMITTED using a specific packet format. |
| **Marker format (ARM-M)** | State transition markers. Positional shorthand, ~20 tokens per packet. Assumes shared context window. |
| **JSON format (ARM-J)** | Compact self-contained JSON envelopes, ~55 tokens per packet. Carries data inline. |
| **Control spec** | A fixed `spec_projeto.md` used identically in both arms. Not generated per-arm. |
| **Token budget** | Total input + output tokens consumed from IDLE to COMMITTED across all personas. |
| **Rejection cost** | Additional tokens consumed by state transitions caused by Reviewer REJECT verdicts. |
| **Security baseline** | The minimum set of checks that the delivered code must pass to be considered secure. |

---

## 1. Hypothesis

**H₀ (null):** Both packet formats produce equivalent total token budgets when delivering a secure, functional Go application through the full Kerux flow.

**H₁ (alternative):** The marker format (ARM-M) consumes fewer total tokens than the JSON format (ARM-J) for the same deliverable quality, because it avoids data duplication within a shared context window.

**H₂ (risk hypothesis):** The marker format produces more Reviewer rejections than the JSON format, because compressed markers may lose contextual signal that the Reviewer needs for accurate auditing.

The benchmark validates or falsifies all three hypotheses with measured data.

---

## 2. Experimental Design

### 2.1 Structure

Two-arm controlled experiment. One independent variable (packet format), one fixed control (the spec), measured outcomes.

```
                    ┌─── ARM-M (marker format) ───── measure ──┐
CONTROL SPEC ──────┤                                            ├─── compare
                    └─── ARM-J (JSON format)  ───── measure ──┘
```

### 2.2 Control Variables (held constant across both arms)

| Variable | Value | Rationale |
|----------|-------|-----------|
| `spec_projeto.md` | Pre-authored, identical in both arms | Eliminates Architect variability. |
| Persona instructions | Identical `.kerux/` files (post-consolidation) | No persona behavioural difference between arms. |
| LLM model | Same model, same temperature, same system prompt | Controls for model variability. |
| Security baseline | Defined in §5 — identical checks for both arms | Same quality bar. |
| Go toolchain | Same `go` version, same linters, same flags | Deterministic tooling output. |
| Project scope | Defined in §3 — identical requirements | Same implementation target. |

### 2.3 Independent Variable

| Arm | Packet format | Schema |
|-----|--------------|--------|
| ARM-M | State transition marker | `→{target}\|{state}\|{delta}\|{focus}` |
| ARM-J | Compact JSON envelope | `{"id","from","to","st","v":{},"s":""}` |

### 2.4 Dependent Variables (measured)

| Metric | Unit | Captured at |
|--------|------|-------------|
| **M1: Total token budget** | tokens (input + output) | End of arm |
| **M2: Tokens per state transition** | tokens / transition | Each state boundary |
| **M3: Number of state transitions** | count | End of arm |
| **M4: Reviewer rejection count** | count | Each REVIEWING state exit |
| **M5: Rejection cost** | tokens | Sum of tokens consumed in reject→rework→re-review cycles |
| **M6: Security baseline pass rate** | % (0–100) | Post-COMMITTED |
| **M7: Functional completeness** | % (0–100) | Post-COMMITTED |
| **M8: Code quality score** | composite (vet + staticcheck + gosec findings) | Post-COMMITTED |

### 2.5 Derived Metrics

| Metric | Formula | Purpose |
|--------|---------|---------|
| **D1: Token efficiency** | M1(ARM-J) − M1(ARM-M) | Raw token savings. |
| **D2: Efficiency ratio** | M1(ARM-M) / M1(ARM-J) | Percentage of baseline cost. |
| **D3: Rejection overhead** | M5 / M1 per arm | Proportion of budget spent on rework. |
| **D4: Cost-adjusted quality** | M1 / (M6 × M7 / 100) per arm | Tokens per unit of quality. Lower is better. |
| **D5: Net efficiency** | D1 − [M5(ARM-M) − M5(ARM-J)] | Token savings minus any rejection penalty from marker format. |

D5 is the decisive metric. If D5 > 0, the marker format wins net. If D5 ≤ 0, the compression savings are eaten by rework cost and JSON wins.

---

## 3. Test Project: `sigcheck`

A minimal Go CLI tool chosen to exercise all Kerux personas and the security baseline without exceeding ~300 LOC.

### 3.1 Functional Requirements

`sigcheck` verifies the integrity of files against pre-computed SHA-256 checksums stored in a manifest file.

```
sigcheck verify --manifest checksums.sha256 --target ./dist/
sigcheck generate --target ./dist/ --output checksums.sha256
```

**`generate`**: Walk a target directory, compute SHA-256 for each file, write a manifest.  
**`verify`**: Read the manifest, recompute hashes for each listed file, report mismatches and missing files. Exit code 0 if all match, exit code 1 if any mismatch or missing, exit code 2 on operational error.

### 3.2 Why This Project

| Kerux persona | Exercise point |
|--------------|---------------|
| **Tracker** | Must map the directory structure, identify `go.mod`, determine project type. Non-trivial because Cobra subcommands create a multi-file `cmd/` layout. |
| **Architect** | Must design the spec with file walking, hash computation, manifest parsing. Decisions around error handling (partial failures vs fail-fast). |
| **Coder** | Must implement file walking with goroutines (security: path traversal), crypto (must use `crypto/sha256`, not `md5`), CLI (Cobra with RunE). |
| **Reviewer** | Must audit: path traversal protection, no race conditions in concurrent hashing, correct exit codes, no swallowed errors, no hardcoded paths. |

### 3.3 Security Surface

The project has a real security surface, not an artificial one:

| Risk | Expected control | Reviewer check |
|------|-----------------|---------------|
| Path traversal via manifest entries | `safePath()` validation — candidate must stay under base dir | Grep for `filepath.Join` without prefix check |
| TOCTOU in verify (file changes between hash and report) | Document as known limitation in spec; no mitigation required for v1 | Reviewer confirms spec acknowledges the limitation |
| Symlink following | `filepath.WalkDir` with `d.Type()` check — skip symlinks | Grep for symlink handling in walk function |
| Large file DoS | `io.LimitReader` or streaming hash (no full-file `os.ReadFile`) | Grep for `os.ReadFile` in hash function — must not exist |
| Manifest injection (malicious paths in manifest) | Parse manifest with path validation per line | Grep for validation in manifest parser |

### 3.4 Pre-authored Control Spec

The `spec_projeto.md` for `sigcheck` is authored once, manually, before either arm runs. It follows the `SPEC_TEMPLATE.md` structure and includes:

- Overview matching §3.1 above.
- Architecture: `cmd/root.go`, `cmd/generate.go`, `cmd/verify.go`, `internal/hasher/hasher.go`, `internal/manifest/manifest.go`.
- Blueprint: five [NEW] files, one [NEW] `go.mod`.
- Pseudocode for the `verify` flow (walk → hash → compare → report).
- Guardrails matching §3.3 above.
- CI Mirror: `go vet`, `staticcheck`, `gosec`, `go test -race`.

This spec is committed to the repository before either arm runs. Neither arm's Architect persona modifies it — the Architect's role in both arms is limited to acknowledging the spec and handing off to the Coder.

---

## 4. Packet Format Specifications

### 4.1 ARM-M: State Transition Marker

```
→{target}|{state}|{delta}|{focus}
```

| Slot | Content | Max length |
|------|---------|-----------|
| `target` | Single-letter persona code: T(racker), A(rchitect), C(oder), R(eviewer), H(erald) | 1 char |
| `state` | Transition in `FROM→TO` notation, using 3-letter state abbreviations | ~7 chars |
| `delta` | What changed. Comma-separated facts. No articles, no filler. | ≤ 80 chars |
| `focus` | What the target persona should prioritize. Imperative fragments. | ≤ 60 chars |

**State abbreviations**: IDL (IDLE), MAP (MAPPING), DES (DESIGNING), SCF (SCAFFOLDING), IMP (IMPLEMENTING), REV (REVIEWING), STG (STAGING), COM (COMMITTED), FAI (FAILED).

**Example — Tracker to Architect:**
```
→A|MAP→DES|5 files mapped, cobra cmd/, go1.22, no lazygo.yml|spec from template, scanner type, high crit
```

**Example — Reviewer PASS:**
```
→H|REV→STG|all checks pass, 0 findings, safePath verified|stage for commit, conventional msg
```

**Example — Reviewer REJECT to Coder:**
```
→C|REV→IMP|REJECT: hasher.go L47 os.ReadFile on untrusted size, need streaming hash|fix hasher.go, keep verify flow, re-submit
```

**Constraints:**
- No data duplication. Never repeat file contents, command outputs, or code blocks that are already in the context window.
- If referencing a prior result, use a position anchor: "see find output above" or "per go.mod L3".
- Total packet size: ≤ 150 chars target, hard ceiling 200 chars.

### 4.2 ARM-J: Compact JSON Envelope

```json
{
  "id": "KRX-NNN",
  "from": "T",
  "to": "A",
  "st": "MAP→DES",
  "v": {
    "files": 5,
    "go": "1.22",
    "cmd": "cobra",
    "scaffold": false
  },
  "s": "5 files mapped. No lazygo.yml. Focus: spec from template.",
  "refs": ["cmd/root.go", "go.mod"]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | Sequential ID: `KRX-001`, `KRX-002`, ... |
| `from` | string | yes | Single-letter origin persona code |
| `to` | string | yes | Single-letter target persona code |
| `st` | string | yes | State transition `FROM→TO` |
| `v` | object | no | Structured variables relevant to the transition. Key names abbreviated. |
| `s` | string | yes | Summary. ≤ 120 chars. Compressed natural language (caveman-tier). |
| `refs` | array | no | File paths referenced. Does not include file content — content is in the context window. |

**Example — Reviewer REJECT:**
```json
{
  "id": "KRX-005",
  "from": "R",
  "to": "C",
  "st": "REV→IMP",
  "v": {"verdict": "REJECT", "file": "hasher.go", "line": 47},
  "s": "os.ReadFile on untrusted size. Need streaming hash.",
  "refs": ["internal/hasher/hasher.go"]
}
```

**Constraints:**
- No pretty-printing. Single line, no indentation.
- `v` keys abbreviated: 3-6 chars max.
- `s` field follows caveman compression rules: no articles, no filler, fragments permitted.
- `refs` does NOT carry content — only paths. Content is in the window.

---

## 5. Security Baseline

Both arms must deliver code that passes every check below. A failure on any check means the arm did not produce a "secure and functional" deliverable, regardless of token efficiency.

### 5.1 Automated Checks

```bash
# All must exit 0
go vet ./...
staticcheck ./...
gosec -quiet ./...
go test -race -count=1 ./...
```

### 5.2 Manual Security Audit (Reviewer persona)

| ID | Check | Pass condition |
|----|-------|----------------|
| S1 | Path traversal protection | `safePath()` or equivalent used on all paths derived from manifest input. |
| S2 | Streaming hash | No `os.ReadFile` or `io.ReadAll` on files of unbounded size. Must use `io.Copy` into `sha256.New()`. |
| S3 | Symlink handling | `WalkDir` callback checks `d.Type()` and skips symlinks, or documents why it follows them. |
| S4 | Error propagation | No `_` on `error` returns in security-critical paths (hash computation, file open, manifest parse). |
| S5 | Exit codes | `0` = all match, `1` = mismatch/missing, `2` = operational error. Tested. |
| S6 | No hardcoded paths | All paths from CLI flags or arguments. No default paths pointing to system directories. |
| S7 | Manifest parsing | Each line validated before use. Malformed lines rejected with error, not silently skipped. |
| S8 | Concurrency safety | If goroutines used for hashing, `sync.WaitGroup` with proper join. No shared mutable state without mutex. |

### 5.3 Functional Completeness

| ID | Requirement | Verification |
|----|-------------|-------------|
| F1 | `generate` produces a valid manifest from a target directory | Test with known directory, compare output to expected manifest. |
| F2 | `verify` detects a tampered file | Modify one file after generate, verify reports mismatch with exit 1. |
| F3 | `verify` detects a missing file | Remove one file after generate, verify reports missing with exit 1. |
| F4 | `verify` succeeds on unmodified directory | Generate then verify without changes, exit 0. |
| F5 | Error on non-existent target | Both commands return exit 2 with error message. |

---

## 6. Measurement Protocol

### 6.1 Token Counting

Tokens are counted at each state transition boundary. The count includes:

- **Input tokens**: All context visible to the model at the start of the turn (system prompt, `.kerux/` files, conversation history, tool outputs).
- **Output tokens**: All tokens generated by the model in that turn (reasoning, tool calls, packet emission, code).

Counting method depends on runtime:
- **API-based runtime**: Read `usage.input_tokens` and `usage.output_tokens` from the API response.
- **Chat-based runtime**: Estimate using `tiktoken` or the model's tokenizer on the raw conversation export. Flag as estimated in results.

### 6.2 Transition Log

Each arm maintains a transition log — a CSV with one row per state transition:

```csv
arm,transition_id,from_state,to_state,persona,input_tokens,output_tokens,verdict,notes
ARM-M,1,IDL,MAP,Tracker,1200,340,,initial mapping
ARM-M,2,MAP,DES,Architect,1540,180,,spec acknowledged
ARM-M,3,DES,IMP,Coder,1720,2800,,full implementation
ARM-M,4,IMP,REV,Reviewer,4520,600,PASS,all checks passed
ARM-M,5,REV,STG,Herald,5120,120,,staged for commit
ARM-M,6,STG,COM,Herald,5240,80,,user approved
```

### 6.3 Execution Protocol

1. **Prepare**: Author the control `spec_projeto.md` for `sigcheck`. Commit to repo.
2. **ARM-J first**: Run the full Kerux flow with JSON packets. Log every transition. Save the delivered code as `sigcheck-arm-j/`.
3. **Clean state**: Reset the workspace. The LLM session starts fresh — no carryover from ARM-J.
4. **ARM-M second**: Run the full Kerux flow with marker packets. Log every transition. Save the delivered code as `sigcheck-arm-m/`.
5. **Validate**: Run §5 checks on both deliverables independently.
6. **Measure**: Compile transition logs. Compute all metrics from §2.4 and §2.5.

ARM-J runs first to avoid contamination: if the model "remembers" an efficient solution from ARM-M, ARM-J might benefit unfairly. Running the more verbose format first means any learning effect biases *against* the hypothesis (ARM-M looks worse, not better), making a positive result more robust.

### 6.4 Confound Mitigation

| Confound | Mitigation |
|----------|-----------|
| LLM non-determinism | Run each arm twice (2 × ARM-M, 2 × ARM-J). Report mean and range. |
| Model learning across arms | Fresh session per run. No conversation history carryover. |
| Spec interpretation variance | Pre-authored control spec with unambiguous pseudocode. |
| Token counting imprecision | Use API token counts where available. Flag estimates. |
| Reviewer subjectivity | Security baseline is a binary checklist, not a judgement call. |

---

## 7. Expected Outcomes

### 7.1 If H₁ is confirmed (marker format cheaper, net)

Adopt marker format as default for `PERSISTENCE_MODE=file` and `PERSISTENCE_MODE=none` (single-session). Keep JSON format available as fallback for `PERSISTENCE_MODE=memory` (cross-session, context may be lost).

Update `KERUX_CONSOLIDATION_SPEC` §5.1 (`packet-schema.md`) to define the marker as primary and JSON as fallback.

### 7.2 If H₁ is falsified (JSON cheaper or equivalent, net)

Adopt JSON format as the universal packet format. The marker format is discarded — the complexity of supporting two formats is not justified without a measurable benefit.

Update `KERUX_CONSOLIDATION_SPEC` §5.1 to define JSON as the sole format.

### 7.3 If H₂ is confirmed (markers cause more rejections)

Even if H₁ is confirmed (markers cheaper gross), the rejection overhead must be quantified. If D3(ARM-M) is significantly higher than D3(ARM-J), markers save on transport but cost on rework. The D5 (net efficiency) metric resolves this: if D5 > 0, markers still win despite rejections. If D5 ≤ 0, the rejection cost negates the compression savings.

### 7.4 Decision Table

| D5 | M6/M7 parity | Decision |
|----|--------------|----------|
| D5 > 0, both arms pass security/functional baseline equally | **Marker format wins.** Adopt as default. |
| D5 > 0, ARM-M has lower M6 or M7 | **Inconclusive.** Markers cheaper but delivered lower quality. Investigate root cause before adopting. |
| D5 ≤ 0, both arms pass equally | **JSON format wins.** Compression savings don't offset rework. |
| D5 ≤ 0, ARM-J has lower M6 or M7 | **Edge case.** JSON cheaper but lower quality. Re-examine the control spec for ambiguity. |

---

## 8. Deliverables

| ID | Artifact | Format | Produced by |
|----|----------|--------|-------------|
| D-1 | Control `spec_projeto.md` for `sigcheck` | Markdown | Author (manual, pre-experiment) |
| D-2 | ARM-J transition log | CSV | Kerux (during ARM-J runs) |
| D-3 | ARM-M transition log | CSV | Kerux (during ARM-M runs) |
| D-4 | ARM-J delivered code | Go source (`sigcheck-arm-j/`) | Coder persona |
| D-5 | ARM-M delivered code | Go source (`sigcheck-arm-m/`) | Coder persona |
| D-6 | Security baseline results | Checklist (pass/fail per check, per arm) | Reviewer persona + automated tools |
| D-7 | Benchmark report | Markdown | Author (post-experiment analysis) |

### 8.1 Benchmark Report Structure

```
# Kerux Packet Benchmark — Results

## Summary
- D5 (net efficiency): [value]
- Decision: [marker wins / JSON wins / inconclusive]

## Raw Metrics
[Table: M1–M8 per arm, mean of 2 runs, with range]

## Derived Metrics
[Table: D1–D5]

## Transition Log Analysis
[Per-state token breakdown. Which states consumed the most? Where did rejections occur?]

## Security Parity
[Side-by-side S1–S8 results]

## Observations
[Qualitative notes on where each format struggled or excelled]

## Recommendation
[Adopt format X. Rationale grounded in D5 and quality parity.]
```

---

## 9. CI Mirror

No CI pipeline for this spec — it's an experiment, not a software project. Validation is manual and follows the execution protocol in §6.3.

**Pre-experiment checklist:**
```
- [ ] Control spec_projeto.md authored and committed
- [ ] .kerux/ files reflect post-consolidation state (KERUX_CONSOLIDATION_SPEC applied)
- [ ] Packet format definitions (§4.1, §4.2) transcribed into the relevant .kerux/ skill files
- [ ] Token counting method confirmed (API or estimate)
- [ ] Workspace reset procedure documented
```

**Post-experiment checklist:**
```
- [ ] 2 × ARM-J runs completed with transition logs
- [ ] 2 × ARM-M runs completed with transition logs
- [ ] Security baseline (§5) run on all 4 deliverables
- [ ] Functional completeness (§5.3) verified on all 4 deliverables
- [ ] Metrics M1–M8 computed per run
- [ ] Derived metrics D1–D5 computed
- [ ] Benchmark report (§8.1) authored
- [ ] Decision recorded and KERUX_CONSOLIDATION_SPEC updated accordingly
```

---

## 10. Scope Boundaries

**IN scope:** Measuring token efficiency of inter-persona packet formats within a single Kerux session on a controlled Go project.

**OUT of scope:**
- Evaluating the Kerux orchestrator itself (personas, flow, rules) — that's the consolidation spec's domain.
- Comparing LLM models — same model in both arms.
- Measuring boot-time token cost of `.kerux/` files — that's a separate experiment (caveman-compress for instruction files).
- Cross-session persistence formats — both arms run in single sessions.
- Cost in dollars — token counts are model-agnostic; pricing depends on the provider.

---

*KERUX_PACKET_BENCHMARK_SPEC v1.0.0 — End of document.*
