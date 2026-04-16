# Flow States

Defines every valid state in the Kerux development cycle, the transitions between them, and failure handling at each stage.

## State Codes

IDL=IDLE MAP=MAPPING DES=DESIGNING SCF=SCAFFOLDING IMP=IMPLEMENTING REV=REVIEWING STG=STAGING COM=COMMITTED FAI=FAILED

## States

### IDLE (IDL)
- **Owner**: Kerux
- **Entry**: Session boot complete, or previous task reached COMMITTED/FAILED.
- **Exit condition**: User provides project-scoped task → MAPPING.
- **Out-of-scope requests**: Kerux answers directly, stays in IDLE.

### MAPPING (MAP)
- **Owner**: Analyst
- **Entry**: Kerux dispatch with user intent.
- **Exit artifact**: Context packet (structure, module, deps, scaffold, spec, security).
- **Transitions**:
  - Success → DESIGNING
  - Empty/greenfield project → DESIGNING (packet is minimal; Architect enters elicitation)
  - Fatal (no shell access, corrupt workspace) → FAILED

### DESIGNING (DES)
- **Owner**: Architect
- **Entry**: Analyst's context packet received.
- **Exit artifact**: `spec_projeto.md` conforming to SPEC_TEMPLATE.md.
- **Sub-phases**:
  - Phase 1: Assess context sufficiency.
  - Phase 2: Elicitation (if needed, max 3 rounds with user).
  - Phase 3: Spec authoring.
- **Transitions**:
  - Success + new project → SCAFFOLDING
  - Success + existing project → IMPLEMENTING
  - Template missing → FAILED
  - User aborts elicitation → IDLE

### SCAFFOLDING (SCF)
- **Owner**: Engineer
- **Entry**: Spec has `new_project: true` flag or empty project structure.
- **Exit artifact**: Project directory, go.mod, initial file tree. Spec copied to project root.
- **Action**: Engineer generates `lazygo.yml` from spec, runs `lazy.go init --from <path>`.
- **Transitions**:
  - Success → IMPLEMENTING
  - Lazy.go error → DESIGNING (Architect adjusts spec)
  - Path conflict → FAILED (unrecoverable filesystem state)

### IMPLEMENTING (IMP)
- **Owner**: Engineer
- **Entry**: Spec is current. Project structure exists.
- **Exit artifact**: Files modified/created per Blueprint. Build passes (`go build ./...`).
- **Transitions**:
  - Success → REVIEWING
  - BLOCKED (spec incomplete/ambiguous) → DESIGNING with BLOCKED packet
  - Build error → self-retry (max 2 attempts), then DESIGNING
  - Fatal → FAILED

### REVIEWING (REV)
- **Owner**: Auditor
- **Entry**: Engineer's handoff packet with change summary.
- **Exit artifact**: Verdict (PASS, COMMENT, REJECT) with evidence.
- **Transitions**:
  - PASS → STAGING
  - COMMENT → STAGING (notes attached)
  - REJECT (code-level) → IMPLEMENTING (Engineer fixes)
  - REJECT (design-level) → DESIGNING (Architect revises spec)
  - REJECT (security finding) → DESIGNING mandatory (no code-only fix for security issues)

### STAGING (STG)
- **Owner**: Kerux
- **Entry**: Auditor PASS or COMMENT.
- **Exit artifact**: Commit message (Conventional Commits format).
- **Action**: Present diff + proposed commit to user.
- **Transitions**:
  - User approves → COMMITTED
  - User requests changes → IMPLEMENTING
  - User abandons → IDLE (changes remain on disk, uncommitted)

### COMMITTED (COM)
- **Owner**: Kerux
- **Entry**: User explicit approval in current turn (Commandment C1).
- **Action**: `git add` + `git commit`. Never `git push`.
- **Transitions**: → IDLE

### FAILED (FAI)
- **Owner**: Kerux
- **Entry**: Any unrecoverable error.
- **Action**: Report full error context to user. Suggest manual intervention.
- **Transitions**: → IDLE (after user acknowledgment)

## Transition Table

| From | To | Trigger | Role |
|------|----|----|------|
| IDL | MAP | User task in scope | Analyst |
| MAP | DES | Context packet emitted | Architect |
| DES | SCF | Spec complete, new project | Engineer |
| DES | IMP | Spec complete, existing project | Engineer |
| DES | FAI | Template missing | Kerux |
| SCF | IMP | Scaffold complete | Engineer |
| SCF | DES | Lazy.go error | Architect |
| IMP | REV | Implementation complete, build passes | Auditor |
| IMP | DES | BLOCKED packet | Architect |
| IMP | FAI | Fatal error after retries | Kerux |
| REV | STG | PASS or COMMENT | Kerux |
| DES | DES | Gate: spec incomplete | Architect |
| REV | REV | Gate: evidence incomplete | Auditor |
| REV | IMP | REJECT code-level | Engineer |
| REV | DES | REJECT design-level or security | Architect |
| STG | COM | User approval | Kerux |
| STG | IMP | User requests changes | Engineer |
| STG | IDL | User abandons | Kerux |
| COM | IDL | Commit complete | Kerux |
| FAI | IDL | User acknowledges | Kerux |
| * | IDL | User abort | Kerux |

## Invariants

1. Only one role holds the lock at a time. No concurrent execution.
2. State transitions are logged in session.json (if PERSISTENCE_MODE=file).
3. Every REJECT specifies the target state — Auditor decides Engineer fix vs Architect revision.
4. STG → COM requires user approval in the CURRENT turn (Commandment C1).
5. IMP → DES via BLOCKED is correct behaviour, not a failure.
6. User can abort at any state; flow returns to IDLE with artifacts preserved.
7. DES→IMP and DES→SCF require the Lead to verify spec completeness per KERUX_CALIBRATION_SPEC §2.1. The Architect does not self-certify.
8. REV→STG PASS requires an evidence block per KERUX_CALIBRATION_SPEC §2.3. The Auditor does not self-certify.
