# Skill: Agent Todo

Task tracking across state transitions and (when persistence allows) across sessions.

## Format

```
- [ ] Task description (Owner: {Role}, State: {state_when_created})
- [x] Completed task (Result: {one-line summary})
- [/] In progress (Blocker: {reason if blocked, else none})
- [!] Failed task (Reason: {why})
```

## Usage

Kerux (Lead) owns the todo list.

- Initialize at boot: load from `memory/session.json` under `tasks` field if PERSISTENCE_MODE=file.
- Append new task when a role emits a packet that implies a follow-up (e.g., Auditor COMMENT with improvement notes).
- Mark complete when the task is delivered or abandoned.
- Include the task delta in every transition packet's focus field when relevant.

## Task Lifecycle

1. Created: by user request or by a role's packet.
2. Pending [ ]: not yet started.
3. In progress [/]: a role is actively working on it.
4. Complete [x]: delivered, verified.
5. Failed [!]: abandoned with a reason. Never silently dropped.

## Storage

PERSISTENCE_MODE=file: `memory/session.json` field `tasks` (array).
PERSISTENCE_MODE=memory: runtime memory system.
PERSISTENCE_MODE=none: in-context only; lost on session end.

## Rules

- Never more than 10 active tasks per session. If exceeding, the user needs to prioritize.
- Tasks have one owner. Shared tasks are split.
- Completed tasks persist in the todo for audit trail until session compression.
- Failed tasks persist until user acknowledges the failure.
