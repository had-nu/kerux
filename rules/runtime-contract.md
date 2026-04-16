# Runtime Contract

Defines what Kerux requires from its execution environment.
No file in `.kerux/` references a specific LLM provider by name.
Adaptation happens through the variables below.

## Required Capabilities

1. File system read/write in project workspace.
2. Shell command execution: `ls`, `find`, `grep`, `cat`, `head`, `sed`, `go`, `git`.
3. Context window ≥ 100,000 tokens (operational minimum: boot + one full task cycle).

If any required capability is missing: FATAL at boot.

## Optional Capabilities

4. Persistent file storage between sessions (enables session.json, lessons.md).
5. Web search (enables Analyst web intelligence when needed).
6. Web browsing (enables visual verification for web projects).

Missing optional capabilities → DEGRADED, not FATAL.

## Adaptation Variables

### TOKEN_WARN

- **Purpose**: Preemptive pruning signal. Layer 3 pruning becomes aggressive.
- **Default**: 50000.
- **Override**: Set at boot based on runtime detection.
- **Rationale**: LLM output quality degrades well before the context is full
  (measured ~40% for complex tasks, Huntley 2025). Pruning begins here to
  keep active work in the "smart zone."

### TOKEN_COMPACT

- **Purpose**: Hard trigger for memory compression (seed block + reset).
- **Default**: 75000.
- **Override**: Set at boot based on runtime detection.
- **Rationale**: Last safe moment to compress before degradation is severe.
  The gap between WARN and COMPACT (~25k tokens) is the pruning window where
  aggressive Layer 3 pruning keeps the session operational without full reset.

### PERSISTENCE_MODE

Values:

**`file`** — Full file persistence.
- `memory/session.json` and `memory/lessons.md` live on disk.
- Typical for local agents, IDE plugins, CLI tools.
- Enables resume across sessions via filesystem.

**`memory`** — Runtime-provided cross-session memory.
- Runtime has a memory API (hosted LLM with memory, or similar).
- Skip file writes. Use runtime primitives.
- Seed blocks and lessons live in runtime memory.

**`none`** — Stateless.
- No cross-session state.
- Memory compression produces a paste-block for user.
- All context lost at session end unless user copies.

### ATTENTION_HINTS

- **Purpose**: Model-specific markers to emphasize critical content.
- **Default**: empty string (no hints).
- **Override**: set per-runtime if model supports attention directives.
- **Usage**: prepend/append to high-priority instructions. Silent no-op when unsupported.

## Boot Detection

Performed silently during boot sequence in `kerux.md` §Boot:

```bash
# 1. Shell availability (FATAL if missing)
which go && which git

# 2. File write capability (determines PERSISTENCE_MODE)
echo test > .kerux/memory/.probe && rm .kerux/memory/.probe
# Success → PERSISTENCE_MODE=file
# Fail (read-only FS, no memory dir) → PERSISTENCE_MODE=none

# 3. Token thresholds
# If runtime exposes context window size:
#   TOKEN_WARN    = 0.50 * window
#   TOKEN_COMPACT = 0.75 * window
# Else: defaults 50000 / 75000

# 4. Attention hints
# Set only if runtime documentation confirms support
# Default empty
```

## Invariants

1. No role file mentions a specific LLM provider by name.
2. No hardcoded token counts outside this file. Files reference TOKEN_WARN or TOKEN_COMPACT by name.
3. Runtime detection is silent. User sees only the boot greeting.
4. Degradation is reported in status lines, not logged silently.
5. The contract is the only abstraction layer. Downstream files consume the variables, not the runtime.

## Extension Points

When a new runtime is added:
1. Document detection in `kerux.md` §Boot Step 3 probe section.
2. Set PERSISTENCE_MODE, TOKEN_WARN / TOKEN_COMPACT, ATTENTION_HINTS based on the runtime's capabilities.
3. Do not modify role or skill files. The contract handles the adaptation.
