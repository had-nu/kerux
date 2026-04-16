# Skill: Memory Compression

Triggered when context ≥ TOKEN_COMPACT. Condenses session into a seed block.

## Trigger Conditions

Invoked by `context-maintenance.md` when threshold exceeded.
Also invokable manually by user: "compress context" or "reset session."

## Compression Protocol

### Step 1: Identify Active State
- Current flow state (from session.json or tracking).
- Active spec (path + version).
- Active role (which is holding the lock).
- Pending packet (if mid-transition).

### Step 2: Role-Specific Summarization

Analyst contribution:
- Summary of context packet already emitted. Drop raw mapping output.

Architect contribution:
- Summary of doc/spec_projeto.md — keep section headers and Blueprint item IDs, drop full pseudocode (reload from file when needed).

Engineer contribution:
- List of files modified with brief description. Drop full file contents (reload from disk when needed).

Auditor contribution:
- List of pending findings with file+line. Drop audit reasoning.

### Step 3: Lessons Extraction

Scan the session for:
- New user preferences expressed explicitly.
- Anti-patterns confirmed (e.g., "never use X library again").
- Critical bug-fix patterns that apply beyond this task.

Append confirmed lessons to `memory/lessons.md` if PERSISTENCE_MODE allows.

### Step 4: Seed Assembly

Produce a fenced Seed Block:

```
=== KERUX SEED BLOCK ===
version: 3.0.0
state: {current_state}
spec: {spec_path} v{spec_version}
active_role: {role}
pending_packet: {packet or none}

files_modified:
  - {path}: {brief description}

pending_findings:
  - {check_id}: {file} L{line} {finding}

lessons_delta:
  - {new lessons added this session}
=== END SEED ===
```

### Step 5: Reset per Runtime Mode

PERSISTENCE_MODE=file:
- Write Seed Block to `.kerux/memory/session.json` under `seed` field.
- Instruct user: "Session compressed. Continue here or start new session to reload from seed."

PERSISTENCE_MODE=memory:
- Rely on runtime's memory system to retain the seed summary.
- Instruct user: "Session compressed. Seed preserved in runtime memory."

PERSISTENCE_MODE=none:
- Output Seed Block as a fenced code block in chat.
- Instruct user: "Session compressed. Paste this Seed Block at the start of the next session to resume."

## On Resume

New session reads the Seed Block (from file, memory, or user paste):
1. Restore current_state from seed.
2. Restore spec path and version.
3. Acknowledge pending findings or packets without re-deriving them.
4. Resume flow from the state indicated.

## Rules

- Compression is destructive to Layer 3 (Task Context) only. Never discard spec, lessons, or commandments.
- Always produce a Seed Block before any reset, even in PERSISTENCE_MODE=file — the seed is the audit trail.
- Lessons must be confirmed, not speculative. Never add a lesson based on a single conversation turn.
- If compression is triggered mid-transition (packet pending), complete the transition first, then compress.
