# Skill: Context Maintenance

Strategy for lean context. Runs continuously during flow.

## Three Layers

Layer 1: Core Directive — `.kerux/` files loaded at boot. Static. Never pruned.

Layer 2: Project Blueprint — `spec_projeto.md` + codebase map from Analyst. Persistent across state transitions in current flow. Pruned on IDLE transition.

Layer 3: Task Context — raw shell output, full file contents, tool results. Volatile. Pruned after each successful state transition.

## Pruning Rules

After every state transition where the packet's summary captured what matters:
- Discard raw `find` output if Analyst packet has `structure` field.
- Discard raw file reads if the relevant content is quoted in the packet's delta.
- Keep: spec, active code diffs, Reviewer findings.
- Drop: shell prompts, intermediate tool output, old transition logs.

Before any dispatch:
- Verify packet contains required fields.
- Drop anything from context that duplicates packet content — packet is now the carrier.

## Threshold

Two-stage response from `rules/runtime-contract.md`:

- At TOKEN_WARN (default 50000): aggressive Layer 3 pruning. Drop all shell
  output, intermediate tool results, and file reads already summarized in
  packets. Continue flow. Emit status line noting context pressure.
- At TOKEN_COMPACT (default 75000): invoke `memory-compression.md`. Halt flow
  until Seed Block is produced and Layer 3 is reset.

## Selective Reading

- Use `grep -n` for targeted line lookups.
- Use `head -N` or `sed -n 'A,Bp'` for ranges.
- Never `cat` files > 100 lines unless the whole content is needed.
- For file discovery, prefer `find . -name 'pattern' -type f` over recursive `ls`.

## Signal vs Noise

High signal (keep):
- Spec sections being actively worked
- Current role's packet input
- Recent code diffs
- Reviewer findings with file+line

Low signal (prune):
- Shell banners, prompts, ANSI codes
- Dependency download logs
- Build output without errors
- Previously-resolved transition logs

## Rules

- Pruning is the Lead's responsibility. Roles do not self-prune.
- When uncertain whether to prune, keep. A re-read is cheaper than a lost fact.
- Never prune spec_projeto.md during active flow.
- Never prune the current role's packet.
