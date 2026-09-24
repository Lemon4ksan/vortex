# Handoff Report: Milestone 3 Complete Worker Execution Plan & Synthesis

- **Agent**: `explorer_m3_3` (teamwork_preview_explorer)
- **Role**: Explorer & Synthesizer
- **Target**: Milestone 3 — Benchmark-Grade Code Documentation & Architecture
- **Workspace**: `d:/CodingProjects/vortex`
- **Output Report**: `d:/CodingProjects/vortex/.agents/explorer_m3_3/handoff.md`
- **Timestamp**: 2026-09-22T20:15:00Z

---

## 1. Observation

### 1.1 Direct Baseline Tool Commands & Output Verifications
1. **Workspace Package Enumeration**:
   - Command: `go list -f '{{.Dir}}: {{.Name}}' ./...`
   - Result: 41 total Go packages discovered across `d:/CodingProjects/vortex`.
     - 1 root/contract package: `ast` (`package ast`)
     - 1 CLI main package: `cmd/vortex` (`package main`)
     - 11 internal packages: `internal/ast`, `internal/base`, `internal/borrow`, `internal/core`, `internal/inspector`, `internal/oracle`, `internal/perf`, `internal/spec`, `internal/text`, `internal/traffic`, `internal/workspace`.
     - 28 packages in `pkg/`: `analysis`, `asyncapi`, `builder`, `cache`, `cfg`, `diff`, `emitter`, `enum`, `git`, `history`, `ingest`, `ir`, `jsbundle`, `lint`, `merge`, `mirror`, `openapi`, `optimizer`, `oracle/gen` (`package gen`), `oracle/spec` (`package spec`), `parser`, `patcher`, `pipeline`, `project`, `spec`, `sys`, `tuple`, `version`.

2. **Existing Documentation Inventory**:
   - `find_by_name` for `doc.go` returned 17 files:
     - 1 in `ast/doc.go` (108 lines; pristine Aoni reference with ASCII diagrams, quickstart, tiers, and `[Type]` links).
     - 9 in `internal/*/doc.go` (missing in `internal/borrow` and `internal/inspector`).
     - 7 in `pkg/*/doc.go`:
       - 3 complete reference files: `pkg/analysis/doc.go` (31 lines), `pkg/asyncapi/doc.go` (49 lines), `pkg/openapi/doc.go` (53 lines).
       - 4 stub files: `pkg/emitter/doc.go` (10 lines), `pkg/ir/doc.go` (14 lines), `pkg/optimizer/doc.go` (20 lines), `pkg/parser/doc.go` (10 lines).
     - 21 missing `doc.go` files in `pkg/`.

3. **Current Error Handling State**:
   - `find_by_name` for `*error*.go` returned only 1 file: `pkg/lint/rules_error.go` (lint rules validating error return types).
   - Zero exported package sentinel errors (`var Err... = errors.New(...)`) exist in `pkg/` or `internal/`.
   - Zero typed error structs implementing `error` exist in `pkg/` or `internal/`.
   - Ad-hoc error creation observed across key subsystems:
     - `pkg/pipeline/ast.go:54, 63, 158, 240, 476`: `return errors.New("target file is required")` repeated verbatim 5 times.
     - `pkg/diff/stack.go:284, 379, 438, 453, 934`: `errors.New("stack is empty")`, `errors.New("at least 2 frames required in stack for adjacent diff")`.
     - `pkg/git/git.go:65, 125, 223, 270, 300, 324`: `fmt.Errorf("git show %s failed: ...")`, `errors.New("not a git repository (or any parent directory)")`.
     - `pkg/project/config.go:414, 782, 858, 1121`: `errors.New("nil configuration")`, `errors.New("no .vortex.work workspace file found in hierarchy")`.
     - `pkg/cache/traffic_store.go:208`: `fmt.Errorf("vortex: traffic session %q not found in cache", idOrHash)`.

4. **Go Runtime & Standard Library Capabilities**:
   - Command: `go version` -> `go version go1.27.0 windows/amd64`.
   - Standard Library `errors.AsType[E error](err error) (E, bool)` confirmed available in Go 1.27.
   - Reference implementations inspected:
     - `d:/CodingProjects/aoni/errors.go` (282 lines): `var Err... = errors.New(...)`, typed structs, and `errors.AsType` predicates.
     - `d:/CodingProjects/foundation/net/dns/error.go` (120 lines): `IsNotFound`, `ResolutionError`, `errors.AsType[*ResolutionError]`.

5. **Test and Lint Baseline Executions**:
   - Command: `go test -count=1 ./...`
     - Result: **PASS** (Exit code 0). All 41 packages pass cleanly.
   - Command: `golangci-lint run --allow-parallel-runners ./...`
     - Result: **PASS** (Exit code 0, 0 issues reported).

6. **Inputs from Peer Explorers**:
   - `explorer_m3_1` (`d:/CodingProjects/vortex/.agents/explorer_m3_1/handoff.md`): Comprehensive Godoc designs for all 27 `pkg/` packages and 2 `internal/` packages, including ASCII diagrams, 3 usage tiers, and concurrency/zero-alloc guarantees.
   - `explorer_m3_2` (`d:/CodingProjects/vortex/.agents/explorer_m3_2/handoff.md`): Concrete implementations of 8 `errors.go` files with typed error structs, sentinels, predicates, and backwards compatibility analysis.

---

## 2. Logic Chain

1. **Alignment with Aoni & Foundation Sovereign Standards**:
   - *Observation*: The user request (R3) dictates: "Align Vortex package documentation with aoni's forever-frozen core standards: Introduce comprehensive doc.go files across all undocumented packages in pkg/ ... with usage tiers, ASCII diagrams, and godoc links. Standardize typed error predicates and sentinel errors across the toolchain."
   - *Logic*: Every newly created `doc.go` must adhere to the 7-part Aoni blueprint (License Header, One-Line Package Summary, Architecture Overview with ASCII diagram, Core Building Blocks with `[Type]` links, Usage Tiers 1-3, Concurrency & Thread-Safety Guarantees, and Performance/Zero-Alloc Profiles).

2. **Deduction of the 43 Target Files**:
   - *Observation*: The codebase contains 21 packages in `pkg/` with no `doc.go`, 4 packages in `pkg/` with 10-20 line stub `doc.go` files, 2 packages in `internal/` with no `doc.go`, 8 key subsystems with no error types, and no test coverage for error predicates.
   - *Logic*: The implementation worker must produce exactly:
     - 21 brand new `doc.go` files in `pkg/`
     - 4 expanded `doc.go` files in `pkg/`
     - 2 new `doc.go` files in `internal/`
     - 8 new `errors.go` files in `pkg/`
     - 8 new companion `errors_test.go` files in `pkg/`
     - Total: 43 files.

3. **Subsystem Error Symmetry (`Op, Path, Key, Err`)**:
   - *Observation*: Each subsystem currently reports disparate error formats.
   - *Logic*: Structuring all 8 subsystem errors with uniform fields (`Op`, `Path`, `Key`, `Err`) provides architectural consistency, enables uniform CLI logging, and maintains clean unwrapping via `Unwrap() error { return e.Err }`.

4. **Go 1.27 High-Performance Error Predicates**:
   - *Observation*: Go 1.27 introduces generic `errors.AsType[E error](err error) (E, bool)`.
   - *Logic*: Implementing `func Is<Predicate>(err error) bool` using `errors.Is(err, ErrSentinel)` and `errors.AsType[*<Subsystem>Error](err)` replaces reflection-based `errors.As` with direct, type-safe traversal. Fallback checks (e.g. `os.ErrNotExist`) guarantee backwards compatibility with existing stdlib errors.

5. **Phased Worker Batching**:
   - *Observation*: Modifying 43 files in a single unbatched pass risks build breaks and difficult debugging if a syntax or lint issue occurs.
   - *Logic*: Partitioning implementation into 6 sequential batches (Batch 1: Error Architecture & Tests, Batch 2: Core Stubs & Internal Docs, Batch 3: Compilation & Toolchain Docs, Batch 4: Spec & Interop Docs, Batch 5: System & Infra Docs, Batch 6: Final Verification) ensures incremental compilation, rapid error localization, and zero test regressions.

---

## 3. Caveats

1. **Explorer Read-Only Constraint**:
   - Per role constraints, Explorer agents must not write or modify production files in `pkg/` or `internal/`. The worker (`worker_m3`) will execute the file creations specified in this plan.
2. **Backwards Compatibility Invariant**:
   - Existing functions currently returning `fmt.Errorf` or `errors.New` must continue to function without breaking existing callers or tests. The new sentinel errors and error structs are exported additions, with predicates designed to recognize both new sentinels and underlying cause errors.
3. **Godoc Heading & Link Syntax**:
   - Go 1.19+ doc comments require headings to start with `#` on their own line preceded by a blank line, and type references to be enclosed in brackets `[TypeName]`. Unmatched brackets or improper indentation can trigger godoc rendering bugs or linter warnings.
4. **No Emoji / Fluff**:
   - In accordance with R2 and the sovereign aesthetic guidelines, all documentation and error strings must be completely free of informal emojis (`⚡`, `✨`, `🔴`, `🚀`, etc.), utilizing only restrained Unicode glyphs (`✔`, `✖`, `◆`, `↳`, `—`) or standard ASCII.

---

## 4. Conclusion & Worker Execution Plan

### 4.1 Master File Inventory Checklist (43 Files Total)

| # | Relative Path | Package | Action | Scope / Key Content |
|---|---|---|---|---|
| **Group A: 21 New `doc.go` Files in `pkg/`** | | | | |
| 1 | `pkg/builder/doc.go` | `builder` | Create | IR Builder, AST assembly pipeline, fixture population, file watching |
| 2 | `pkg/cache/doc.go` | `cache` | Create | 3-pillar cache: LintCache, SecretsVault (AES-256-GCM), TrafficStore (gzip) |
| 3 | `pkg/cfg/doc.go` | `cfg` | Create | Control Flow Graph, basic blocks, dataflow analysis, reachability |
| 4 | `pkg/diff/doc.go` | `diff` | Create | AST diffing, 3-way structural diffs, checkpoint stack snapshots |
| 5 | `pkg/enum/doc.go` | `enum` | Create | OpenAPI/TS/HAR enum inference, type-safe enum constants, validators |
| 6 | `pkg/git/doc.go` | `git` | Create | In-memory Git inspection, worktree state, branch tracking, zero tempfiles |
| 7 | `pkg/history/doc.go` | `history` | Create | Version evolution timeline, undo/redo journal ledger, semantic version bump |
| 8 | `pkg/ingest/doc.go` | `ingest` | Create | HAR 1.2 network capture ingestion, OpenAPI 3.x document synthesis |
| 9 | `pkg/jsbundle/doc.go` | `jsbundle` | Create | Minified JS/TS bundle scanner, Protobuf & JSPB schema extraction |
| 10 | `pkg/lint/doc.go` | `lint` | Create | 5-category contract linter (E, W, P, S, B), automated fixes, registry |
| 11 | `pkg/merge/doc.go` | `merge` | Create | 3-way specification merging, conflict resolution, schema reconciliation |
| 12 | `pkg/mirror/doc.go` | `mirror` | Create | `@mirror` upstream contract synchronization, shadow drift detection |
| 13 | `pkg/oracle/gen/doc.go` | `gen` | Create | JS/Playwright oracle code emitter, Go RPC bridge generator |
| 14 | `pkg/oracle/spec/doc.go` | `spec` | Create | Declarative test harness validation specs, DOM & network assertions |
| 15 | `pkg/patcher/doc.go` | `patcher` | Create | Non-destructive in-place Go AST patching preserving comments & formatting |
| 16 | `pkg/pipeline/doc.go` | `pipeline` | Create | High-level compilation pipeline, AST deobfuscation, enum extraction |
| 17 | `pkg/project/doc.go` | `project` | Create | `.vortex.yml` workspace manifest parsing, contract lifecycle, status engine |
| 18 | `pkg/spec/doc.go` | `spec` | Create | Single-source-of-truth directive registry, argument scopes, pipeline stages |
| 19 | `pkg/sys/doc.go` | `sys` | Create | OS runtime capability inspection (AVX2, AVX-512, io_uring), CPU affinity |
| 20 | `pkg/tuple/doc.go` | `tuple` | Create | Positional JSON array & Protobuf/JSPB tuple schema inference |
| 21 | `pkg/version/doc.go` | `version` | Create | Vortex toolchain version metadata, semantic build tags, capability flags |
| **Group B: 4 Stub `doc.go` Expansions in `pkg/`** | | | | |
| 22 | `pkg/emitter/doc.go` | `emitter` | Expand | DTO encoders, client structs, sub-requesters, harness gen, ASCII pipeline |
| 23 | `pkg/ir/doc.go` | `ir` | Expand | 4 hierarchical scopes (Service, Method, Param, Return), type invariants |
| 24 | `pkg/optimizer/doc.go` | `optimizer` | Expand | Sub-requester clustering, 64-byte stack sizing, query canonicalization |
| 25 | `pkg/parser/doc.go` | `parser` | Expand | Directive parsing (@service, @get, etc.), lexer, AST traversal, binder |
| **Group C: 2 Internal `doc.go` Files** | | | | |
| 26 | `internal/borrow/doc.go` | `borrow` | Create | Linear resource borrow checker (`analysis.Analyzer`), Acquire/Release |
| 27 | `internal/inspector/doc.go` | `inspector` | Create | Live HTTP capture dashboard, WebSocket telemetry, off-heap ring buffer |
| **Group D: 8 `errors.go` Implementations** | | | | |
| 28 | `pkg/project/errors.go` | `project` | Create | Sentinels (`ErrWorkspaceNotFound`, etc.), `ProjectError`, `IsNotFound`, `IsStale` |
| 29 | `pkg/parser/errors.go` | `parser` | Create | Sentinels (`ErrSyntaxError`, etc.), `ParseError`, `IsSyntaxError`, `IsNotFound` |
| 30 | `pkg/diff/errors.go` | `diff` | Create | Sentinels (`ErrStackEmpty`, etc.), `DiffError`, `IsNotFound`, `IsConflict` |
| 31 | `pkg/git/errors.go` | `git` | Create | Sentinels (`ErrNotRepository`, etc.), `GitError`, `IsNotRepository`, `IsNotFound` |
| 32 | `pkg/cache/errors.go` | `cache` | Create | Sentinels (`ErrSessionNotFound`, etc.), `CacheError`, `IsNotFound`, `IsCorrupt` |
| 33 | `pkg/lint/errors.go` | `lint` | Create | Sentinels (`ErrLintFailure`, etc.), `LintError`, `IsLintFailure`, `IsNotFound` |
| 34 | `pkg/spec/errors.go` | `spec` | Create | Sentinels (`ErrSpecNotFound`, etc.), `SpecError`, `IsNotFound`, `IsUnsupportedFormat` |
| 35 | `pkg/pipeline/errors.go` | `pipeline` | Create | Sentinels (`ErrTargetFileRequired`, etc.), `PipelineError`, `IsPipelineAborted`, `IsNotFound` |
| **Group E: 8 `errors_test.go` Unit Test Suites** | | | | |
| 36 | `pkg/project/errors_test.go` | `project` | Create | Unit assertions for `ProjectError`, sentinels, and predicates |
| 37 | `pkg/parser/errors_test.go` | `parser` | Create | Unit assertions for `ParseError`, sentinels, and predicates |
| 38 | `pkg/diff/errors_test.go` | `diff` | Create | Unit assertions for `DiffError`, sentinels, and predicates |
| 39 | `pkg/git/errors_test.go` | `git` | Create | Unit assertions for `GitError`, sentinels, and predicates |
| 40 | `pkg/cache/errors_test.go` | `cache` | Create | Unit assertions for `CacheError`, sentinels, and predicates |
| 41 | `pkg/lint/errors_test.go` | `lint` | Create | Unit assertions for `LintError`, sentinels, and predicates |
| 42 | `pkg/spec/errors_test.go` | `spec` | Create | Unit assertions for `SpecError`, sentinels, and predicates |
| 43 | `pkg/pipeline/errors_test.go` | `pipeline` | Create | Unit assertions for `PipelineError`, sentinels, and predicates |

---

### 4.2 Error Architecture Specifications

Each `errors.go` must implement:
1. Exported sentinel errors with standard prefix: `var Err<Concept> = errors.New("<pkg>: <detail>")`
2. Uniform typed error struct `<Subsystem>Error`:
   ```go
   type <Subsystem>Error struct {
       Op   string // Operation being executed (e.g. "load", "parse", "show")
       Path string // Filesystem path or target URL
       Key  string // Subsystem identifier (e.g. contract name, rule ID, frame query)
       Err  error  // Underlying cause or sentinel error
   }

   func (e *<Subsystem>Error) Error() string { ... }
   func (e *<Subsystem>Error) Unwrap() error { return e.Err }
   ```
3. High-performance typed error predicates using Go 1.27 `errors.AsType`:
   ```go
   func Is<Predicate>(err error) bool {
       if err == nil {
           return false
       }
       if errors.Is(err, Err<Sentinel>) {
           return true
       }
       if typedErr, ok := errors.AsType[*<Subsystem>Error](err); ok {
           return errors.Is(typedErr.Err, Err<Sentinel>)
       }
       return false
   }
   ```

#### Subsystem Mapping Matrix

| Subsystem | File Path | Sentinel Errors | Struct Type | Predicate Functions |
|---|---|---|---|---|
| `project` | `pkg/project/errors.go` | `ErrWorkspaceNotFound`, `ErrConfigNotFound`, `ErrInvalidConfig`, `ErrContractNotFound`, `ErrStaleCodegen` | `ProjectError` | `IsNotFound(err)`, `IsStale(err)` |
| `parser` | `pkg/parser/errors.go` | `ErrSyntaxError`, `ErrContractNotFound`, `ErrInvalidDirective`, `ErrUnresolvedType` | `ParseError` | `IsSyntaxError(err)`, `IsNotFound(err)` |
| `diff` | `pkg/diff/errors.go` | `ErrStackEmpty`, `ErrFrameNotFound`, `ErrInsufficientFrames`, `ErrConflict` | `DiffError` | `IsNotFound(err)`, `IsConflict(err)` |
| `git` | `pkg/git/errors.go` | `ErrNotRepository`, `ErrBranchNotFound`, `ErrGitCommandFailed` | `GitError` | `IsNotRepository(err)`, `IsNotFound(err)` |
| `cache` | `pkg/cache/errors.go` | `ErrSessionNotFound`, `ErrSecretNotFound`, `ErrCorruptCache` | `CacheError` | `IsNotFound(err)`, `IsCorrupt(err)` |
| `lint` | `pkg/lint/errors.go` | `ErrLintFailure`, `ErrRuleNotFound`, `ErrFixFailed` | `LintError` | `IsLintFailure(err)`, `IsNotFound(err)` |
| `spec` | `pkg/spec/errors.go` | `ErrSpecNotFound`, `ErrUnsupportedFormat`, `ErrEmptySpec` | `SpecError` | `IsNotFound(err)`, `IsUnsupportedFormat(err)` |
| `pipeline` | `pkg/pipeline/errors.go` | `ErrTargetFileRequired`, `ErrNoContractsFound`, `ErrPipelineAborted` | `PipelineError` | `IsPipelineAborted(err)`, `IsNotFound(err)` |

#### Unit Test Specification (`errors_test.go`)
Every `errors_test.go` must execute 5 standard test cases per predicate:
1. **Direct Sentinel**: `require.True(t, Is<Pred>(Err<Sentinel>))`
2. **Wrapped Sentinel**: `require.True(t, Is<Pred>(fmt.Errorf("wrap: %w", Err<Sentinel>)))`
3. **Typed Struct**: `require.True(t, Is<Pred>(&<Subsystem>Error{Op: "op", Err: Err<Sentinel>}))`
4. **Negative Case**: `require.False(t, Is<Pred>(errors.New("unrelated error")))`
5. **Nil Error**: `require.False(t, Is<Pred>(nil))`

---

### 4.3 Godoc Standard Structural Anatomy

Every `doc.go` must follow this uniform structure:
```go
// Copyright (c) 2026 Lemon4ksan All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package <pkg> provides <one-line executive mission statement>.
//
// # Architecture Overview
//
// <High-level technical description explaining how this package fits into the Vortex compiler/toolchain.>
//
//	+-----------------------+       +-----------------------+
//	|     Input Stage       | ----> |     Process Stage     |
//	+-----------------------+       +-----------------------+
//	                                            |
//	                                            v
//	                                +-----------------------+
//	                                |     Output Stage      |
//	                                +-----------------------+
//
// # Core Building Blocks
//
//   - [TypeName] represents ...
//   - [OtherType] provides ...
//
// # Usage Tiers
//
// ## Tier 1: Basic / Quick Start
//
// <Single-call idiomatic entry points and simple examples.>
//
// ## Tier 2: Advanced Configuration
//
// <Custom options, pipeline composition, and complex operational flows.>
//
// ## Tier 3: Low-Level / Internal Optimizations
//
// <Memory allocation profiles, AST node manipulation, affine borrow checking, zero-copy buffers.>
//
// # Concurrency & Thread Safety
//
// <Explicit guarantees: goroutine safety of instances, read-only guarantees after initialization, mutex locking behavior.>
//
// # Performance & Zero-Allocation Profile
//
// <Specific memory allocation guarantees, stack allocation boundaries, L1 cache-line alignment.>
package <pkg>
```

---

### 4.4 Worker Ownership Matrix & Batch Execution Plan

To execute Milestone 3 cleanly without build conflicts, the implementation worker must execute in 6 sequential batches:

```
+------------------------------------------------------------------------------------+
|  Batch 1: Error Architecture & Tests                                              |
|  - 8 errors.go files (project, parser, diff, git, cache, lint, spec, pipeline)     |
|  - 8 errors_test.go files (5 assertions per predicate)                             |
|  -> Verification: go test ./pkg/project ./pkg/parser ./pkg/diff ...               |
+------------------------------------------------------------------------------------+
                                          |
                                          v
+------------------------------------------------------------------------------------+
|  Batch 2: Core Stubs Overhaul & Internal Docs                                      |
|  - Overhaul 4 stubs: pkg/emitter, pkg/ir, pkg/optimizer, pkg/parser              |
|  - Author 2 internal docs: internal/borrow, internal/inspector                    |
|  -> Verification: go test ./pkg/emitter ./pkg/ir ... && go doc ./pkg/emitter      |
+------------------------------------------------------------------------------------+
                                          |
                                          v
+------------------------------------------------------------------------------------+
|  Batch 3: Compilation & Toolchain Docs (7 packages)                               |
|  - pkg/builder, pkg/cfg, pkg/diff, pkg/merge, pkg/patcher, pkg/pipeline,          |
|    pkg/history                                                                     |
|  -> Verification: go test ./pkg/builder ./pkg/pipeline ...                         |
+------------------------------------------------------------------------------------+
                                          |
                                          v
+------------------------------------------------------------------------------------+
|  Batch 4: Specifications & Interop Docs (7 packages)                              |
|  - pkg/enum, pkg/jsbundle, pkg/lint, pkg/oracle/gen, pkg/oracle/spec,             |
|    pkg/project, pkg/spec                                                           |
|  -> Verification: go test ./pkg/lint ./pkg/project ...                             |
+------------------------------------------------------------------------------------+
                                          |
                                          v
+------------------------------------------------------------------------------------+
|  Batch 5: System, Infrastructure & Utility Docs (7 packages)                      |
|  - pkg/cache, pkg/git, pkg/ingest, pkg/mirror, pkg/sys, pkg/tuple, pkg/version    |
|  -> Verification: go test ./pkg/cache ./pkg/git ...                               |
+------------------------------------------------------------------------------------+
                                          |
                                          v
+------------------------------------------------------------------------------------+
|  Batch 6: Workspace-Wide Acceptance & Lint Verification                           |
|  - go test -count=1 ./...                                                         |
|  - golangci-lint run --allow-parallel-runners ./...                                |
|  - godoc verification across all packages                                         |
+------------------------------------------------------------------------------------+
```

#### Detailed Batch Breakdown:

- **Batch 1: Error Architecture (16 files)**:
  - Create `pkg/project/errors.go` and `pkg/project/errors_test.go`
  - Create `pkg/parser/errors.go` and `pkg/parser/errors_test.go`
  - Create `pkg/diff/errors.go` and `pkg/diff/errors_test.go`
  - Create `pkg/git/errors.go` and `pkg/git/errors_test.go`
  - Create `pkg/cache/errors.go` and `pkg/cache/errors_test.go`
  - Create `pkg/lint/errors.go` and `pkg/lint/errors_test.go`
  - Create `pkg/spec/errors.go` and `pkg/spec/errors_test.go`
  - Create `pkg/pipeline/errors.go` and `pkg/pipeline/errors_test.go`
  - *Checkpoint*: Run `go test -v ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline`.

- **Batch 2: Core Stubs & Internal Docs (6 files)**:
  - Overhaul `pkg/emitter/doc.go` (DTO encoders, sub-requesters, harness gen, ASCII pipeline diagram, 3 tiers)
  - Overhaul `pkg/ir/doc.go` (4 scopes: Service, Method, Param, Return, type invariants, 3 tiers)
  - Overhaul `pkg/optimizer/doc.go` (clustering, 64-byte stack sizing, query canonicalization, 3 tiers)
  - Overhaul `pkg/parser/doc.go` (directive parsing, lexer, AST traversal, binder, 3 tiers)
  - Create `internal/borrow/doc.go` (borrow-checker analyzer, Acquire/Release tracking, 3 tiers)
  - Create `internal/inspector/doc.go` (HTTP capture dashboard, WebSocket telemetry, 3 tiers)
  - *Checkpoint*: Verify with `go doc ./pkg/emitter`, `go doc ./pkg/ir`, `go doc ./internal/borrow`.

- **Batch 3: Compilation & Toolchain Docs (7 files)**:
  - Create `pkg/builder/doc.go`
  - Create `pkg/cfg/doc.go`
  - Create `pkg/diff/doc.go`
  - Create `pkg/merge/doc.go`
  - Create `pkg/patcher/doc.go`
  - Create `pkg/pipeline/doc.go`
  - Create `pkg/history/doc.go`
  - *Checkpoint*: Run `go doc ./pkg/pipeline`, `go test ./pkg/pipeline`.

- **Batch 4: Specifications & Interop Docs (7 files)**:
  - Create `pkg/enum/doc.go`
  - Create `pkg/jsbundle/doc.go`
  - Create `pkg/lint/doc.go`
  - Create `pkg/oracle/gen/doc.go` (`package gen`)
  - Create `pkg/oracle/spec/doc.go` (`package spec`)
  - Create `pkg/project/doc.go`
  - Create `pkg/spec/doc.go`
  - *Checkpoint*: Run `go doc ./pkg/lint`, `go test ./pkg/lint`.

- **Batch 5: System, Infrastructure & Utility Docs (7 files)**:
  - Create `pkg/cache/doc.go`
  - Create `pkg/git/doc.go`
  - Create `pkg/ingest/doc.go`
  - Create `pkg/mirror/doc.go`
  - Create `pkg/sys/doc.go`
  - Create `pkg/tuple/doc.go`
  - Create `pkg/version/doc.go`
  - *Checkpoint*: Run `go doc ./pkg/cache`, `go test ./pkg/cache`.

- **Batch 6: Workspace-Wide Verification & Quality Gate**:
  - Run `go test -count=1 ./...`
  - Run `golangci-lint run --allow-parallel-runners ./...`
  - Verify zero test failures, zero lint issues, and clean Godoc output across all 41 packages.

---

## 5. Verification Method

To independently verify the Milestone 3 implementation:

1. **Verify All 43 Target Files Exist**:
   ```powershell
   $files = @(
       "pkg/builder/doc.go", "pkg/cache/doc.go", "pkg/cfg/doc.go", "pkg/diff/doc.go",
       "pkg/enum/doc.go", "pkg/git/doc.go", "pkg/history/doc.go", "pkg/ingest/doc.go",
       "pkg/jsbundle/doc.go", "pkg/lint/doc.go", "pkg/merge/doc.go", "pkg/mirror/doc.go",
       "pkg/oracle/gen/doc.go", "pkg/oracle/spec/doc.go", "pkg/patcher/doc.go",
       "pkg/pipeline/doc.go", "pkg/project/doc.go", "pkg/spec/doc.go", "pkg/sys/doc.go",
       "pkg/tuple/doc.go", "pkg/version/doc.go", "pkg/emitter/doc.go", "pkg/ir/doc.go",
       "pkg/optimizer/doc.go", "pkg/parser/doc.go", "internal/borrow/doc.go",
       "internal/inspector/doc.go",
       "pkg/project/errors.go", "pkg/parser/errors.go", "pkg/diff/errors.go",
       "pkg/git/errors.go", "pkg/cache/errors.go", "pkg/lint/errors.go",
       "pkg/spec/errors.go", "pkg/pipeline/errors.go",
       "pkg/project/errors_test.go", "pkg/parser/errors_test.go", "pkg/diff/errors_test.go",
       "pkg/git/errors_test.go", "pkg/cache/errors_test.go", "pkg/lint/errors_test.go",
       "pkg/spec/errors_test.go", "pkg/pipeline/errors_test.go"
   )
   $missing = $files | Where-Object { -not (Test-Path $_) }
   if ($missing.Count -eq 0) { Write-Output "✔ All 43 files exist" } else { Write-Error "Missing: $($missing -join ', ')" }
   ```

2. **Verify Error Predicates & Test Coverage**:
   ```bash
   go test -v -count=1 ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline
   ```
   *Expected Outcome*: All test cases in `errors_test.go` pass (exit code 0).

3. **Verify Full Workspace Test Suite**:
   ```bash
   go test -count=1 ./...
   ```
   *Expected Outcome*: Exit code 0, all 41 packages pass cleanly.

4. **Verify Zero Linter Violations**:
   ```bash
   golangci-lint run --allow-parallel-runners ./...
   ```
   *Expected Outcome*: Exit code 0, "0 issues."

5. **Verify Clean Godoc Output**:
   ```powershell
   go list ./pkg/... ./internal/... | ForEach-Object { go doc $_ }
   ```
   *Expected Outcome*: Clean documentation output for every package without syntax or formatting errors.
