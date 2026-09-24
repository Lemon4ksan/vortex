# Dispatch: worker_m3 (Milestone 3 Implementation Worker)

- Target: Milestone 3 — Benchmark-Grade Code Documentation & Architecture
- Working directory: `d:/CodingProjects/vortex/.agents/worker_m3/`
- Workspace root: `d:/CodingProjects/vortex`

## MANDATORY INTEGRITY WARNING
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_preview_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

## Mandatory Inputs (Read Before Starting)
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md` — Authoritative user requirements
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md` — Project architecture & contracts
3. `d:/CodingProjects/vortex/.agents/explorer_m3_1/handoff.md` — Complete Godoc blueprints for 27 `pkg/` + 2 `internal/` packages
4. `d:/CodingProjects/vortex/.agents/explorer_m3_2/handoff.md` — Verbatim Go source code for 8 `errors.go` and 8 `errors_test.go`
5. `d:/CodingProjects/vortex/.agents/explorer_m3_3/handoff.md` — Master 43-file inventory and 6-batch execution plan

## Scope of Work (43 Files Total)

Execute the 6 sequential batches described in `explorer_m3_3/handoff.md`:

### Batch 1: Error Architecture & Tests (16 files)
Author the following 8 `errors.go` and 8 `errors_test.go` files using the verbatim source code in `explorer_m3_2/handoff.md`:
1. `pkg/project/errors.go` & `pkg/project/errors_test.go`
2. `pkg/parser/errors.go` & `pkg/parser/errors_test.go`
3. `pkg/diff/errors.go` & `pkg/diff/errors_test.go`
4. `pkg/git/errors.go` & `pkg/git/errors_test.go`
5. `pkg/cache/errors.go` & `pkg/cache/errors_test.go`
6. `pkg/lint/errors.go` & `pkg/lint/errors_test.go`
7. `pkg/spec/errors.go` & `pkg/spec/errors_test.go`
8. `pkg/pipeline/errors.go` & `pkg/pipeline/errors_test.go`

Each `errors.go` must implement:
- Package sentinel errors (`var Err... = errors.New("<pkg>: ...")`)
- Uniform typed struct `<Subsystem>Error` (`Op`, `Path`, `Key`, `Err`) with `Error() string` and `Unwrap() error`
- Go 1.27 high-performance error predicates using `errors.AsType` (e.g. `IsNotFound(err)`, `IsSyntaxError(err)`, etc.)
- Each `errors_test.go` must verify the 5 standard assertions per predicate (direct sentinel, wrapped sentinel, typed struct, negative unrelated, nil).

*Verification Checkpoint*: Run `go test -v -count=1 ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline`.

### Batch 2: Core Stubs Overhaul & Internal Docs (6 files)
Author/expand using the blueprints from `explorer_m3_1/handoff.md`:
- `pkg/emitter/doc.go` (expand stub: DTO encoders, sub-requesters, harness gen, ASCII pipeline diagram, 3 tiers)
- `pkg/ir/doc.go` (expand stub: 4 scopes: Service, Method, Param, Return, type invariants, 3 tiers)
- `pkg/optimizer/doc.go` (expand stub: clustering, 64-byte stack sizing, query canonicalization, 3 tiers)
- `pkg/parser/doc.go` (expand stub: directive parsing, lexer, AST traversal, binder, 3 tiers)
- `internal/borrow/doc.go` (create: borrow-checker analyzer, Acquire/Release tracking, 3 tiers)
- `internal/inspector/doc.go` (create: HTTP capture dashboard, WebSocket telemetry, 3 tiers)

*Verification Checkpoint*: Verify with `go doc ./pkg/emitter`, `go doc ./pkg/ir`, `go doc ./internal/borrow`.

### Batch 3: Compilation & Toolchain Docs (7 files)
Author using blueprints from `explorer_m3_1/handoff.md`:
- `pkg/builder/doc.go`
- `pkg/cfg/doc.go`
- `pkg/diff/doc.go`
- `pkg/merge/doc.go`
- `pkg/patcher/doc.go`
- `pkg/pipeline/doc.go`
- `pkg/history/doc.go`

### Batch 4: Specifications & Interop Docs (7 files)
Author using blueprints from `explorer_m3_1/handoff.md`:
- `pkg/enum/doc.go`
- `pkg/jsbundle/doc.go`
- `pkg/lint/doc.go`
- `pkg/oracle/gen/doc.go` (`package gen`)
- `pkg/oracle/spec/doc.go` (`package spec`)
- `pkg/project/doc.go`
- `pkg/spec/doc.go`

### Batch 5: System, Infrastructure & Utility Docs (7 files)
Author using blueprints from `explorer_m3_1/handoff.md`:
- `pkg/cache/doc.go`
- `pkg/git/doc.go`
- `pkg/ingest/doc.go`
- `pkg/mirror/doc.go`
- `pkg/sys/doc.go`
- `pkg/tuple/doc.go`
- `pkg/version/doc.go`

### Batch 6: Full Workspace Verification
- Run `$env:GOWORK="off"; go test -count=1 ./...` — must PASS across all 41 packages.
- Run `golangci-lint run --allow-parallel-runners ./...` — must report 0 issues.
- Verify `go doc` renders cleanly for updated packages.

## Completion Criteria
1. All 43 target files exist and compile cleanly.
2. All 8 `errors_test.go` test suites pass 100%.
3. Full workspace `go test -count=1 ./...` passes (exit code 0).
4. Full workspace `golangci-lint run --allow-parallel-runners ./...` passes with 0 issues.
5. Write your complete handoff report to `d:/CodingProjects/vortex/.agents/worker_m3/handoff.md`.
6. Send a message to parent notifying that your handoff is ready.

## 2026-09-22T20:11:01Z

You are worker_m3 (teamwork_preview_worker).
Your working directory is d:/CodingProjects/vortex/.agents/worker_m3/.
Workspace root: d:/CodingProjects/vortex.

MANDATORY INTEGRITY WARNING:
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_preview_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

You MUST read:
1. d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md
2. d:/CodingProjects/vortex/.agents/worker_m3/DISPATCH.md
3. d:/CodingProjects/vortex/.agents/explorer_m3_1/handoff.md
4. d:/CodingProjects/vortex/.agents/explorer_m3_2/handoff.md
5. d:/CodingProjects/vortex/.agents/explorer_m3_3/handoff.md

Your task is to implement Milestone 3 (Benchmark-Grade Code Documentation & Architecture):
Author all 43 files in 6 sequential batches as specified in DISPATCH.md and explorer_m3_3/handoff.md:
- Batch 1: 8 errors.go and 8 errors_test.go files using verbatim source from explorer_m3_2/handoff.md.
- Batch 2: Overhaul 4 stubs (emitter, ir, optimizer, parser) and author 2 internal doc.go (borrow, inspector) using blueprints from explorer_m3_1/handoff.md.
- Batch 3: Author 7 compilation doc.go files (builder, cfg, diff, merge, patcher, pipeline, history).
- Batch 4: Author 7 spec & interop doc.go files (enum, jsbundle, lint, oracle/gen, oracle/spec, project, spec).
- Batch 5: Author 7 system & infra doc.go files (cache, git, ingest, mirror, sys, tuple, version).
- Batch 6: Run full workspace test and lint verification:
  $env:GOWORK="off"; go test -count=1 ./...
  golangci-lint run --allow-parallel-runners ./...
  go doc verification.

Document all created/modified files, test commands, and lint outputs in d:/CodingProjects/vortex/.agents/worker_m3/handoff.md.
When complete, send a message to parent notifying that your handoff is ready.

## 2026-09-23T04:22:42Z

**Context**: Milestone 3 Execution Resume after Server Restart

**Content**: A server restart occurred. Your progress log shows that Batches 1–4 are fully complete and passing. Please resume execution and finish the remaining work:
1. **Batch 5: System, Infrastructure & Utility Docs (7 files)**:
   Author `doc.go` for:
   - `pkg/cache/doc.go`
   - `pkg/git/doc.go`
   - `pkg/ingest/doc.go`
   - `pkg/mirror/doc.go`
   - `pkg/sys/doc.go`
   - `pkg/tuple/doc.go`
   - `pkg/version/doc.go`
   using the exact structural blueprints, ASCII diagrams, and 3-tier usage sections from `d:/CodingProjects/vortex/.agents/explorer_m3_1/handoff.md`.
2. **Batch 6: Full Workspace Verification**:
   - Run `$env:GOWORK="off"; go test -count=1 ./...` — verify all 41 packages pass cleanly (exit code 0).
   - Run `golangci-lint run --allow-parallel-runners ./...` — verify 0 issues reported (exit code 0).
   - Validate `go doc` renders cleanly for all modified packages.
3. **Handoff Report**:
   - Document all 43 created/modified files, test commands, and linter results in `d:/CodingProjects/vortex/.agents/worker_m3/handoff.md`.

**Action**: Complete Batches 5 & 6, author `handoff.md`, and send a completion message back to parent when ready.

