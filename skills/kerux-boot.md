# Skill: Boot Protocol (kerux-boot.md)

> **Objective**: Establish the baseline state and load critical constraints before processing any user requests.

## 🚀 Boot Sequence
When Kerux is invoked in a new session, execute these steps silently before responding:

0. **Version Check**: Read `.kerux/VERSION`. Log the orchestrator version.
   If VERSION is missing, warn the user and create it with `1.0.0`.
1. **Load Memory Rules**: Ingest `rules/memory-rules.md` to define the context boundaries for this session.
2. **Check the Layers**: Verify if a `spec_projeto.md` exists in the local workspace layer.
3. **Template Sync**: Acknowledge the location of `templates/SPEC_TEMPLATE.md`.
4. **Read Edicts**: Load `rules/edicts.md` to internalize the overarching technical commandments.
5. **Runtime Detection**: Execute the probes defined in `rules/runtime-contract.md` §Boot Detection.
   Set PERSISTENCE_MODE and TOKEN_THRESHOLD silently.

## 🏁 Final Step
Once the boot sequence is complete, greet the User:
"Kerux Tech House v{VERSION} Online. Layers initialized. {PERSISTENCE_MODE} persistence active.
 What is the architecture we are building today?"
