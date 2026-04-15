# Edicts: Authoritative Guidance

> These rules apply within specific scopes. Deviations require justification.

## 🛠️ E1: Clean Architecture (Standard)
- Prefer composition over inheritance.
- Keep components focused and reusable.
- Any "Portfolio" project modification must respect the existing codebase style.

## 📝 E2: Commit Style
- Use [Conventional Commits](https://www.conventionalcommits.org/).
- Example: `feat: add kerux traffic skill` or `fix(wardex): calibrate score math`.

## 📂 E3: File Management
- Group related files in domain-specific directories.
- Avoid "god files" that handle multiple unrelated responsibilities.

## 🚦 E4: Review Loop
- All non-trivial code changes MUST be audited by the Reviewer persona.
- Rejection by the Reviewer triggers a mandatory redesign phase.
