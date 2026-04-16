# Packet Schema

Contract for all inter-role communication. Every packet sent or received MUST conform.
Primary format: state transition marker (validated by packet benchmark, v1.0.0).
Fallback format: compact JSON (for cross-session persistence scenarios).

## Primary Format — State Transition Marker

```
→{target}|{state}|{delta}|{focus}
```

| Slot | Content | Max |
|------|---------|-----|
| `target` | Single-letter role code | 1 char |
| `state` | Transition `FROM→TO` in 3-letter codes | 7 chars |
| `delta` | What changed. Comma-separated facts. No articles, no filler. | 80 chars |
| `focus` | What the target should prioritize. Imperative fragments. | 60 chars |

Hard ceiling: 200 chars per packet.

### Role Codes

K=Kerux(Lead) N=Analyst A=Architect E=Engineer U=Auditor

### State Codes

IDL=IDLE MAP=MAPPING DES=DESIGNING SCF=SCAFFOLDING IMP=IMPLEMENTING REV=REVIEWING STG=STAGING COM=COMMITTED FAI=FAILED

### Examples

Analyst → Architect:
```
→A|MAP→DES|12 files mapped, cobra cmd/, go1.22, no lazygo.yml|spec from template, CLI type, high crit
```

Engineer → Auditor:
```
→U|IMP→REV|5 files created per spec, go build clean, race passes|audit security baseline S1-S8
```

Auditor PASS:
```
→K|REV→STG|all checks pass, 0 findings, safePath verified|stage for commit, conventional msg
```

Auditor REJECT to Engineer (code-level):
```
→E|REV→IMP|REJECT: hasher.go L47 os.ReadFile on untrusted size|fix to streaming, re-submit
```

Auditor REJECT to Architect (design-level):
```
→A|REV→DES|REJECT: spec missing symlink control in T3|amend threat model, re-handoff
```

Engineer BLOCKED:
```
→A|IMP→DES|BLOCKED: spec §3.2 silent on error aggregation|codebase has logger, suggest slog in spec|propose fix
```

### Validation Rules

1. `target` must be a valid role code.
2. `state` must be a valid transition per flow-states.md transition table.
3. `delta` and `focus` must not duplicate data already in context (spec, previous packets).
4. Packets over 200 chars must be rewritten or escalated to JSON fallback.
5. Malformed packet = DEGRADED log entry. Kerux fixes format before dispatch.

## Fallback Format — Compact JSON

Use when:
- PERSISTENCE_MODE=memory (cross-session, shared context may not be visible).
- Data volume exceeds marker capacity legitimately (not filler).
- Structured fields needed for audit trail beyond session.

```json
{"id":"KRX-NNN","from":"U","to":"E","st":"REV→IMP","v":{"verdict":"REJECT","file":"hasher.go","line":47},"s":"os.ReadFile unbounded. Need streaming.","refs":["internal/hasher/hasher.go"]}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | Sequential: `KRX-001`, `KRX-002`, ... |
| `from` | string | yes | Single-letter origin role code |
| `to` | string | yes | Single-letter target role code |
| `st` | string | yes | State transition `FROM→TO` |
| `v` | object | no | Structured variables. Keys ≤ 6 chars abbreviated. |
| `s` | string | yes | Summary. ≤ 120 chars. Caveman-compressed. |
| `refs` | array | no | File paths referenced. Paths only, not contents. |

Constraints:
- Single line, no indentation.
- Keys abbreviated aggressively.
- `s` field follows compression rules: no articles, fragments OK.
- `refs` carries paths, never contents.

## Format Selection Rule

PERSISTENCE_MODE=file or PERSISTENCE_MODE=none (single session):
- Default: marker format.
- Data is in context; packet is routing + delta.

PERSISTENCE_MODE=memory (cross-session):
- Default: JSON format.
- Context may not carry forward; packet must be self-contained enough.

Kerux selects the format at dispatch time based on runtime mode. Roles produce the format Kerux dispatched in.

## Invariants

1. Every transition produces exactly one packet.
2. Packets never carry file contents or full command output — only paths or summaries.
3. Packet size is measured after compression; verbose packets are a smell.
4. When in doubt between marker and JSON, choose marker. Upgrade only if content exceeds capacity.
