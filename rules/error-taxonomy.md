# Error Taxonomy v1

> Defines failure categories and the expected response for each.
> Personas do not improvise error handling — they match the situation to this taxonomy.

## Severity Levels

### FATAL
- **Definition**: The task cannot continue and no automated recovery is possible.
- **Examples**: Required tool missing (go, git not in PATH), project directory deleted mid-task,
  corrupt `spec_projeto.md` that fails to parse.
- **Response**: Transition to FAILED state. Report full error context to user.
  Do not attempt workarounds.

### BLOCKING
- **Definition**: The current persona cannot proceed, but another persona or the user can resolve it.
- **Examples**: Spec incomplete (Coder → Architect), dependency not declared (Coder → Architect),
  ambiguous user intent (any persona → Herald → user).
- **Response**: Emit a BLOCKED packet to Herald with the blocking reason and suggested resolver.
  Herald routes to the appropriate persona or asks the user.

### DEGRADED
- **Definition**: The task can continue but with reduced quality or missing optional features.
- **Examples**: Web search unavailable (Tracker proceeds with local-only mapping),
  PERSISTENCE_MODE=none (memory compression skipped, session state in-context only).
- **Response**: Continue execution. Log the degradation in the handoff packet summary.
  Do not escalate unless the degradation compounds.

### RECOVERABLE
- **Definition**: A transient failure that the current persona can retry.
- **Examples**: Shell command timeout, file lock contention, git index.lock stale.
- **Response**: Retry once. If second attempt fails, escalate to BLOCKING.

## Conflict Resolution

### File Conflict
Two personas want to modify the same file (e.g., Architect updates spec while Coder is implementing).
- **Resolution**: The persona with the lock (current flow state owner) has priority.
  The other persona's changes are queued as a follow-up packet.

### Verdict Disagreement
Reviewer rejects, but the rejection rationale contradicts the spec.
- **Resolution**: Herald presents both the spec text and the rejection rationale to the user.
  User decides. This is not automatable.
