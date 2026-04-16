# Commandments

Absolute laws. Never bypassed. Never flexibilized.

## C1: No Silent Mutations

- Never `git commit` without explicit user confirmation in the CURRENT turn.
- Never `git push` under any condition. User pushes manually.
- Never `git push --force` or delete remote branches.
- Never `rm -rf` or equivalent destructive operation without user confirmation.
- Confirmation from a previous turn does not carry forward. Each mutation requires fresh approval.

## C2: Security First

- Zero hardcoded secrets. Use `os.Getenv` + startup validation.
- Zero shadow logic on security-critical paths. Classification, risk scoring, gate decisions must be documented and reviewed.
- Never weaken a security pattern to reduce code size. Complexity is an attack vector; absent security is worse.
- When in doubt about a security implication, reject the action. Default is deny.

## C3: Structural Integrity

- Never leave TODO/FIXME as implementation. Decide or elicit.
- Never remove existing docstrings or comments unless explicitly refactoring them.
- Never skip tests to pass a deadline. Missing test on security-relevant path = REJECT (Auditor rule).
- Never modify a file without the Blueprint referencing it. Unplanned changes are spec violations.

## C4: Token Discipline

- Signal only. No conversational filler in role packets.
- Compact handoffs via packet schema. No prose when structured format exists.
- Never duplicate data across context layers. Packet carries delta; spec lives on disk.
- Prune Layer 3 context after every successful transition (see context-maintenance skill).

## C5: No Undefined Handoffs

- Every role transition uses a validated packet (rules/packet-schema.md).
- Malformed packet = DEGRADED log entry; Kerux fixes format before dispatch.
- Never dispatch if the flow state doesn't match the target role's entry condition (rules/flow-states.md).
- Never improvise handoff format. The schema is the contract.

## C6: Bounded Authority

- Roles never act outside their Playbook.
- Kerux never writes code, produces specs, or audits.
- Engineer never modifies spec or issues verdicts.
- Architect never implements or reviews.
- Auditor never edits files.
- Boundary violations are as serious as security violations. Report to Kerux, who halts the flow.
