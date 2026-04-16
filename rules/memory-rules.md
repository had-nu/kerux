# Memory Rules

Context boundaries and persistence. Ensures lean, high-signal reasoning across sessions.

## M1: Context Layering

Three layers prevent noise saturation:

Layer 1 — Core Directive (static):
- `.kerux/` files loaded at boot.
- Commandments, edicts, flow-states, packet-schema, runtime-contract.
- Never pruned during active flow.

Layer 2 — Project Blueprint (persistent):
- `spec_projeto.md` + Analyst's context packet.
- Source of truth for the current objective.
- Pruned only on transition to IDLE.

Layer 3 — Task Context (volatile):
- Raw shell output, full file reads, intermediate tool results.
- Force-pruned after every successful state transition.

## M2: Token Discipline

- Signal-to-noise ratio is the primary metric. Pseudocode for explanation, shorthand for handoffs.
- Selective reading: `grep -n`, `head -N`, `sed -n 'A,Bp'` over `cat` for files > 100 lines.
- Chunk complex reasoning into role-specific segments.
- Packets carry delta, never full state — the receiving role reads from disk when needed.

## M3: Adaptive Pruning

Two thresholds from `rules/runtime-contract.md`:

- At TOKEN_WARN: Layer 3 pruning becomes aggressive. Discard all shell output,
  intermediate tool results, and file reads already summarized in transition
  packets. Layers 1 and 2 untouched. Status line emitted so the user sees
  context pressure.
- At TOKEN_COMPACT: Invoke `skills/memory-compression.md`. Produce Seed Block,
  reset Layer 3 entirely. Layers 1 and 2 survive.

The gap between WARN and COMPACT is the pruning window. If pruning at WARN
brings context back below WARN, no compression fires. If context continues to
grow past COMPACT despite pruning, compression is mandatory.

## M4: Memory Seal

- `memory/lessons.md`: update only on task completion or explicit user instruction.
- Never add speculative lessons. A single turn is not enough evidence.
- Never hallucinate state. If session memory is ambiguous, re-verify paths via `ls`.
- Lessons are append-only within a session. Revision happens at task boundaries.

## M5: Persistence Protocol

Persistence behaviour depends on PERSISTENCE_MODE (from runtime-contract.md):

PERSISTENCE_MODE=file:
- `memory/session.json`: current state, task IDs, todo deltas, transient vars.
- Updated at every state transition by Kerux.
- `memory/lessons.md`: persistent across sessions, append-only.

PERSISTENCE_MODE=memory:
- Runtime provides cross-session memory (e.g., hosted LLM memory API).
- Skip file writes. Use runtime memory primitives.
- Seed blocks and lessons live in runtime memory.

PERSISTENCE_MODE=none:
- Stateless. All context is in-session.
- Memory compression produces a text block the user pastes into the next session.
- No file writes to memory/ directory.

### Sync Rules (when persistence is active)

1. Identify delta: what changed since the last state transition?
2. Update session: sync current_state, task IDs, active spec to session.json.
3. Commit lessons: only on task completion or major preference established.
4. Verify: after write, re-read to confirm persistence succeeded. On failure, degrade to in-session only.
