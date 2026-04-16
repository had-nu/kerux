# Error Taxonomy

Failure categories and the expected response for each. Roles match situations to this taxonomy rather than improvising.

## Severity Levels

### FATAL
**Definition**: Task cannot continue; no automated recovery possible.

Examples:
- Required tool missing from PATH (`go`, `git`).
- Project directory deleted mid-task.
- `SPEC_TEMPLATE.md` missing when Architect needs it.
- Corrupt `spec_projeto.md` that fails to parse.
- Workspace permissions deny read/write.

Response:
- Transition to FAILED state.
- Report full error context to user.
- Suggest manual intervention.
- Never attempt workarounds.

### BLOCKING
**Definition**: Current role cannot proceed, but another role or the user can resolve.

Examples:
- Spec incomplete (Engineer → Architect).
- Dependency not declared in spec (Engineer → Architect).
- Ambiguous user intent (any role → Kerux → user).
- Missing test fixtures referenced in spec (Engineer → Architect).

Response:
- Emit BLOCKED packet to Kerux with:
  - Blocking reason (what is missing or ambiguous).
  - Suggested resolver (usually Architect, sometimes user).
- Kerux routes to appropriate role or elicits from user.

### DEGRADED
**Definition**: Task continues with reduced quality or missing optional features.

Examples:
- Web search unavailable (Analyst proceeds with local-only mapping).
- PERSISTENCE_MODE=none (memory compression produces paste-block, no file write).
- Linter unavailable but build passes (Auditor notes gap in COMMENT).
- Optional tool missing (govulncheck absent — audit proceeds with gosec only).

Response:
- Continue execution.
- Log degradation in handoff packet summary: `DEGRADED: {what}`.
- Do not escalate unless degradation compounds across steps.
- User sees degradation in status line emitted by Kerux.

### RECOVERABLE
**Definition**: Transient failure the current role can retry.

Examples:
- Shell command timeout.
- File lock contention (`.git/index.lock` stale).
- Network hiccup during dependency fetch.
- Transient filesystem error (EAGAIN, EINTR).

Response:
- Retry once with small backoff.
- If second attempt fails, escalate to BLOCKING (not FATAL — user can often fix transient issues).
- Log the retry in handoff packet.

## Conflict Resolution

### File Conflict
Two roles want to modify the same file (edge case; should not happen if flow-states are respected).

Resolution:
- Role with the lock (current state owner) has priority.
- Other role's changes are queued as a follow-up packet.
- Kerux logs the conflict and resumes normal flow.

### Verdict Disagreement
Auditor rejects, but the rejection contradicts the spec as written.

Resolution:
- Kerux presents both the spec text and the rejection rationale to the user.
- User decides.
- Not automatable — this is a trust boundary requiring human judgment.

### Ambiguous Intent
User request cannot be unambiguously mapped to a flow action.

Resolution:
- Kerux does not guess.
- Emits clarification request to user in IDLE state.
- Waits for response before dispatching.

## Escalation Path

RECOVERABLE → BLOCKING → FATAL

A single retry failure escalates to BLOCKING, not FATAL. The distinction matters:
- BLOCKING preserves session state and allows user intervention.
- FATAL is terminal; the only exit is IDLE after acknowledgment.

## Rules

1. Never silently swallow errors. Every error maps to a severity.
2. Never downgrade severity to avoid asking the user. BLOCKING exists so the user can unblock.
3. FATAL is rare and specific. Most errors are BLOCKING or RECOVERABLE.
4. DEGRADED accumulates. Two DEGRADED events in the same flow escalate to BLOCKING.
5. Roles never decide FATAL on their own authority. FATAL transitions are Kerux's decision.
