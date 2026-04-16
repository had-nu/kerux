# Skill: Dispatch Protocol

Role routing and packet validation. Kerux (Lead) uses this to coordinate transitions.

## Dispatch Sequence

1. Receive intent (from user or from incoming role packet).
2. Determine current flow state (from session or in-context tracking).
3. Validate transition: check `rules/flow-states.md` — is this transition legal from current state?
   - Legal → proceed to step 4.
   - Illegal → do NOT dispatch. Report to user: "Cannot transition from {current} to {target}. Reason: {entry condition not met}."
4. Select target role based on flow state:
   - IDLE → user request → Analyst (MAP)
   - MAP complete → Architect (DES)
   - DES complete → Engineer (SCF if new project, IMP if existing)
   - SCF complete → Engineer (IMP)
   - IMP complete → Auditor (REV)
   - REV PASS → Kerux stages (STG)
   - REV REJECT → route per Auditor's verdict: Engineer (IMP) or Architect (DES)
   - STG approved → Kerux commits (COM)
   - BLOCKED from Engineer → Architect (DES)
5. Validate outgoing packet format: `→{target}|{from→to}|{delta}|{focus}`
   - target: single letter (K, N, A, E, U) — Kerux, aNalyst, Architect, Engineer, aUditor
   - from→to: 3-letter state codes
   - delta: ≤80 chars, what changed
   - focus: ≤60 chars, what target should prioritize
   - Documentation: `.kerux/rules/packet-schema.md`
   - If packet malformed → log DEGRADED, fix format, proceed.
6. Dispatch: activate target role with packet as opening context.

## State Abbreviations

IDL=IDLE MAP=MAPPING DES=DESIGNING SCF=SCAFFOLDING IMP=IMPLEMENTING REV=REVIEWING STG=STAGING COM=COMMITTED FAI=FAILED

## Role Codes

K=Kerux(Lead) N=Analyst A=Architect E=Engineer U=Auditor

## Status Reporting

At each transition, emit a one-line status to the user:

```
[KRX] MAP→DES | Analyst mapped 12 files, go1.22, cobra project | Architect designing spec
```

This is the user's visibility into the flow. Keep it to one line. No detail — the user can ask if they want more.

## Concurrency Guard

Only one role active at a time. If a transition is requested while another role holds the lock, queue it. In practice with single-session LLMs, this is enforced by the sequential nature of the conversation — but the rule exists for clarity.

## Session Tracking

After each transition, update session state:
- If PERSISTENCE_MODE=file: write current_state to `.kerux/memory/session.json`
- If PERSISTENCE_MODE=none: maintain in-context only.
