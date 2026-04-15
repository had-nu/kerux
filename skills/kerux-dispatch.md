# Skill: Traffic (Dispatch & Handoff)

> **Objective**: Efficiently route work between personas and automate project scaffolding via local ecosystem tools.

## 📡 Dispatch Protocol
1. **Selection**: Kerux selects the Persona based on the current step in the Organic Flow.
2. **Packet Assembly**: Build a `<packet>` with `intent`, `context`, and `vars`.
3. **Execution**: Invoke the Persona's instructions.

## 🤝 Handoff Protocol
- **State Preservation**: The Persona must return its new state and any artifacts created.
- **Summary**: Always prefix output with a 1-2 sentence summary.

## 🔁 Review Loop
If the Persona is `Coder`:
1. Execute `Reviewer` (The Guard).
2. If `Reviewer.verdict == FAIL`, loop back to `Architect` or `Coder`.
