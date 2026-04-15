# sigcheck — Technical Specification
<!-- Version: 1.0 | Status: Active | Author: hadnu | Date: 2026-04-15 -->

> **RFC 2119 Convention:** The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHOULD", "SHOULD NOT", "MAY", and "OPTIONAL" in this document are to be interpreted as described in [RFC 2119](https://www.ietf.org/rfc/rfc2119.txt).

> **Context**: This spec is a _control document_ for the Kerux Packet Benchmark (KERUX_PACKET_BENCHMARK_SPEC v1.0.0). It is authored once and used identically in both experimental arms. Neither arm's Architect modifies it.

---

## 1. Overview

### 1.1 Problem Statement
Verifying file integrity after distribution requires computing and comparing cryptographic hashes against a known-good manifest. Existing tools (`sha256sum`) lack structured error handling, path traversal protection, and programmatic exit codes suitable for CI pipelines.

### 1.2 Proposed Solution
`sigcheck` is a Go CLI tool that generates SHA-256 manifest files from a target directory and verifies file integrity against those manifests. It provides structured exit codes (0/1/2), path traversal protection, streaming hash computation, and symlink-safe directory walking.

### 1.3 Scope
- **IN**: SHA-256 manifest generation, manifest verification, CLI interface via Cobra, security hardening.
- **OUT**: Other hash algorithms, remote file verification, GUI, daemon mode, manifest signing (PKI).

---

## 2. Goals & Non-Goals

### 2.1 Goals
1. Generate a SHA-256 manifest from all regular files in a target directory (recursive).
2. Verify all files listed in a manifest against their actual SHA-256 hashes.
3. Report mismatches and missing files with clear, actionable output.
4. Exit with deterministic codes: `0` = all match, `1` = mismatch/missing, `2` = operational error.
5. Protect against path traversal, symlink following, and large-file DoS.
6. Pass all automated checks: `go vet`, `staticcheck`, `gosec`, `go test -race`.

### 2.2 Non-Goals
1. Supporting hash algorithms other than SHA-256.
2. Concurrent hashing with goroutines (simplicity over performance for v1).
3. Manifest signing or PKI integration.
4. Watch mode or daemon operation.

---

## 3. Architecture

### 3.1 System Diagram

```
CLI (Cobra)
├── cmd/root.go        → Entry point, persistent flags
├── cmd/generate.go    → generate subcommand
├── cmd/verify.go      → verify subcommand
│
├── internal/hasher/
│   └── hasher.go      → Streaming SHA-256 computation
│
└── internal/manifest/
    └── manifest.go    → Manifest read/write/parse with validation
```

### 3.2 Component Inventory

| Component | Responsibility | Technology | Notes |
|-----------|---------------|------------|-------|
| `cmd/root.go` | CLI entry point, persistent flags | Cobra | Sets up root command |
| `cmd/generate.go` | Walk directory, compute hashes, write manifest | Cobra + internal | Uses hasher + manifest |
| `cmd/verify.go` | Read manifest, recompute hashes, report results | Cobra + internal | Uses hasher + manifest |
| `internal/hasher/hasher.go` | Stream file contents into SHA-256, return hex digest | `crypto/sha256`, `io` | MUST use streaming — no `os.ReadFile` |
| `internal/manifest/manifest.go` | Parse/write manifest format, validate paths per line | `bufio`, `strings`, `filepath` | Each line: `<hex>  <path>` (sha256sum compat) |

### 3.3 Data Flow

**Generate flow:**
1. User runs `sigcheck generate --target ./dist/ --output checksums.sha256`
2. `cmd/generate.go` calls `filepath.WalkDir` on `--target`
3. For each regular file (skip symlinks, skip directories): call `hasher.HashFile(path)`
4. `hasher.HashFile` opens the file, streams via `io.Copy` into `sha256.New()`, returns hex digest
5. Collect `hash  relativePath` pairs
6. Write pairs to `--output` via `manifest.Write()`
7. Exit 0 on success, exit 2 on any error

**Verify flow:**
1. User runs `sigcheck verify --manifest checksums.sha256 --target ./dist/`
2. `cmd/verify.go` calls `manifest.Parse(manifestPath)` → returns `[]Entry{Hash, Path}`
3. `manifest.Parse` validates each line: format check, path validation via `safePath()`
4. For each entry: compute `hasher.HashFile(filepath.Join(target, entry.Path))`
5. Compare computed hash with manifest hash
6. Collect results: `MATCH`, `MISMATCH`, `MISSING`, `ERROR`
7. Print results to stdout
8. Exit 0 if all MATCH, exit 1 if any MISMATCH or MISSING, exit 2 if any ERROR

### 3.4 External Dependencies

| Dependency | Version | Purpose | Risk if Unavailable |
|-----------|---------|---------|---------------------|
| `github.com/spf13/cobra` | latest stable | CLI framework | Cannot build — hard dependency |
| Go stdlib (`crypto/sha256`, `io`, `filepath`, `bufio`) | go1.22+ | Core logic | N/A — stdlib |

---

## 4. Data Model

### 4.1 Core Entities

```go
// internal/manifest/manifest.go

// Entry represents one line in the manifest file.
type Entry struct {
    Hash string // Lowercase hex-encoded SHA-256 digest (64 chars)
    Path string // Relative file path (validated, no traversal)
}
```

### 4.2 Storage
Manifest files are plain text, one entry per line, format: `<64-char-hex>  <relative-path>` (two spaces, sha256sum compatible).

### 4.3 Data Lifecycle
Manifest files are generated on demand. No database. No retention policy — user manages lifecycle.

---

## 5. Interfaces

### 5.1 CLI / API Surface

```
sigcheck generate --target <dir> --output <file>
sigcheck verify   --manifest <file> --target <dir>
```

| Command | Flag | Type | Required | Default | Description |
|---------|------|------|----------|---------|-------------|
| `generate` | `--target` | string | yes | — | Directory to hash |
| `generate` | `--output` | string | yes | — | Manifest output path |
| `verify` | `--manifest` | string | yes | — | Manifest file to verify against |
| `verify` | `--target` | string | yes | — | Base directory for file resolution |

### 5.2 Configuration
No configuration files. No environment variables. All inputs via CLI flags.

### 5.3 Output Formats

**Generate**: Writes manifest file. No stdout output on success.

**Verify**: Prints verification results to stdout:
```
OK   path/to/file1.txt
FAIL path/to/file2.txt (hash mismatch)
MISS path/to/file3.txt (file not found)
```

Exit codes: `0` (all OK), `1` (any FAIL or MISS), `2` (operational error).

---

## 6. Performance & Capacity

### 6.1 Targets

| Metric | Target | Boundary Condition |
|--------|--------|---------------------|
| Throughput | ≥ 100 MB/s hashing | Measured on local SSD |
| Memory | ≤ 32 KB per file | Streaming hash, no full-file read |

### 6.2 Bottlenecks & Limits
- Single-threaded file walking and hashing (non-goal: concurrency in v1).
- IO-bound on large directories. Acceptable for the benchmark scope.

### 6.3 Scaling Strategy
Not applicable. Single-user CLI tool.

---

## 7. Security & Compliance

### 7.1 Threat Model

| ID | Threat | Attack Vector | Control |
|----|--------|--------------|---------|
| T1 | Path traversal | Malicious manifest contains `../../etc/passwd` | `safePath()` validates every path resolves under `--target` |
| T2 | Large file DoS | Target dir contains a 100GB file | Streaming hash via `io.Copy` — O(1) memory |
| T3 | Symlink escape | Symlink in target dir points outside base | `WalkDir` callback checks `d.Type()`, skips symlinks |
| T4 | Manifest injection | Malformed lines in manifest | `manifest.Parse` rejects malformed lines with error |
| T5 | TOCTOU | File changes between hash and report | **Known limitation.** Documented. No mitigation in v1. |

### 7.2 Data Handling
No secrets. No PII. All data is file hashes and paths. No network communication.

### 7.3 Compliance Requirements
None. Internal tooling.

---

## 8. Deployment & Operations

### 8.1 Infrastructure
Local binary. `go build -o sigcheck .`

### 8.2 Build & Release
```bash
go build -o sigcheck .
```
No CI pipeline for this benchmark project. Manual build.

### 8.3 Monitoring & Observability
Not applicable. CLI tool with stdout output.

### 8.4 Disaster Recovery
Not applicable.

---

## 9. Testing Strategy

### 9.1 Unit Tests

| Package | Test | Description |
|---------|------|-------------|
| `internal/hasher` | `TestHashFile` | Hash a known file, compare to pre-computed digest |
| `internal/hasher` | `TestHashFile_NotFound` | Returns error on non-existent file |
| `internal/manifest` | `TestParse` | Parse valid manifest, verify entries |
| `internal/manifest` | `TestParse_Malformed` | Reject lines with wrong format |
| `internal/manifest` | `TestParse_Traversal` | Reject entries with `..` path components |
| `internal/manifest` | `TestWrite` | Write entries, verify format matches sha256sum |
| `internal/manifest` | `TestSafePath` | Valid and invalid paths against base dir |

### 9.2 Integration Tests
Not applicable for this scope. Unit + E2E covers the surface.

### 9.3 End-to-End / Acceptance Tests

| ID | Test | Steps | Expected |
|----|------|-------|----------|
| F1 | Generate produces valid manifest | Create tmpdir with 3 files → run `generate` → parse output | Manifest has 3 entries, hashes match `sha256sum` output |
| F2 | Verify detects tampered file | Generate → modify one file → run `verify` | Exit 1, output shows FAIL for modified file |
| F3 | Verify detects missing file | Generate → delete one file → run `verify` | Exit 1, output shows MISS for deleted file |
| F4 | Verify succeeds on clean dir | Generate → verify immediately | Exit 0, all OK |
| F5 | Error on non-existent target | Run `generate --target /nonexistent` | Exit 2, error message on stderr |

### 9.4 Performance / Load Tests
Not applicable for benchmark scope.

---

## 10. Milestones & Deliverables

| Phase | Deliverable | Success Criteria | Target Date |
|-------|------------|------------------|-------------|
| 1 | `internal/hasher` + tests | `TestHashFile` passes, streaming verified | — |
| 2 | `internal/manifest` + tests | Parse/Write/SafePath tests pass | — |
| 3 | `cmd/` (Cobra CLI) | All E2E tests (F1–F5) pass | — |
| 4 | Security audit | All S1–S8 checks pass | — |
| 5 | Tooling checks | `go vet`, `staticcheck`, `gosec`, `go test -race` all exit 0 | — |

---

## 11. Risks & Mitigations

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| TOCTOU race in verify | False positive/negative | Low (single-user tool) | Document as known limitation |
| Cobra version mismatch | Build failure | Low | Pin version in go.mod |

---

## 12. Decision Log

| ID | Decision | Rationale | Date | Status |
|----|----------|-----------|------|--------|
| D-001 | No concurrency | Simplicity for benchmark; v1 scope | 2026-04-15 | Accepted |
| D-002 | sha256sum-compatible manifest format | Interoperability, familiar format | 2026-04-15 | Accepted |
| D-003 | Streaming hash only | Prevents large-file DoS, O(1) memory | 2026-04-15 | Accepted |
| D-004 | Skip symlinks silently | Safest default; no symlink escape | 2026-04-15 | Accepted |

---

## 13. Open Questions

- [x] All questions resolved for benchmark control spec.

---

## Appendices

### A. Glossary

| Term | Definition |
|------|-----------|
| Manifest | Plain text file mapping SHA-256 hashes to relative file paths |
| Streaming hash | Computing a hash by reading a file in chunks, not loading it entirely into memory |
| Path traversal | Attack where `../` sequences escape intended directory boundaries |

### B. References

- `KERUX_PACKET_BENCHMARK_SPEC v1.0.0` — Experiment design
- `SPEC_TEMPLATE.md` — Template this spec follows

---

## Blueprint

### Files to create

#### [NEW] `go.mod`
Module: `github.com/had-nu/sigcheck`. Go 1.22+. Dependency: `github.com/spf13/cobra`.

#### [NEW] `main.go`
Entry point. Calls `cmd.Execute()`.

#### [NEW] `cmd/root.go`
Root command setup. No persistent flags beyond help.

#### [NEW] `cmd/generate.go`
```pseudocode
func runGenerate(cmd, args):
    target = flag("target")
    output = flag("output")
    
    entries = []
    err = filepath.WalkDir(target, func(path, d, err):
        if err: return err
        if d.IsDir(): return nil
        if d.Type() is symlink: return nil  // T3: skip symlinks
        
        relPath = relativize(path, target)
        hash = hasher.HashFile(path)
        entries.append(Entry{hash, relPath})
        return nil
    )
    if err: exit(2)
    
    manifest.Write(output, entries)
    exit(0)
```

#### [NEW] `cmd/verify.go`
```pseudocode
func runVerify(cmd, args):
    manifestPath = flag("manifest")
    target = flag("target")
    
    entries, err = manifest.Parse(manifestPath, target)  // validates paths
    if err: exit(2)
    
    hasFailure = false
    for entry in entries:
        fullPath = filepath.Join(target, entry.Path)
        if not exists(fullPath):
            print("MISS", entry.Path)
            hasFailure = true
            continue
        
        computed = hasher.HashFile(fullPath)
        if computed != entry.Hash:
            print("FAIL", entry.Path, "(hash mismatch)")
            hasFailure = true
        else:
            print("OK", entry.Path)
    
    if hasFailure: exit(1)
    exit(0)
```

#### [NEW] `internal/hasher/hasher.go`
```pseudocode
func HashFile(path string) (string, error):
    f = os.Open(path)
    defer f.Close()
    
    h = sha256.New()
    io.Copy(h, f)  // T2: streaming, O(1) memory
    
    return hex.EncodeToString(h.Sum(nil)), nil
```

#### [NEW] `internal/manifest/manifest.go`
```pseudocode
type Entry struct { Hash, Path string }

func Parse(manifestPath, baseDir string) ([]Entry, error):
    scanner = bufio.NewScanner(open(manifestPath))
    entries = []
    for scanner.Scan():
        line = scanner.Text()
        parts = split line by "  " (two spaces)
        if len(parts) != 2: return error("malformed line")
        hash, path = parts[0], parts[1]
        if len(hash) != 64: return error("invalid hash length")
        if not safePath(path, baseDir): return error("path traversal: " + path)  // T1
        entries.append(Entry{hash, path})
    return entries, nil

func Write(outputPath string, entries []Entry) error:
    f = create(outputPath)
    for entry in entries:
        fmt.Fprintf(f, "%s  %s\n", entry.Hash, entry.Path)
    return nil

func safePath(path, baseDir string) bool:
    // T1: Resolve absolute, verify prefix under baseDir
    abs = filepath.Join(baseDir, path)
    resolved = filepath.Clean(abs)
    return strings.HasPrefix(resolved, filepath.Clean(baseDir))
```

### Guardrails (for Reviewer)

| ID | Check | Pass condition |
|----|-------|----------------|
| S1 | Path traversal protection | `safePath()` used on all manifest-derived paths |
| S2 | Streaming hash | No `os.ReadFile` or `io.ReadAll` in hasher. Must use `io.Copy` → `sha256.New()` |
| S3 | Symlink handling | `WalkDir` skips symlinks via `d.Type()` check |
| S4 | Error propagation | No `_` on error returns in hash, file open, manifest parse |
| S5 | Exit codes | 0=match, 1=mismatch/missing, 2=error. Tested in E2E. |
| S6 | No hardcoded paths | All paths from CLI flags |
| S7 | Manifest parsing | Malformed lines return error, not silently skipped |
| S8 | Concurrency safety | N/A for v1 (sequential). If goroutines introduced, require WaitGroup. |

### CI Mirror

```bash
go vet ./...
staticcheck ./...
gosec -quiet ./...
go test -race -count=1 ./...
```
