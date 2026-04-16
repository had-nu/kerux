# Skill: Boot Protocol

Session initialization. Execute silently before responding to user.

## Sequence

0. Read `.kerux/VERSION`. Log version. If missing, warn user, create with `1.0.0`.
1. Load `rules/commandments.md` — absolute constraints.
2. Load `rules/edicts.md` — authoritative guidance.
3. Load `rules/memory-rules.md` — context boundaries.
4. Load `rules/flow-states.md` — state machine.
5. Load `rules/packet-schema.md` — handoff format.
6. Check `spec_projeto.md` in project root. If present, note version + status.
7. Check `.kerux/templates/SPEC_TEMPLATE.md` exists. If missing, flag as degraded — Architect cannot produce specs.
8. Runtime detection (silent):
   - Shell: `which go && which git` — if missing, FATAL.
   - Persistence probe: attempt write/read to `.kerux/memory/test_probe`. 
     - Success → PERSISTENCE_MODE=file, delete probe.
     - Fail → PERSISTENCE_MODE=none.
   - TOKEN_WARN / TOKEN_COMPACT: defaults 50000 / 75000. Override if runtime
     provides context window size (0.50 / 0.75 of window respectively).
9. Load `.kerux/memory/session.json` if exists and PERSISTENCE_MODE=file.
10. Load `.kerux/memory/lessons.md` if exists.

## Boot Greeting

After sequence completes:

```
Kerux v{VERSION} | {PERSISTENCE_MODE} persistence | State: {current_state from session.json or IDLE}
What are we building?
```

## Failure Modes

- `go` or `git` not in PATH → FATAL. Report and stop.
- SPEC_TEMPLATE missing → DEGRADED. Boot continues, but Architect cannot produce specs until resolved.
- session.json corrupt → warn, reset to IDLE, continue.
- lessons.md missing → normal (first session). Continue.
