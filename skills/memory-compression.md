# Skill: Memory Compression (memory-compression.md)

> **Objective**: Condense high-volume sessions into a "Seed Context" to reset the token window.

## 🌀 Compression Protocol
1. **Identify**: Extract the current goal and progress from `agent-todo.md`.
2. **Summarize**:
   - **Architect**: Summarize the current `spec_projeto.md`.
   - **Coder**: State current file changes and uncommitted edits.
   - **Reviewer**: List pending audit items.
3. **Capture Lessons**: Identify any new preferences or anti-patterns for `lessons.md`.
4. **Seed Creation**: Assemble these summaries into a single "Seed Block" within `.kerux/memory/session.json`.

## 🔄 The Reset
Once the Seed Block is created:
1. Advise the user to start a new chat session.
2. In the new session, Kerux's boot sequence will prioritize the Seed Block to resume with full context efficiency.
