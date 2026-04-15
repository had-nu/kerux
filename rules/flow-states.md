# Flow States v1

> **Authority**: Defines every valid state in the Kerux development cycle,
> the transitions between them, and the failure handling at each stage.

## States

### IDLE
- **Owner**: Kerux
- **Entry**: Session boot complete, or previous task committed/abandoned.
- **Exit**: User provides a task. → MAPPING

### MAPPING
- **Owner**: Analyst
- **Entry**: Kerux dispatches a mapping packet.
- **Exit artifacts**: Context Packet (file paths, dependencies, project metadata).
- **Transitions**:
  - Success → DESIGNING
  - Failure (project not found, ambiguous scope) → IDLE + user clarification request

### DESIGNING
- **Owner**: Architect
- **Entry**: Context Packet received from Analyst.
- **Exit artifacts**: `spec_projeto.md` conforming to `templates/SPEC_TEMPLATE.md`.
- **Transitions**:
  - Success → SCAFFOLDING (if new project) or IMPLEMENTING (if existing project)
  - Failure (template unavailable, conflicting requirements) → IDLE + user escalation

### SCAFFOLDING
- **Owner**: Engineer (under Kerux supervision)
- **Entry**: New project flag set in spec. `lazygo.yml` generated.
- **Exit artifacts**: Project directory created, `spec_projeto.md` copied to root.
- **Transitions**:
  - Success → IMPLEMENTING
  - Failure (lazy.go error, path conflict) → DESIGNING + error context

### IMPLEMENTING
- **Owner**: Engineer
- **Entry**: `spec_projeto.md` exists and is current.
- **Exit artifacts**: Modified/created files as listed in blueprint.
- **Transitions**:
  - Success → REVIEWING
  - Blocked (spec incomplete, missing dependency) → DESIGNING + BLOCKED packet
  - Failure (build error, test failure) → self-retry (max 2), then DESIGNING

### REVIEWING
- **Owner**: Auditor
- **Entry**: Engineer handoff packet with change summary.
- **Exit artifacts**: Verdict (PASS | REJECT | COMMENT) with evidence.
- **Transitions**:
  - PASS → STAGING
  - COMMENT → STAGING (with notes attached)
  - REJECT → route based on rejection type:
    - Blueprint deviation → IMPLEMENTING (Engineer fixes)
    - Architectural flaw → DESIGNING (Architect revises)
    - Security finding → DESIGNING (mandatory — security issues require spec-level response)

### STAGING
- **Owner**: Kerux
- **Entry**: Auditor PASS or COMMENT.
- **Exit artifacts**: Commit message (Conventional Commits format).
- **Action**: Present diff summary + proposed commit message to user.
- **Transitions**:
  - User approves → COMMITTED
  - User requests changes → IMPLEMENTING
  - User abandons → IDLE

### COMMITTED
- **Owner**: Kerux
- **Entry**: User explicit approval in current turn.
- **Action**: Execute git add + git commit. Never git push (user does this manually).
- **Transitions**: → IDLE

### FAILED
- **Owner**: Kerux
- **Entry**: Any unrecoverable error (missing tools, corrupt state, repeated failures).
- **Action**: Report full error context to user. Suggest manual intervention.
- **Transitions**: → IDLE (after user acknowledgment)

## Invariants

1. Only one role holds the lock at any time. No concurrent role execution.
2. State transitions are logged in the session todo (skills/agent-todo.md).
3. Every REJECT must include the target state for the rollback — the Auditor decides
   whether the fix is a Engineer task or an Architect task.
4. STAGING → COMMITTED requires user approval in the CURRENT turn (Commandment C1).
5. The IMPLEMENTING → DESIGNING (BLOCKED) transition is the Engineer's escape hatch
   for spec incompleteness. It is not a failure — it is correct behaviour.
