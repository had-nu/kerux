# Skill: Kerux Boot (kerux-boot.md)

> **Objective**: Efficiently initialize the Kerux session.

1. **Self-Check**: Verify path `.kerux/` presence.
2. **Persistence**: Load `.kerux/memory/lessons.md` and `.kerux/memory/session.json`.
3. **Environment**: Identify active workspace and relevant `project-xiphos` or `wardex` context.
4. **Handoff**: Set session variables and yield to Herald for intent parsing.

---

# Skill: Kerux Dispatch (kerux-dispatch.md)

> **Objective**: Assembly and routing of persona packets.

1. **Synthesis**: Compile user request + filtered context into a `<packet>`.
2. **Routing**: Select persona based on complexity and authority level.
3. **Internal Log**: Record dispatch reasoning (cost/token estimation).

---

# Skill: Context Maintenance (context-maintenance.md)

> **Objective**: Keep the mental map of the project updated and lean.

1. **Topology Scan**: Periodic `find` to detect new files/directories.
2. **Pruning**: Identify and ignore `.log`, `.tmp`, or large binary files not already in `.gitignore`.
3. **Focus Mapping**: Update `CURRENT_DOMAIN` variable based on user interaction.
