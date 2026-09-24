# Dispatch: reviewer_m3_2 (Error Architecture & Conformance Reviewer)

- Target: Milestone 3 — Benchmark-Grade Code Documentation & Architecture
- Working directory: `d:/CodingProjects/vortex/.agents/reviewer_m3_2/`
- Workspace root: `d:/CodingProjects/vortex`

## Mandatory Inputs (Read Before Starting)
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md` — Authoritative user requirements
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md` — Project architecture & contracts
3. `d:/CodingProjects/vortex/.agents/worker_m3/handoff.md` — Implementation report (43 files created/updated)

## Review Task
1. Inspect the 8 `errors.go` implementations and companion `errors_test.go` suites:
   - `pkg/project/errors.go`
   - `pkg/parser/errors.go`
   - `pkg/diff/errors.go`
   - `pkg/git/errors.go`
   - `pkg/cache/errors.go`
   - `pkg/lint/errors.go`
   - `pkg/spec/errors.go`
   - `pkg/pipeline/errors.go`
2. Verify architectural symmetry:
   - Exported package sentinel errors: `var Err... = errors.New("<pkg>: ...")`
   - Uniform typed struct `<Subsystem>Error` with fields `Op`, `Path`/`Key`, `Err`
   - `Error() string` and `Unwrap() error` methods implemented
   - Go 1.27 generic `errors.AsType[*<Subsystem>Error](err)` used in predicate functions
   - 5 standard assertions verified in each companion `errors_test.go`
3. Verify backwards compatibility: no breaking changes to preexisting package signatures or behavior.
4. Run full test suite: `$env:GOWORK="off"; go test -count=1 ./...`
5. Run linter: `golangci-lint run --allow-parallel-runners ./...`
6. Write your comprehensive review report with an explicit verdict (**APPROVE** or **REQUEST_CHANGES**) in `d:/CodingProjects/vortex/.agents/reviewer_m3_2/handoff.md`.
7. Send a message to parent notifying that your handoff is ready.

## 2026-09-23T04:28:22Z

You are reviewer_m3_2 (teamwork_preview_reviewer).
Your working directory is d:/CodingProjects/vortex/.agents/reviewer_m3_2/.
Workspace root: d:/CodingProjects/vortex.

You MUST read:
1. d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md
2. d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md
3. d:/CodingProjects/vortex/.agents/worker_m3/handoff.md
4. d:/CodingProjects/vortex/.agents/reviewer_m3_2/DISPATCH.md

Task:
Perform independent review of Milestone 3 Error Architecture & Conformance:
1. Inspect the 8 errors.go and 8 errors_test.go files (pkg/project, pkg/parser, pkg/diff, pkg/git, pkg/cache, pkg/lint, pkg/spec, pkg/pipeline). Verify sentinels, uniform <Subsystem>Error struct, Unwrap, Go 1.27 errors.AsType predicates, and 5 standard assertions per test suite.
2. Run tests: $env:GOWORK="off"; go test -count=1 ./...
3. Run linter: golangci-lint run --allow-parallel-runners ./...
4. Write your review report and final verdict (APPROVE or REQUEST_CHANGES) in d:/CodingProjects/vortex/.agents/reviewer_m3_2/handoff.md.
5. When complete, send a message to parent notifying that your handoff is ready.

