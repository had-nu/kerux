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
