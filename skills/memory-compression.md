# Skill: Memory Compression (memory-compression.md)

> **Objective**: Condense high-volume sessions into a "Seed Context" to reset the token window.

## 🌀 Compression Protocol
1. **Identify**: Extract the current goal and progress from `agent-todo.md`.
2. **Summarize**:
   - **Architect**: Summarize the current `spec_projeto.md`.
   - **Engineer**: State current file changes and uncommitted edits.
   - **Auditor**: List pending audit items.
3. **Capture Lessons**: Identify any new preferences or anti-patterns for `memory/lessons.md`.
4. **Seed Creation**: Assemble these summaries into a single "Seed Block" within `.kerux/memory/session.json`.

## 🔄 The Reset
Once the Seed Block is created:
- If PERSISTENCE_MODE=file: Write to `.kerux/memory/session.json`. Advise user to start new session.
- If PERSISTENCE_MODE=memory: Summarize the seed block in-conversation. The runtime memory system
  will persist the relevant context.
- If PERSISTENCE_MODE=none: Output the seed block as a fenced code block. Instruct the user
  to paste it at the start of the next session.
