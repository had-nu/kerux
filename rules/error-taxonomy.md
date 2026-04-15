# Error Taxonomy v1

> Defines failure categories and the expected response for each.
> Roles do not improvise error handling — they match the situation to this taxonomy.

## Severity Levels

### FATAL
- **Definition**: The task cannot continue and no automated recovery is possible.
- **Examples**: Required tool missing (go, git not in PATH), project directory deleted mid-task,
  corrupt `spec_projeto.md` that fails to parse.
- **Response**: Transition to FAILED state. Report full error context to user.
  Do not attempt workarounds.

### BLOCKING
- **Definition**: The current role cannot proceed, but another role or the user can resolve it.
- **Examples**: Spec incomplete (Engineer → Architect), dependency not declared (Engineer → Architect),
  ambiguous user intent (any role → Kerux → user).
- **Response**: Emit a BLOCKED packet to Kerux with the blocking reason and suggested resolver.
  Kerux routes to the appropriate role or asks the user.

### DEGRADED
- **Definition**: The task can continue but with reduced quality or missing optional features.
- **Examples**: Web search unavailable (Analyst proceeds with local-only mapping),
  PERSISTENCE_MODE=none (memory compression skipped, session state in-context only).
- **Response**: Continue execution. Log the degradation in the handoff packet summary.
  Do not escalate unless the degradation compounds.

### RECOVERABLE
- **Definition**: A transient failure that the current role can retry.
- **Examples**: Shell command timeout, file lock contention, git index.lock stale.
- **Response**: Retry once. If second attempt fails, escalate to BLOCKING.

## Conflict Resolution

### File Conflict
Two roles want to modify the same file (e.g., Architect updates spec while Engineer is implementing).
- **Resolution**: The role with the lock (current flow state owner) has priority.
  The other role's changes are queued as a follow-up packet.

### Verdict Disagreement
Auditor rejects, but the rejection rationale contradicts the spec.
- **Resolution**: Kerux presents both the spec text and the rejection rationale to the user.
  User decides. This is not automatable.
