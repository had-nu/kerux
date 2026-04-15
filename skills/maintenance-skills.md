# Skill: Agent Todo (agent-todo.md)

> **Objective**: Persistent task tracking across sessions.

## 📝 Format
- [ ] Task Description (Owner: PersonaName)
- [x] Completed Task (Result: Summary)
- [/] In-Progress Task

## 🛠 Usage
- **Kerux**: Initialize at boot.
- **Coder/Architect**: Update when state changes.
- **Handoff**: Include the delta of the todo list.

---

# Skill: Context Maintenance (context-maintenance.md)

> **Objective**: Strategy for lean context management.

## 📁 Hierarchy
1. **Core**: `.kerux/` and `.gitignore`.
2. **Structural**: Directories and file names (`ls -R`).
3. **Active**: Files currently being modified and their dependencies.

## 🧹 Pruning Logic
- If context exceeds 1M tokens, summarize "Active" files and discard their raw content from the current prompt.
- Tag important code blocks with `[!!]` for Gemini's attention.
