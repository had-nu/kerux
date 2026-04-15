# Runtime Contract v1

> Defines what Kerux requires from its execution environment.
> No file in .kerux/ may reference a specific LLM provider by name.
> Adaptation is handled through the variables below.

## Required Capabilities
1. File system read/write in the project workspace.
2. Shell command execution: ls, find, grep, cat, head, go, git.
3. Context window ≥ 100,000 tokens (operational minimum for boot + one task cycle).

## Optional Capabilities
4. Persistent file storage between sessions (enables session.json / lessons.md).
5. Web search (enables Analyst web intelligence skill).
6. Web browsing (enables visual verification skill).

## Adaptation Variables

### TOKEN_THRESHOLD
- **Purpose**: Trigger point for memory compression.
- **Default**: 80% of available context window.
- **Override**: Set at boot based on runtime detection.

### PERSISTENCE_MODE
- **Values**: `file` | `memory` | `none`
- `file`: Full persistence via session.json and lessons.md (local agent, IDE plugin).
- `memory`: Runtime provides cross-session memory.
  Skip file writes for lessons; use runtime memory API.
- `none`: Stateless. All context is in-session only. Memory compression produces
  a text block the user must paste into the next session.

### ATTENTION_HINTS
- **Purpose**: Model-specific markers to emphasize critical content.
- **Default**: Empty string (no hints).
- **Override**: Set per-runtime if the model supports attention directives.

## Boot Detection

At boot, `kerux-boot.md` probes the environment:
1. Check shell availability: `which go && which git`
2. Check file write: attempt to write/read a temp file in `.kerux/memory/`
3. If file write fails → PERSISTENCE_MODE=none
4. TOKEN_THRESHOLD: set by the hosting system or default to 80k.

Runtime detection is silent. The user sees only the boot greeting.
