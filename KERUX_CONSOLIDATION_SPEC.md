# KERUX_CONSOLIDATION_SPEC v1.0.0

**Project**: Kerux Orchestrator  
**Author**: hadnu  
**Date**: 2026-04-15  
**Status**: DRAFT — Pending owner approval  
**Scope**: Full structural consolidation, flow formalization, runtime abstraction, ecosystem alignment.  
**Precondition**: Audit document `KERUX_AUDIT.md` (2026-04-15) accepted as diagnostic baseline.

---

## 0. Glossary

| Term | Definition |
|------|-----------|
| **Herald** | Kerux core — the orchestrator identity that routes traffic between personas. |
| **Persona** | A scoped agent role (Architect, Coder, Reviewer, Tracker) with defined inputs, outputs, and constraints. |
| **Packet** | The structured handoff unit between personas. Currently undefined — this spec formalizes it. |
| **Flow State** | A named stage in the development cycle with an owner, entry condition, and exit artifact. |
| **Runtime** | The LLM environment executing Kerux (Gemini, Claude, local agent, API). |
| **SPEC** | The `spec_projeto.md` — the source of truth for any project-level task. |
| **SoT** | Source of Truth. |

---

## 1. Problem Statement

The Kerux audit identified five categories of structural debt:

1. **Redundancy**: 4 files duplicate content that exists elsewhere, inflating boot-time token cost by ~30%.
2. **Undefined contracts**: The `<packet>` handoff format is referenced in 4 files but never specified. Personas cannot validate what they receive.
3. **Linear pipeline assumption**: The flow handles the happy path only. Rejection, rollback, spec incompleteness, and concurrent access have no defined behaviour.
4. **Runtime coupling**: Memory management, token thresholds, and attention hints are hardcoded for Gemini. The system cannot run on Claude or stateless APIs without modification.
5. **Ecosystem drift**: The Reviewer lacks supply chain and envelope checks despite the owner's established work in those areas. The Tracker doesn't parse Go-specific metadata. The SPEC template is referenced via an absolute `file:///` URI.

The correction converts Kerux from a well-structured intent document into a functional orchestrator specification.

---

## 2. Goals

| ID | Goal | Measurable Outcome |
|----|------|--------------------|
| G1 | Zero content duplication across `.kerux/` | Every paragraph of instructional text exists in exactly one file. |
| G2 | Every persona handoff uses a validated packet schema | `packet-schema.md` exists and is referenced by all persona files. |
| G3 | Flow handles rejection, rollback, and spec incompleteness | `flow-states.md` defines transitions for all identified failure modes. |
| G4 | Runtime-agnostic operation | No file references a specific LLM provider by name. A `runtime-contract.md` defines the abstraction. |
| G5 | Ecosystem-aligned audit capabilities | Reviewer checks supply chain and JSON envelope compatibility. Tracker parses `go.mod` and `lazygo.yml`. |
| G6 | SPEC template is portable | Template lives inside `.kerux/templates/` with a relative path reference. |

---

## 3. Architecture: Before / After

### 3.1 Before — Current File Tree (20 files)

```
.kerux/
├── kerux.md
├── personas/
│   ├── architect.md
│   ├── coder.md
│   ├── reviewer.md
│   └── tracker.md
├── rules/
│   ├── commandments.md
│   ├── edicts.md
│   └── memory-rules.md
├── skills/
│   ├── README.md                    ← EMPTY
│   ├── agent-memory.md              ← DUPLICATES memory-rules.md
│   ├── agent-todo.md                ← EMPTY
│   ├── code-review/
│   │   └── protocol.md
│   ├── context-maintenance.md
│   ├── core-orchestration.md        ← DUPLICATES boot + dispatch + context
│   ├── kerux-boot.md
│   ├── kerux-dispatch.md
│   ├── maintenance-skills.md        ← DUPLICATES todo + context
│   ├── memory-compression.md
│   ├── web-browser/
│   │   └── protocol.md
│   └── web-search/
│       └── protocol.md
```

### 3.2 After — Target File Tree (21 files)

```
.kerux/
├── kerux.md                          [MODIFY] — revised flow reference, scope section
├── VERSION                           [NEW]    — orchestrator version tracking
│
├── personas/
│   ├── architect.md                  [MODIFY] — relative template path
│   ├── coder.md                      [MODIFY] — spec-incomplete escape hatch
│   ├── reviewer.md                   [MODIFY] — supply chain + envelope checks
│   └── tracker.md                    [MODIFY] — go.mod, lazygo.yml, SPEC header parsing
│
├── rules/
│   ├── commandments.md               [KEEP]   — no changes
│   ├── edicts.md                     [KEEP]   — no changes
│   ├── memory-rules.md               [MODIFY] — absorb persistence protocol from agent-memory.md
│   ├── flow-states.md                [NEW]    — state machine definition
│   ├── packet-schema.md              [NEW]    — handoff format specification
│   ├── error-taxonomy.md             [NEW]    — failure modes and responses
│   └── runtime-contract.md           [NEW]    — environment abstraction layer
│
├── skills/
│   ├── kerux-boot.md                 [MODIFY] — version check at boot, runtime detection
│   ├── kerux-dispatch.md             [MODIFY] — state-aware dispatch, packet validation
│   ├── context-maintenance.md        [MODIFY] — single canonical source, remove Gemini refs
│   ├── memory-compression.md         [MODIFY] — runtime-agnostic reset protocol
│   ├── agent-todo.md                 [MODIFY] — fill with task tracking content
│   ├── code-review/
│   │   └── protocol.md              [MODIFY] — supply chain section
│   ├── web-browser/
│   │   └── protocol.md              [KEEP]   — no changes
│   └── web-search/
│       └── protocol.md              [KEEP]   — no changes
│
├── templates/
│   └── SPEC_TEMPLATE.md             [NEW/MOVE] — relocated from external file:/// path
│
└── memory/
    ├── lessons.md                    [NEW]    — explicit creation (referenced but never existed in tree)
    └── session.json                  [NEW]    — explicit creation (skeleton)
```

### 3.3 Operation Summary

| Operation | Count | Files |
|-----------|-------|-------|
| **DELETE** | 4 | `skills/core-orchestration.md`, `skills/maintenance-skills.md`, `skills/agent-memory.md`, `skills/README.md` |
| **NEW** | 8 | `VERSION`, `flow-states.md`, `packet-schema.md`, `error-taxonomy.md`, `runtime-contract.md`, `templates/SPEC_TEMPLATE.md`, `memory/lessons.md`, `memory/session.json` |
| **MODIFY** | 11 | `kerux.md`, `architect.md`, `coder.md`, `reviewer.md`, `tracker.md`, `memory-rules.md`, `kerux-boot.md`, `kerux-dispatch.md`, `context-maintenance.md`, `memory-compression.md`, `agent-todo.md`, `code-review/protocol.md` |
| **KEEP** | 4 | `commandments.md`, `edicts.md`, `web-browser/protocol.md`, `web-search/protocol.md` |

---

## 4. Blueprint — Phase 1: Eliminate Redundancy (P0)

Phase 1 is mechanical. No design decisions. Pure deletion and content migration.

### 4.1 DELETE `skills/core-orchestration.md`

**Reason**: Contains three sections that are verbatim summaries of `kerux-boot.md`, `kerux-dispatch.md`, and `context-maintenance.md`. Zero unique content.

**Dependency check**: No other file imports or references `core-orchestration.md` by name. `kerux.md` references boot and dispatch by their individual file names.

**Action**: Delete. No migration needed.

### 4.2 DELETE `skills/maintenance-skills.md`

**Reason**: Section 1 (Agent Todo) is the content that should be in the empty `agent-todo.md`. Section 2 (Context Maintenance) duplicates `context-maintenance.md` with a Gemini-specific `[!!]` tag addition.

**Action**:
1. Copy §Agent Todo content into `skills/agent-todo.md` (see §4.5).
2. Note the `[!!]` tag reference for removal during Phase 3 (runtime abstraction).
3. Delete `maintenance-skills.md`.

### 4.3 DELETE `skills/agent-memory.md`

**Reason**: Adds a Persistence Protocol and an Efficiency Protocol on top of `rules/memory-rules.md`. The rules file is the canonical location for memory behaviour.

**Action**:
1. Migrate §Persistence Protocol into `rules/memory-rules.md` as a new section `M5: Persistence Protocol` (see §5.3).
2. Migrate §Efficiency Protocol reference to compression into `memory-rules.md` M3 (it already mentions compression — add the cross-reference).
3. Delete `agent-memory.md`.

### 4.4 DELETE `skills/README.md`

**Reason**: Empty file. Zero signal.

**Action**: Delete.

### 4.5 FILL `skills/agent-todo.md`

**Source**: Content from `maintenance-skills.md` §Agent Todo.

**Target content**:

```markdown
# Skill: Agent Todo

> **Objective**: Persistent task tracking across sessions.

## Format
- [ ] Task Description (Owner: PersonaName)
- [x] Completed Task (Result: Summary)
- [/] In-Progress Task (Blocker: reason, if any)

## Usage
- **Kerux**: Initialize at boot. Load from `memory/session.json` if available.
- **Coder/Architect**: Update when state changes.
- **Handoff**: Include the delta of the todo list in the packet.

## Storage
Tasks persist in `memory/session.json` under the `tasks` key.
When PERSISTENCE_MODE=none (see runtime-contract.md), tasks exist only in-session.
```

---

## 5. Blueprint — Phase 2: Define Missing Contracts (P1)

### 5.1 NEW `rules/packet-schema.md`

This is the most-referenced undefined contract. Every persona file mentions packets. This file defines the shape.

```markdown
# Packet Schema v1

> **Authority**: This schema is the contract for all inter-persona communication.
> Every packet sent or received by any persona MUST conform to this structure.

## Schema

<packet>
  <id>Unique task identifier (format: KRX-YYYYMMDD-NNN)</id>
  <origin>Sending persona name</origin>
  <target>Receiving persona name</target>
  <state>Current flow state (reference: rules/flow-states.md)</state>
  <intent>Imperative verb phrase: what the receiver must do</intent>
  <context>
    <files>Ordered list of file paths relevant to this task</files>
    <dependencies>External requirements (tools, env vars, APIs)</dependencies>
    <constraints>Guardrails specific to this task — overrides nothing in Commandments</constraints>
  </context>
  <vars>
    Key=value pairs. Project name, target path, branch, flags.
  </vars>
  <summary>1-2 sentence description of what happened before this handoff</summary>
</packet>

## Validation Rules

1. `id` must be unique within a session. Kerux assigns IDs; personas do not.
2. `origin` and `target` must be valid persona names: Herald, Architect, Coder, Reviewer, Tracker.
3. `state` must be a valid state from `flow-states.md`.
4. `intent` must start with an imperative verb (map, design, implement, audit, scaffold).
5. `context.files` paths must be verified (ls/stat) before inclusion. No stale paths.
6. `constraints` cannot weaken Commandments. They can only add task-specific restrictions.

## Compact Mode

For simple handoffs where full context is unnecessary (e.g., Reviewer PASS → Herald):

<packet>
  <id>KRX-20260415-007</id>
  <origin>Reviewer</origin>
  <target>Herald</target>
  <state>REVIEWED</state>
  <intent>approve implementation</intent>
  <summary>All blueprint items verified. No security findings. PASS.</summary>
</packet>
```

### 5.2 NEW `rules/flow-states.md`

Replaces the implicit linear pipeline with an explicit state machine.

```markdown
# Flow States v1

> **Authority**: Defines every valid state in the Kerux development cycle,
> the transitions between them, and the failure handling at each stage.

## States

### IDLE
- **Owner**: Herald
- **Entry**: Session boot complete, or previous task committed/abandoned.
- **Exit**: User provides a task. → MAPPING

### MAPPING
- **Owner**: Tracker
- **Entry**: Herald dispatches a mapping packet.
- **Exit artifacts**: Context Packet (file paths, dependencies, project metadata).
- **Transitions**:
  - Success → DESIGNING
  - Failure (project not found, ambiguous scope) → IDLE + user clarification request

### DESIGNING
- **Owner**: Architect
- **Entry**: Context Packet received from Tracker.
- **Exit artifacts**: `spec_projeto.md` conforming to `templates/SPEC_TEMPLATE.md`.
- **Transitions**:
  - Success → SCAFFOLDING (if new project) or IMPLEMENTING (if existing project)
  - Failure (template unavailable, conflicting requirements) → IDLE + user escalation

### SCAFFOLDING
- **Owner**: Coder (under Herald supervision)
- **Entry**: New project flag set in spec. `lazygo.yml` generated.
- **Exit artifacts**: Project directory created, `spec_projeto.md` copied to root.
- **Transitions**:
  - Success → IMPLEMENTING
  - Failure (lazy.go error, path conflict) → DESIGNING + error context

### IMPLEMENTING
- **Owner**: Coder
- **Entry**: `spec_projeto.md` exists and is current.
- **Exit artifacts**: Modified/created files as listed in blueprint.
- **Transitions**:
  - Success → REVIEWING
  - Blocked (spec incomplete, missing dependency) → DESIGNING + BLOCKED packet
  - Failure (build error, test failure) → self-retry (max 2), then DESIGNING

### REVIEWING
- **Owner**: Reviewer
- **Entry**: Coder handoff packet with change summary.
- **Exit artifacts**: Verdict (PASS | REJECT | COMMENT) with evidence.
- **Transitions**:
  - PASS → STAGING
  - COMMENT → STAGING (with notes attached)
  - REJECT → route based on rejection type:
    - Blueprint deviation → IMPLEMENTING (Coder fixes)
    - Architectural flaw → DESIGNING (Architect revises)
    - Security finding → DESIGNING (mandatory — security issues require spec-level response)

### STAGING
- **Owner**: Herald
- **Entry**: Reviewer PASS or COMMENT.
- **Exit artifacts**: Commit message (Conventional Commits format).
- **Action**: Present diff summary + proposed commit message to user.
- **Transitions**:
  - User approves → COMMITTED
  - User requests changes → IMPLEMENTING
  - User abandons → IDLE

### COMMITTED
- **Owner**: Herald
- **Entry**: User explicit approval in current turn.
- **Action**: Execute git add + git commit. Never git push (user does this manually).
- **Transitions**: → IDLE

### FAILED
- **Owner**: Herald
- **Entry**: Any unrecoverable error (missing tools, corrupt state, repeated failures).
- **Action**: Report full error context to user. Suggest manual intervention.
- **Transitions**: → IDLE (after user acknowledgment)

## Invariants

1. Only one persona holds the lock at any time. No concurrent persona execution.
2. State transitions are logged in the session todo (skills/agent-todo.md).
3. Every REJECT must include the target state for the rollback — the Reviewer decides
   whether the fix is a Coder task or an Architect task.
4. STAGING → COMMITTED requires user approval in the CURRENT turn (Commandment C1).
5. The IMPLEMENTING → DESIGNING (BLOCKED) transition is the Coder's escape hatch
   for spec incompleteness. It is not a failure — it is correct behaviour.
```

### 5.3 MODIFY `rules/memory-rules.md`

**Change**: Absorb persistence protocol from `skills/agent-memory.md`. Add M5 section.

**Append after M4**:

```markdown
## 💾 M5: Persistence Protocol
Context is organized into persistent layers (follows the Layer model in M1):
- **Long-Term**: `.kerux/memory/lessons.md`. Confirmed user preferences and critical anti-patterns.
  Updated only on task completion or explicit user instruction.
- **Session**: `.kerux/memory/session.json`. Current state, task IDs, todo deltas, and transient variables.
  Updated at every state transition.

### Sync Rules
1. **Identify Delta**: What changed since the last state transition? (Handoff summary, lesson learned, task progress.)
2. **Update Session**: Sync task IDs and flow state to `session.json`.
3. **Commit Lessons**: Only on task completion or when a major preference is established.
4. **Runtime dependency**: When PERSISTENCE_MODE=none (see rules/runtime-contract.md),
   skip file writes. Maintain state in-context only.
```

**Remove from the file**: Any reference to specific token counts for specific models. Replace with:
- M3 threshold: "TOKEN_THRESHOLD as defined in `rules/runtime-contract.md` (default: 80% of available context)."

### 5.4 NEW `rules/error-taxonomy.md`

```markdown
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
```

### 5.5 NEW `rules/runtime-contract.md`

```markdown
# Runtime Contract v1

> Defines what Kerux requires from its execution environment.
> No file in .kerux/ may reference a specific LLM provider by name.
> Adaptation is handled through the variables below.

## Required Capabilities
1. File system read/write in the project workspace.
2. Shell command execution: ls, find, grep, cat, head, go, git.
3. Context window ≥ 100,000 tokens (operational minimum for boot + one task cycle).

## Optional Capabilities
4. Persistent file storage between sessions (enables session.json / lessons.md).
5. Web search (enables Tracker web intelligence skill).
6. Web browsing (enables visual verification skill).

## Adaptation Variables

### TOKEN_THRESHOLD
- **Purpose**: Trigger point for memory compression.
- **Default**: 80% of available context window.
- **Override**: Set at boot based on runtime detection.

### PERSISTENCE_MODE
- **Values**: `file` | `memory` | `none`
- `file`: Full persistence via session.json and lessons.md (local agent, IDE plugin).
- `memory`: Runtime provides cross-session memory (Claude memory system, similar).
  Skip file writes for lessons; use runtime memory API.
- `none`: Stateless. All context is in-session only. Memory compression produces
  a text block the user must paste into the next session.

### ATTENTION_HINTS
- **Purpose**: Model-specific markers to emphasize critical content.
- **Default**: Empty string (no hints).
- **Override**: Set per-runtime if the model supports attention directives.

## Boot Detection

At boot, `kerux-boot.md` probes the environment:
1. Check shell availability: `which go && which git`
2. Check file write: attempt to write/read a temp file in `.kerux/memory/`
3. If file write fails → PERSISTENCE_MODE=none
4. TOKEN_THRESHOLD: set by the hosting system or default to 80k.

Runtime detection is silent. The user sees only the boot greeting.
```

---

## 6. Blueprint — Phase 3: Revise Existing Files (P2–P4)

### 6.1 MODIFY `kerux.md`

**Changes**:

1. **Replace** §Organic Flow numbered list with a reference to `rules/flow-states.md`:
   ```
   ## Organic Flow (DevSecOps)
   The Herald routes traffic through the state machine defined in `rules/flow-states.md`.
   Each state has an owning persona, entry conditions, exit artifacts, and failure paths.
   ```

2. **Add** §Scope section after §Boot Sequence:
   ```
   ## Scope
   Kerux operates on project-scoped development tasks that produce version-controlled artifacts.
   The operational boundary:
   - IN: Code implementation, spec authoring, code review, scaffolding, commit preparation.
   - OUT: Production deployment, CI/CD pipeline execution, external service configuration,
     content authoring (blog posts, documentation outside the project).
   - BOUNDARY: Infrastructure-as-code changes are IN if they live in the project repo.
   ```

3. **Replace** §Traffic Protocol with a reference to `rules/packet-schema.md`:
   ```
   ## Traffic Protocol
   All inter-persona communication uses the packet format defined in `rules/packet-schema.md`.
   ```

4. **Add** to §Red Lines:
   ```
   - **NO UNDEFINED HANDOFFS**: Every persona transition must use a validated packet.
   - **NO STALE STATE**: If the flow state doesn't match the expected entry condition, halt and report.
   ```

5. **Update** version footer: `Kerux v3.0 | Consolidated`

### 6.2 MODIFY `personas/architect.md`

**Changes**:

1. **Replace** the absolute `file:///` SPEC template path with:
   ```
   3. **Template Compliance**: Read `.kerux/templates/SPEC_TEMPLATE.md`.
      Use it as the mandatory base for any `spec_projeto.md`.
   ```

2. **Add** to §Output:
   ```
   - **CI Mirror**: For every requirement in the spec, identify the corresponding
     CI check or test that enforces it. If no check exists, flag it as a gap.
   ```

### 6.3 MODIFY `personas/coder.md`

**Changes**:

1. **Add** §Escape Hatch after §Playbook step 4:
   ```
   5. **Blocked Path**: If during implementation the spec is incomplete, ambiguous,
      or contradicted by the actual codebase:
      - Do NOT improvise outside the blueprint.
      - Emit a BLOCKED packet to Herald with:
        - The specific spec section that is incomplete/wrong.
        - What the codebase actually requires.
        - A suggested amendment (informational, not authoritative — the Architect decides).
      - Flow transitions to DESIGNING. This is correct behaviour, not a failure.
   ```

2. **Modify** §Scaffolding step 1:
   ```
   - Verify lazy.go compatibility: run `go run main.go version` (or equivalent)
     to confirm the CLI accepts the generated lazygo.yml format.
   ```

3. **Renumber** subsequent steps.

### 6.4 MODIFY `personas/reviewer.md`

**Changes**:

1. **Add** §Supply Chain Audit after §Ops Audit:
   ```
   6. **Supply Chain Audit (SCA)**:
      - Verify no unsigned or unpinned dependencies were introduced.
      - Check for mutable tag references in CI workflows (e.g., `@latest`, `@main`).
      - If the project has a go.sum, verify it was updated consistently with go.mod.
      - Flag any new external dependency that lacks a documented justification in the spec.
   ```

2. **Add** §Envelope Compatibility after §Supply Chain Audit:
   ```
   7. **Envelope Compatibility**: If the modified code produces or consumes a JSON envelope
      (Vexil → Wardex → Vigil contract), verify:
      - Struct field additions are backward-compatible (new fields only, with omitempty).
      - No field was renamed or removed without a version bump.
      - The envelope validation function (ParseEnvelope pattern) still accepts the new shape.
   ```

3. **Modify** §Verdicts — add REJECT routing guidance:
   ```
   - **REJECT**: Fails a Commandment, deviates from the blueprint, or fails template/SCA compliance.
     The rejection MUST specify the target state:
     - REJECT → IMPLEMENTING: if the fix is a code-level correction within the existing spec.
     - REJECT → DESIGNING: if the fix requires a spec amendment (architectural flaw, security finding).
   ```

### 6.5 MODIFY `personas/tracker.md`

**Changes**:

1. **Expand** §Playbook with Go-specific intelligence:
   ```
   5. **Go Metadata**: If the project contains `go.mod`:
      - Parse module path and Go version.
      - List direct dependencies (exclude indirect).
      - Flag any dependency not in the stdlib that lacks a comment in go.mod.
   6. **Scaffold Metadata**: If the project contains `lazygo.yml`:
      - Parse project type, criticality level, and enabled features.
      - Include in the Context Packet — the Architect needs this for design decisions.
   7. **SPEC Header**: If the project contains `spec_projeto.md` or a SPEC file:
      - Extract version, status, and scope.
      - Report whether the spec is current or potentially stale (based on last-modified vs recent commits).
   ```

2. **Add** to §Accuracy Standard:
   ```
   - **Metadata freshness**: When reporting go.mod dependencies or SPEC status,
     always state the file's last-modified timestamp.
   ```

### 6.6 MODIFY `skills/kerux-boot.md`

**Changes**:

1. **Add** version check as step 0:
   ```
   0. **Version Check**: Read `.kerux/VERSION`. Log the orchestrator version.
      If VERSION is missing, warn the user and create it with `3.0.0`.
   ```

2. **Add** runtime detection as step 5 (before final greeting):
   ```
   5. **Runtime Detection**: Execute the probes defined in `rules/runtime-contract.md` §Boot Detection.
      Set PERSISTENCE_MODE and TOKEN_THRESHOLD silently.
   ```

3. **Modify** final greeting to include version:
   ```
   "Kerux Tech House v{VERSION} Online. Layers initialized. {PERSISTENCE_MODE} persistence active.
    What is the architecture we are building today?"
   ```

### 6.7 MODIFY `skills/kerux-dispatch.md`

**Changes**:

1. **Replace** §Dispatch Protocol step 2 with:
   ```
   2. **Packet Assembly**: Build a packet conforming to `rules/packet-schema.md`.
      Validate all required fields before dispatch. Reject malformed packets
      with a DEGRADED log entry — never send an incomplete handoff.
   ```

2. **Add** §State Validation:
   ```
   ## State Validation
   Before dispatching to any persona, Kerux verifies:
   1. The current flow state (from flow-states.md) allows this transition.
   2. The target persona's entry conditions are met.
   3. If conditions are not met, Kerux does NOT dispatch.
      Instead, it emits a BLOCKING error to the user explaining the gap.
   ```

### 6.8 MODIFY `skills/context-maintenance.md`

**Changes**:

1. **Remove** any reference to `[!!]` tags or model-specific attention mechanisms.
2. **Replace** hardcoded token thresholds with:
   ```
   - **Threshold**: If context exceeds TOKEN_THRESHOLD (rules/runtime-contract.md),
     invoke `skills/memory-compression.md`.
   ```
3. **Remove** the Gemini parenthetical in the current text.

### 6.9 MODIFY `skills/memory-compression.md`

**Changes**:

1. **Replace** §The Reset with a runtime-aware version:
   ```
   ## The Reset
   Once the Seed Block is created:
   - If PERSISTENCE_MODE=file: Write to `.kerux/memory/session.json`. Advise user to start new session.
   - If PERSISTENCE_MODE=memory: Summarize the seed block in-conversation. The runtime memory system
     will persist the relevant context.
   - If PERSISTENCE_MODE=none: Output the seed block as a fenced code block. Instruct the user
     to paste it at the start of the next session.
   ```

### 6.10 MODIFY `skills/code-review/protocol.md`

**Changes**:

Expand the minimal checklist into a structured protocol:

```markdown
# Skill: Code Review Protocol

Standard checklist for the Reviewer persona. Applied at REVIEWING state.

## Checklist

### Security Audit (Sec)
- [ ] No hardcoded secrets (Commandment C2).
- [ ] Input validation at trust boundaries.
- [ ] No shell injection vectors (exec.Command with user input).
- [ ] Cryptographic operations use stdlib (crypto/rand, crypto/subtle).

### Supply Chain Audit (SCA)
- [ ] No new unsigned or unpinned dependencies.
- [ ] No mutable tag references in CI workflows.
- [ ] go.sum consistent with go.mod.
- [ ] New external dependencies justified in spec.

### Logic Verification (Dev)
- [ ] Implementation matches spec_projeto.md blueprint exactly.
- [ ] No TODO/FIXME left as implementation (Commandment C3).
- [ ] Error handling: no silenced errors on security-critical paths.
- [ ] Tests cover the changed logic paths.

### Envelope Compatibility (if applicable)
- [ ] JSON struct changes are backward-compatible.
- [ ] No field renames/removals without version bump.
- [ ] ParseEnvelope-pattern validators updated.

### Ops Audit
- [ ] Config changes are logged and safe.
- [ ] No destructive operations without explicit user confirmation path.
```

---

## 7. Blueprint — Phase 4: New Files (P5)

### 7.1 NEW `VERSION`

```
3.0.0
```

Single line. Follows semver. Bumped on any structural change to `.kerux/`.

### 7.2 NEW `templates/SPEC_TEMPLATE.md`

**Action**: Copy the current `SPEC_TEMPLATE.md` from `~/Documentos/Projects/portfolio/SPEC_TEMPLATE.md` into `.kerux/templates/`. This is a relocation, not a creation — the content already exists.

If the source file is unavailable at implementation time, create a minimal placeholder that the owner fills manually:

```markdown
# SPEC_TEMPLATE.md
# TODO: Copy from portfolio/SPEC_TEMPLATE.md
# This placeholder must be replaced before the Architect can produce a spec.
```

Note: This is the ONE acceptable TODO in the system — it's a user action item, not an implementation gap.

### 7.3 NEW `memory/lessons.md`

```markdown
# Lessons

> Persistent preferences and confirmed anti-patterns. Updated only on task completion
> or explicit user instruction. Never updated speculatively.

<!-- Entries added below this line -->
```

### 7.4 NEW `memory/session.json`

```json
{
  "version": "3.0.0",
  "persistence_mode": "unknown",
  "current_state": "IDLE",
  "tasks": [],
  "last_updated": null
}
```

---

## 8. Acceptance Criteria

Each criterion maps to a Goal from §2.

| ID | Criterion | Verification Method |
|----|-----------|-------------------|
| AC-1 (G1) | `grep -r` across `.kerux/` returns no paragraph-level duplicate content. | Manual grep + diff after implementation. |
| AC-2 (G2) | Every persona file references `rules/packet-schema.md`. No persona defines its own handoff format. | Grep for "packet-schema" in all persona files. |
| AC-3 (G3) | `flow-states.md` defines transitions for: REJECT routing, BLOCKED (spec incomplete), SCAFFOLDING failure, concurrent access prevention. | Manual review of state machine completeness against FM-1 through FM-4 from the audit. |
| AC-4 (G4) | `grep -ri "gemini\|claude\|openai\|anthropic\|google"` across `.kerux/` returns zero hits (excluding this spec and the audit). | Automated grep. |
| AC-5 (G5) | Reviewer checklist includes SCA and envelope compatibility sections. Tracker playbook includes go.mod, lazygo.yml, and SPEC header parsing. | Manual review of persona files. |
| AC-6 (G6) | No file in `.kerux/` contains a `file:///` URI. SPEC template path is relative. | Grep for `file:///`. |
| AC-7 | All 4 files on the kill list are absent from the tree. | `ls` verification. |
| AC-8 | `VERSION` file exists and contains a valid semver string. | `cat .kerux/VERSION`. |

---

## 9. CI Mirror

This section maps spec requirements to verifiable checks. Since Kerux is not a code project with a CI pipeline, the "CI" equivalent is a post-implementation validation script.

```bash
#!/bin/bash
# .kerux/validate.sh — Post-implementation verification
# Run from project root containing .kerux/

set -euo pipefail

PASS=0
FAIL=0

check() {
  if eval "$2" > /dev/null 2>&1; then
    echo "  PASS: $1"
    ((PASS++))
  else
    echo "  FAIL: $1"
    ((FAIL++))
  fi
}

echo "=== Kerux Consolidation Validation ==="

# AC-1: No duplicate files
check "core-orchestration.md deleted"    "! test -f .kerux/skills/core-orchestration.md"
check "maintenance-skills.md deleted"    "! test -f .kerux/skills/maintenance-skills.md"
check "agent-memory.md deleted"          "! test -f .kerux/skills/agent-memory.md"
check "README.md deleted"               "! test -f .kerux/skills/README.md"

# AC-2: Packet schema referenced
check "architect refs packet-schema"     "grep -q 'packet-schema' .kerux/personas/architect.md"
check "coder refs packet-schema"         "grep -q 'packet-schema' .kerux/personas/coder.md"
check "reviewer refs packet-schema"      "grep -q 'packet-schema' .kerux/personas/reviewer.md"
check "dispatch refs packet-schema"      "grep -q 'packet-schema' .kerux/skills/kerux-dispatch.md"

# AC-4: No provider names
check "no provider names in tree"        "! grep -ri 'gemini\|anthropic\|openai' .kerux/ --include='*.md' | grep -v 'KERUX_AUDIT\|KERUX_CONSOLIDATION_SPEC'"

# AC-6: No file:/// URIs
check "no file:/// URIs"                 "! grep -r 'file:///' .kerux/ --include='*.md'"

# AC-8: VERSION exists
check "VERSION file exists"              "test -f .kerux/VERSION"
check "VERSION is semver"                "grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$' .kerux/VERSION"

# New files exist
check "flow-states.md exists"            "test -f .kerux/rules/flow-states.md"
check "packet-schema.md exists"          "test -f .kerux/rules/packet-schema.md"
check "error-taxonomy.md exists"         "test -f .kerux/rules/error-taxonomy.md"
check "runtime-contract.md exists"       "test -f .kerux/rules/runtime-contract.md"
check "templates dir exists"             "test -d .kerux/templates"
check "memory dir exists"               "test -d .kerux/memory"

echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
```

---

## 10. Implementation Order

Execute phases in sequence. Each phase is independently committable.

| Phase | Commit Message | Files Touched | Depends On |
|-------|---------------|---------------|------------|
| **P0** | `refactor(kerux): eliminate redundant skill files` | DELETE 4, MODIFY 2 (memory-rules.md, agent-todo.md) | Nothing |
| **P1** | `feat(kerux): define packet schema and flow state machine` | NEW 4 (packet-schema, flow-states, error-taxonomy, runtime-contract) | P0 |
| **P2** | `refactor(kerux): revise core files for state-aware flow` | MODIFY 7 (kerux.md, boot, dispatch, context-maintenance, compression, architect, coder) | P1 |
| **P3** | `feat(kerux): upgrade reviewer and tracker for ecosystem alignment` | MODIFY 3 (reviewer, tracker, code-review/protocol) | P1 |
| **P4** | `feat(kerux): add version tracking, templates dir, memory skeleton` | NEW 4 (VERSION, templates/SPEC_TEMPLATE, memory/lessons, memory/session.json) | P0 |

P2 and P3 are independent of each other — they can be parallelized or reordered.
P4 can run any time after P0.

---

*KERUX_CONSOLIDATION_SPEC v1.0.0 — End of document.*
