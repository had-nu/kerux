# Skill: Code Review Protocol

Standard checklist for the Reviewer persona. Applied at REVIEWING state.

## Checklist

### Security Audit (Sec)
- [ ] No hardcoded secrets (Commandment C2).
- [ ] Input validation at trust boundaries.
- [ ] No shell injection vectors (exec.Command with user input).
- [ ] Cryptographic operations use stdlib (crypto/rand, crypto/subtle).

### Supply Chain Audit (SCA)
- [ ] No new unsigned or unpinned dependencies.
- [ ] No mutable tag references in CI workflows.
- [ ] go.sum consistent with go.mod.
- [ ] New external dependencies justified in spec.

### Logic Verification (Dev)
- [ ] Implementation matches spec_projeto.md blueprint exactly.
- [ ] No TODO/FIXME left as implementation (Commandment C3).
- [ ] Error handling: no silenced errors on security-critical paths.
- [ ] Tests cover the changed logic paths.

### Envelope Compatibility (if applicable)
- [ ] JSON struct changes are backward-compatible.
- [ ] No field renames/removals without version bump.
- [ ] ParseEnvelope-pattern validators updated.

### Ops Audit
- [ ] Config changes are logged and safe.
- [ ] No destructive operations without explicit user confirmation path.
