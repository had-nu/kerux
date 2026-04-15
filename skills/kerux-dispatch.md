# Skill: Traffic (Dispatch, Handoff & Scaffolding)

> **Objective**: Efficiently route work between roles and automate project scaffolding via local ecosystem tools.

## 📡 Dispatch Protocol
1. **Selection**: Kerux selects the Role based on the current step in the Organic Flow.
2. **Packet Assembly**: Build a packet conforming to `rules/packet-schema.md`.
   Validate all required fields before dispatch. Reject malformed packets
   with a DEGRADED log entry — never send an incomplete handoff.
3. **Execution**: Invoke the Role's instructions.

## 🔍 State Validation
Before dispatching to any role, Kerux verifies:
1. The current flow state (from `rules/flow-states.md`) allows this transition.
2. The target role's entry conditions are met.
3. If conditions are not met, Kerux does NOT dispatch.
   Instead, it emits a BLOCKING error to the user explaining the gap.

## 🏗️ Scaffolding Protocol (lazy.go)
If a `spec_projeto.md` defines a new project:
1. **Detection**: Identify the new project name and module path from the spec.
2. **Configuration**: Generate a `lazygo.yml` based on the spec's technical inventory.
3. **Invoke Scaffolding**: 
   - Path: `/home/hadnu/Documentos/Projects/homelab/lazy.go`
   - Command: `go run main.go init --from <path_to_generated_yml>`
4. **Integration**: Save a copy of `spec_projeto.md` into the newly created project folder.

## 🤝 Handoff Protocol
- **State Preservation**: The Role must return its new state and any artifacts created.
- **Summary**: Always prefix output with a 1-2 sentence summary.

## 🔁 Review Loop
If the Role is `Engineer`:
1. Execute `Auditor` (The Guard).
2. If `Auditor.verdict == FAIL`, loop back to `Architect` or `Engineer` per `rules/flow-states.md` REJECT routing.
