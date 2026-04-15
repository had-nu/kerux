#!/bin/bash
# .kerux/validate.sh — Post-implementation verification
# Run from the kerux/ directory

PASS=0
FAIL=0

check() {
  if eval "$2" > /dev/null 2>&1; then
    echo "  PASS: $1"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: $1"
    FAIL=$((FAIL + 1))
  fi
}

echo "=== Kerux Consolidation Validation ==="

# AC-1: No duplicate files
check "core-orchestration.md deleted"    "! test -f skills/core-orchestration.md"
check "maintenance-skills.md deleted"    "! test -f skills/maintenance-skills.md"
check "agent-memory.md deleted"          "! test -f skills/agent-memory.md"
check "README.md deleted"               "! test -f skills/README.md"

# AC-2: Packet schema referenced
check "architect refs packet-schema"     "grep -q 'packet-schema' roles/architect.md"
check "engineer refs packet-schema"         "grep -q 'packet-schema' roles/engineer.md"
check "auditor refs packet-schema"      "grep -q 'packet-schema' roles/auditor.md"
check "dispatch refs packet-schema"      "grep -q 'packet-schema' skills/kerux-dispatch.md"

# AC-4: No provider names
check "no provider names in tree"        "! grep -ri 'gemini\|anthropic\|openai' --include='*.md' . | grep -v 'KERUX_AUDIT\|KERUX_CONSOLIDATION_SPEC'"

# AC-6: No file:/// URIs
check "no file:/// URIs"                 "! grep -r 'file:///' --include='*.md' . | grep -v 'KERUX_AUDIT\|KERUX_CONSOLIDATION_SPEC'"

# AC-8: VERSION exists
check "VERSION file exists"              "test -f VERSION"
check "VERSION is semver"                "grep -qE '^[0-9]+\.[0-9]+\.[0-9]+$' VERSION"

# New files exist
check "flow-states.md exists"            "test -f rules/flow-states.md"
check "packet-schema.md exists"          "test -f rules/packet-schema.md"
check "error-taxonomy.md exists"         "test -f rules/error-taxonomy.md"
check "runtime-contract.md exists"       "test -f rules/runtime-contract.md"
check "templates dir exists"             "test -d templates"
check "SPEC_TEMPLATE.md exists"          "test -f templates/SPEC_TEMPLATE.md"
check "memory dir exists"               "test -d memory"
check "lessons.md exists"               "test -f memory/lessons.md"
check "session.json exists"             "test -f memory/session.json"

echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
