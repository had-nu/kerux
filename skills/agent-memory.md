# Skill: Agent Memory (agent-memory.md)

> **Objective**: Manage layered context and persistent knowledge to optimize the reasoning window.

## 💾 Layers (follows: `rules/memory-rules.md`)
- **Long-Term**: `.kerux/memory/lessons.md`. Persistent preferences/anti-patterns.
- **Project**: `spec_projeto.md`. Active objective and structural state.
- **Session**: `.kerux/memory/session.json`. Current state, task IDs, and transient variables.

## 📝 Persistence Protocol
1. **Identify Delta**: What changed since the last turn? (Handoff summary, lesson learned, task progress).
2. **Update Session**: Sync Task IDs and step counts to `session.json`.
3. **Commit Lessons**: If a task is complete or a major preference is established, update `lessons.md`.

## 🧹 Efficiency Protocol
- **Consolidation**: When context is heavy, invoke `skills/memory-compression.md`.
- **Garbage Collection**: Purge Task Context (Layer 3) after every major review PASS.
