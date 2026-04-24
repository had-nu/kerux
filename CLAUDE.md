# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Kerux

Kerux is a **prompt-driven orchestrator** that lives inside a project as a `.kerux/` directory (similar to `.git/`). Its purpose is to coordinate complex development workflows using role-based context isolation within a **single LLM session**, replacing multi-agent platforms (CrewAI, AutoGen, LangGraph) without the cost of parallel model instances.

To activate it in any project: `Please comply with .kerux/kerux.md`

## Validation

```bash
./validate.sh