# Dispatch: explorer_m3_2

- Identity: explorer_m3_2 (teamwork_preview_explorer)
- Working directory: d:/CodingProjects/vortex/.agents/explorer_m3_2/
- Milestone: Milestone 3 — Benchmark-Grade Code Documentation & Architecture

## Context & Inputs
You MUST read:
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md`
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md`
3. `d:/CodingProjects/vortex/.agents/survey_arch_bench_1/handoff.md`

## Objectives
Explore and formulate the standardized Sentinel Errors and Typed Error Predicates architecture:
1. Standard design pattern per Aoni / Foundation:
   - Exported package sentinel errors: `var Err... = errors.New("<pkg>: <concept>")`
   - Typed error struct: `type <Subsystem>Error struct { Op, Path, Key string; Err error }` with `Error() string` and `Unwrap() error`.
   - Typed error predicates: `func Is<Predicate>(err error) bool` checking `errors.Is` and `errors.AsType[*<Subsystem>Error]`.
2. Map out concrete `errors.go` implementations across key subsystems:
   - `pkg/project/errors.go` (`ErrWorkspaceNotFound`, `ErrConfigNotFound`, `ErrInvalidConfig`, `ErrContractNotFound`, `ErrStaleCodegen`, `IsNotFound`, `IsStale`)
   - `pkg/parser/errors.go` (`ErrSyntaxError`, `ErrContractNotFound`, `ErrInvalidDirective`, `ErrUnresolvedType`, `IsSyntaxError`, `IsNotFound`)
   - `pkg/diff/errors.go` (`ErrStackEmpty`, `ErrFrameNotFound`, `ErrInsufficientFrames`, `ErrConflict`, `IsNotFound`, `IsConflict`)
   - `pkg/git/errors.go` (`ErrNotRepository`, `ErrBranchNotFound`, `ErrGitCommandFailed`, `IsNotRepository`)
   - `pkg/cache/errors.go` (`ErrSessionNotFound`, `ErrSecretNotFound`, `ErrCorruptCache`, `IsNotFound`)
   - `pkg/lint/errors.go` (`ErrLintFailure`, `ErrRuleNotFound`, `ErrFixFailed`, `IsLintFailure`)
   - `pkg/spec/errors.go` (`ErrSpecNotFound`, `ErrUnsupportedFormat`, `ErrEmptySpec`, `IsNotFound`)
   - `pkg/pipeline/errors.go` (`ErrTargetFileRequired`, `ErrNoContractsFound`, `ErrPipelineAborted`, `IsPipelineAborted`)
3. Ensure backwards compatibility: existing error messages and tests must remain intact.

Formulate the error declarations and types in `d:/CodingProjects/vortex/.agents/explorer_m3_2/handoff.md`.
Do NOT modify production files (Explorer is read-only).
Send notification message back to parent when done.

## 2026-09-22T20:04:21Z
You are explorer_m3_2.
Your working directory is d:/CodingProjects/vortex/.agents/explorer_m3_2/.
You MUST read:
1. d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md
2. d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md
3. d:/CodingProjects/vortex/.agents/survey_arch_bench_1/handoff.md
4. d:/CodingProjects/vortex/.agents/explorer_m3_2/DISPATCH.md

Task:
Explore and design the standardized Sentinel Errors and Typed Error Predicates architecture:
- Standard design: exported sentinels `var Err... = errors.New(...)`, typed structs `type <Subsystem>Error struct { Op, Path, Key string; Err error }` with `Error() string` and `Unwrap() error`, typed predicates `func Is<Predicate>(err error) bool` using Go 1.27 `errors.AsType`.
- Map out concrete errors.go for key subsystems: project, parser, diff, git, cache, lint, spec, pipeline.
- Ensure backwards compatibility with existing error strings and tests.
Write your report in d:/CodingProjects/vortex/.agents/explorer_m3_2/handoff.md.
Do NOT modify production files (Explorer is read-only).
Send notification message back to parent when done.
