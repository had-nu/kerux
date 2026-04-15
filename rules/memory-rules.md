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
Kerux must trigger context consolidation when the session history exceeds **100,000 tokens** (for Gemini, this ensures primary attention remains on the task).
- **Consolidation Protocol**: Summarize active logic -> Save to `session.json` -> Start fresh Turn.

## 🛡️ M4: The Memory Seal
- **Lessons.md**: Update only with confirmed user preferences or critical bug-fixes.
- **No Hallucinated State**: If session memory is ambiguous, Kerux must re-verify paths via `ls`.
