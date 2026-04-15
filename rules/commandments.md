# Commandments: Absolute Laws

> These rules are the foundation of Kerux. They are never bypassed, nunca flexibilizadas.

## 🔴 C1: No Silent Mutations
- **NEVER** `git commit` without explicit, unambiguous user confirmation in the *current* turn.
- **NEVER** `git push --force` or delete remote branches without double confirmation.

## 🛡️ C2: Security First
- **Zero Hardcoded Secrets**: Use environment variables or placeholders.
- **No Shadow Logic**: Any logic that handles classification, risk, or sensitive data (like in `wardex`) must be explicitly documented and reviewed.

## 🏗️ C3: Structural Integrity
- **No Placeholders**: Never leave `TODO` comments as an implementation.
- **Maintain Docs**: Do not remove existing docstrings or comments unless explicitly refactoring them.

## 📉 C4: Token Discipline
- **Signal Only**: Avoid conversational filler in role outputs.
- **Compact Handoffs**: Use `<packet>` tags for structured data.
