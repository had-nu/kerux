# Skill: Boot Protocol (kerux-boot.md)

> **Objective**: Establish the baseline state and load critical constraints before processing any user requests.

## 🚀 Boot Sequence
When Kerux is invoked in a new session, you MUST execute these steps silently before responding:
1. **Load Memory Rules**: Ingest `rules/memory-rules.md` to define the context boundaries for this session.
2. **Check the Layers**: Verify if a `spec_projeto.md` exists in the local workspace layer. 
3. **Template Sync**: Acknowledge the location of the `SPEC_TEMPLATE.md`.
4. **Read Edicts**: Load `rules/edicts.md` to internalize the overarching technical commandments.

## 🏁 Final Step
Once the boot sequence is complete, ask the User:
"Kerux Tech House Online. Layers initialized. What is the architecture we are building today?"
