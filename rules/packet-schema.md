# Packet Schema v1

> **Authority**: This schema is the contract for all inter-persona communication.
> Every packet sent or received by any persona MUST conform to this structure.

## Schema

<packet>
  <id>Unique task identifier (format: KRX-YYYYMMDD-NNN)</id>
  <origin>Sending persona name</origin>
  <target>Receiving persona name</target>
  <state>Current flow state (reference: rules/flow-states.md)</state>
  <intent>Imperative verb phrase: what the receiver must do</intent>
  <context>
    <files>Ordered list of file paths relevant to this task</files>
    <dependencies>External requirements (tools, env vars, APIs)</dependencies>
    <constraints>Guardrails specific to this task — overrides nothing in Commandments</constraints>
  </context>
  <vars>
    Key=value pairs. Project name, target path, branch, flags.
  </vars>
  <summary>1-2 sentence description of what happened before this handoff</summary>
</packet>

## Validation Rules

1. `id` must be unique within a session. Kerux assigns IDs; personas do not.
2. `origin` and `target` must be valid persona names: Herald, Architect, Coder, Reviewer, Tracker.
3. `state` must be a valid state from `flow-states.md`.
4. `intent` must start with an imperative verb (map, design, implement, audit, scaffold).
5. `context.files` paths must be verified (ls/stat) before inclusion. No stale paths.
6. `constraints` cannot weaken Commandments. They can only add task-specific restrictions.

## Compact Mode

For simple handoffs where full context is unnecessary (e.g., Reviewer PASS → Herald):

<packet>
  <id>KRX-20260415-007</id>
  <origin>Reviewer</origin>
  <target>Herald</target>
  <state>REVIEWED</state>
  <intent>approve implementation</intent>
  <summary>All blueprint items verified. No security findings. PASS.</summary>
</packet>
