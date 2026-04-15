# Skill: Context Maintenance (context-maintenance.md)

> **Objective**: Strategy for lean, high-signal context management and project mapping.

## 📁 Intelligence Layer
Kerux uses the following markers to map the domain:
1. **Core**: `.kerux/` and `.gitignore`.
2. **Scaffold**: `lazygo.yml` (MUST be parsed to understand project type, criticality, and features).
3. **Hierarchy**: Directory structure and entry points (e.g., `main.go`, `SPEC.md`).

## 📁 Hierarchy (follows: `rules/memory-rules.md`)
1. **Primary**: All files in the current active domain (determined by task).
2. **References**: Documentation identified by `Tracker`.

## 🧹 Pruning Logic
- **Delta Focus**: After every persona turn, discard raw task logs that have been summarized in the handoff.
- **Threshold**: If context exceeds TOKEN_THRESHOLD (`rules/runtime-contract.md`),
  invoke `skills/memory-compression.md`.
- **Ignore**: Always exclude heavy binary paths or known `.gitignore` matches.
