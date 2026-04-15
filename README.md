# Kerux

Kerux is a prompt-driven orchestrator that lives inside your project as a `.kerux/` directory — the same way `.git/` manages version control. It runs coordinated development workflows inside a single LLM session: specialized roles (Analyst, Architect, Engineer, Auditor) coordinated by Kerux — the Lead — under strict rules, communicating through compressed handoffs, following a spec-first development cycle. One context window, one model, no platform fees.

## The problem

Multi-agent development platforms (CrewAI, AutoGen, LangGraph, Devin) run multiple LLM instances in parallel, each with its own context window, each billing tokens separately. A four-agent pipeline where each agent consumes 50k tokens costs you 200k tokens per task — plus orchestration overhead, API routing, and the serial latency of agent-to-agent calls.

If you're working with a single-agent tool — Antigravity, Claude Code, Cursor, Codex — you can't run those platforms. You get one model, one session, one context window. But you still need the coordination: someone to map the codebase before writing code, someone to design before implementing, someone to audit after implementing. Without that structure, the model drifts — it refactors files it shouldn't touch, skips security checks, forgets the spec existed.

Kerux gives you multi-agent coordination on a single-agent budget.

## How it works

Kerux is a `.kerux/` directory inside your project root. It contains markdown files that define how the LLM should coordinate development work. When you load `kerux.md` at the start of a session, the model reads the roles, rules, and skills, then operates as a structured development team instead of an unconstrained assistant. The files define:

**Roles** — scoped functions with specific inputs, outputs, and constraints. Kerux acts as the Lead — it coordinates the flow and routes work. The Analyst maps the codebase. The Architect produces a spec. The Engineer implements against the spec. The Auditor verifies the implementation. Each role operates within defined boundaries and hands off to the next through a structured packet.

**Rules** — non-negotiable constraints that apply across all roles. No silent git commits. No hardcoded secrets. No implementation without a spec. No bypassing an Auditor rejection. These prevent the failure mode where the model quietly does something destructive because nothing told it not to.

**A state machine** — instead of a linear pipeline that only handles the happy path, Kerux defines explicit states (MAPPING, DESIGNING, IMPLEMENTING, REVIEWING, STAGING) with transitions for success, rejection, and failure. When the Auditor rejects, the flow routes back to the Engineer or the Architect depending on the type of defect. When the Engineer discovers the spec is incomplete, it signals BLOCKED instead of improvising.

**Compressed inter-role packets** — because all roles share the same context window, handoffs don't need to carry data. The data is already visible. A packet is a state transition marker that tells the next role where to look and what to do, in ~20 tokens instead of ~180. On a constrained context window, this matters.

## Roles

Each role maps to a real engineering function. Kerux itself is the Lead — it doesn't write code or specs, it coordinates who does what and when.

| Role | Function | Real-world equivalent |
|------|----------|----------------------|
| **Kerux** (Lead) | Coordinates the flow, routes work between roles, enforces rules, manages state transitions. | Tech Lead / Engineering Manager |
| **Analyst** | Maps the codebase, gathers project metadata (`go.mod`, `lazygo.yml`), produces context packets. | Systems Analyst / CTI Analyst |
| **Architect** | Designs the system, produces the `spec_projeto.md`, defines the blueprint for implementation. | Solutions Architect |
| **Engineer** | Implements code against the spec. Owns the build, tests, and scaffolding. | Software Engineer |
| **Auditor** | Verifies implementation against spec and security baseline. Issues PASS, REJECT, or COMMENT verdicts. | Security Auditor / QA Auditor |

Roles are additive. New roles (Writer, Tester, Ops) can be introduced by adding a file to `roles/` and registering the role in the state machine — no existing files need to change.

## Where the tokens go

In a typical development task with Kerux, token consumption breaks down roughly like this:

| Category | Proportion | What it is |
|----------|-----------|------------|
| Boot (`.kerux/` files) | ~15% | Role definitions, rules, skills loaded at session start. |
| Spec + context | ~20% | The `spec_projeto.md` and codebase context from Analyst mapping. |
| Implementation | ~45% | The Engineer generating actual code. This is the productive spend. |
| Coordination | ~10% | Packets, handoffs, state transitions between roles. |
| Review + rework | ~10% | Auditor verification and any rejection-triggered rework cycles. |

Kerux targets the 15% (boot) and 10% (coordination) categories. The implementation cost is irreducible — you need those tokens to write code. But the overhead is compressible.

**Boot cost reduction**: The `.kerux/` instruction files are optimized for AI consumption, not human readability. Compressed prose, no decorative formatting, no filler. Human-readable originals are kept as `.original.md` for editing. Estimated savings: 40–50% of boot tokens.

**Coordination cost reduction**: Inter-role packets use state transition markers instead of verbose XML or natural language summaries. An Auditor PASS handoff is `→K|REV→STG|all checks pass, 0 findings|stage for commit` — one line, ~20 tokens. The equivalent in prose would cost 80–150 tokens. Across 5–8 transitions per task, this compounds.

**Rework cost reduction**: The state machine routes Auditor rejections to the right role (Engineer for code bugs, Architect for design flaws) instead of restarting the full pipeline. A targeted fix costs less than a full re-implementation.

## Architecture

```
.kerux/
├── kerux.md                    ← Entry point. Kerux is the Lead.
├── VERSION
│
├── roles/
│   ├── analyst.md              ← Codebase intelligence
│   ├── architect.md            ← Spec authoring
│   ├── engineer.md             ← Implementation
│   └── auditor.md              ← Security + quality verification
│
├── rules/
│   ├── commandments.md         ← Absolute laws (never bypassed)
│   ├── edicts.md               ← Authoritative guidance
│   ├── memory-rules.md         ← Context management
│   ├── flow-states.md          ← State machine definition
│   ├── packet-schema.md        ← Handoff format
│   ├── error-taxonomy.md       ← Failure handling
│   └── runtime-contract.md     ← Environment abstraction
│
├── skills/
│   ├── kerux-boot.md           ← Session initialization
│   ├── kerux-dispatch.md       ← Role routing
│   ├── context-maintenance.md  ← Context pruning
│   ├── memory-compression.md   ← Session reset protocol
│   └── ...
│
├── templates/
│   └── SPEC_TEMPLATE.md        ← Blueprint template for specs
│
└── memory/
    ├── lessons.md              ← Persistent preferences
    └── session.json            ← Session state
```

## The development cycle

```
User request
    │
    ▼
MAPPING ──── Analyst maps the codebase, parses go.mod, locates the spec
    │
    ▼
DESIGNING ── Architect produces or validates spec_projeto.md
    │
    ▼
IMPLEMENTING ── Engineer writes code against the spec
    │         │
    │         └── BLOCKED? → back to DESIGNING (spec incomplete)
    │
    ▼
REVIEWING ── Auditor verifies implementation against spec + security baseline
    │         │
    │         ├── PASS → STAGING (prepare commit for user approval)
    │         ├── REJECT (code bug) → back to IMPLEMENTING
    │         └── REJECT (design flaw) → back to DESIGNING
    │
    ▼
STAGING ──── User reviews diff and approves
    │
    ▼
COMMITTED ── git commit (never push — user does that)
```

Every transition produces a packet. Every packet fits on one line.

## Getting started

`.kerux/` lives inside the project it coordinates — the same way `.git/` lives inside the repository it tracks. Each project gets its own `.kerux/` with its own memory, lessons, and session state.

**Option A: scaffold with lazy.go**

```bash
lazy.go init --kerux
```

This generates the project structure with `.kerux/` included and pre-configured for the project type and criticality level defined in `lazygo.yml`.

**Option B: manual setup**

Copy the `.kerux/` directory into your project root. Then start your LLM session and load the orchestrator:

```
Please comply with .kerux/kerux.md
```

Kerux boots, loads its rules, detects the runtime, and asks what you're building.

## Runtime compatibility

Kerux is a set of text files. It runs on any LLM agent that can read files and execute shell commands.

| Runtime | Tested | Notes |
|---------|--------|-------|
| Antigravity (Gemini) | Yes | Primary development environment. Full persistence support. |
| Claude Code | Planned | Memory system maps to `PERSISTENCE_MODE=memory`. |
| Cursor / Windsurf | Planned | Skill-file compatible via project rules. |
| Codex (OpenAI) | Untested | Should work — same file-reading + shell model. |

The `runtime-contract.md` file defines what Kerux needs from the environment. If the runtime can read files, run `go`, and run `git`, Kerux works.

## Design decisions

**Spec-driven development.** Nothing gets implemented without a `spec_projeto.md`. The spec is authored by the Architect (or the user), and every downstream decision traces back to it. The Auditor verifies implementation *against the spec*, not against vibes. The SPEC template lives inside `.kerux/templates/` — portable with the project, never an external dependency.

**Security as a flow property, not a gate.** The Auditor doesn't run at the end as a checkbox. It runs as a state in the machine, with the authority to reject and route back. Security findings that require spec changes go to the Architect, not the Engineer. This prevents the pattern where a developer patches a security bug locally without addressing the design flaw that caused it.

**Token economy is a first-class concern.** Every file in `.kerux/` is written with token cost in mind. Compressed packets, pruned context, tiered verbosity (internal communication is terse; user-facing output is normal). The system includes a benchmark spec for empirically measuring packet format efficiency — because the claim "this saves tokens" should be backed by data, not assumptions.

**No vendor lock-in.** Kerux doesn't import libraries, call APIs, or depend on any specific model's features. It's prompt engineering with structure. The `.kerux/` directory travels with the project — open it with Antigravity today, Claude Code tomorrow, Cursor next week. The `runtime-contract.md` abstracts the three things that vary between environments: persistence, token budget, and attention hints.

## What Kerux is not

It's not a framework. There's no `npm install`, no binary, no runtime dependency. It's a directory of instructions that lives inside your project.

It's not a standalone tool. Kerux doesn't have its own repo or CLI. It's embedded — like `.git/`, `.github/`, or `.vscode/`. The project owns the orchestrator, not the other way around.

It's not a multi-agent platform. There's one model, one session, one context window. The "agents" are roles — behavioural modes that the same model switches between, governed by rules that constrain each mode.

It's not autonomous. Every git mutation requires explicit user approval in the current turn. The model proposes; the user disposes.

## Project status

Kerux is under active development. Two internal specs govern the current work:

- `KERUX_CONSOLIDATION_SPEC.md` — structural overhaul: eliminating redundancy, formalizing the state machine, abstracting the runtime layer.
- `KERUX_PACKET_BENCHMARK_SPEC.md` — empirical evaluation of inter-role communication formats (compressed markers vs. compact JSON) to validate token savings with measured data.

Planned: integration as a scaffold option in [lazy.go](https://github.com/had-nu/lazy.go), so new projects can include `.kerux/` from day one.

Current version: see `.kerux/VERSION`.

## License

[TBD by owner]
