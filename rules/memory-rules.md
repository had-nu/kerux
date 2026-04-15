# Rules: Memory, Tokens & Context Management

> **Objective**: Ensure high-fidelity reasoning by maintaining a lean, high-signal context window.

## 💾 M1: Context Layering
Context is organized into three distinct layers to prevent "noise saturation":
- **Layer 1: Core Directive (Static)**: `.kerux` core instructions. Loaded at boot.
- **Layer 2: Project Blueprint (Persistent)**: `spec_projeto.md`. The source of truth for the current objective. 
- **Layer 3: Task Context (Volatile)**: Raw code and logs. Force-pruned after each successful sub-task completion.

## 🪙 M2: Token Discipline
- **Signal-to-Noise Ratio**: Prefer pseudocode for explanation and shorthand for handoffs.
- **Selective Reading**: Use `grep` and `sed` for targeted file reads instead of reading whole files > 300 lines.
- **Chunking**: Break complex reasoning into persona-specific segments to keep individual prompt sizes low.

## 🧹 M3: Adaptive Pruning
Kerux must trigger context consolidation when the session exceeds TOKEN_THRESHOLD
as defined in `rules/runtime-contract.md` (default: 80% of available context).
- **Consolidation Protocol**: Summarize active logic -> Save to `session.json` -> Start fresh Turn.

## 🛡️ M4: The Memory Seal
- **Lessons.md**: Update only with confirmed user preferences or critical bug-fixes.
- **No Hallucinated State**: If session memory is ambiguous, Kerux must re-verify paths via `ls`.

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
