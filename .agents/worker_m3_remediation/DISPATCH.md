# Dispatch: worker_m3_remediation (Milestone 3 Remediation Worker)

- Target: Milestone 3 Remediation (Typed Nil Safety, Sibling Comment Cleanup, Godoc API Alignment)
- Working directory: `d:/CodingProjects/vortex/.agents/worker_m3_remediation/`
- Workspace root: `d:/CodingProjects/vortex`

## MANDATORY INTEGRITY WARNING
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_preview_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

## Mandatory Inputs (Read Before Starting)
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md` — Authoritative user requirements
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md` — Project architecture & contracts
3. `d:/CodingProjects/vortex/.agents/orchestrator_3/GATE_STATUS.md` — Gate verdicts & failure reasons
4. `d:/CodingProjects/vortex/.agents/explorer_m3_fix_1/handoff.md` & `typed_nil_remediation.patch`
5. `d:/CodingProjects/vortex/.agents/explorer_m3_fix_3/handoff.md` — Master 28-file remediation blueprint

## Scope of Work (28 Files Across 4 Batches)

Execute the 4 batches as specified in `explorer_m3_fix_3/handoff.md`:

### Batch 1: Sibling Comment Cleanups (4 Files)
Remove redundant `// Package <name>` comments preceding `package <name>`:
1. `pkg/emitter/emitter.go` (remove line 5)
2. `pkg/ingest/namer.go` (remove line 5)
3. `pkg/lint/rule.go` (remove line 5)
4. `pkg/openapi/importer.go` (remove lines 5-6)

### Batch 2: Error Predicate Typed Nil Guards (8 Files, 17 Predicates)
Add `&& <target> != nil` guards before accessing `.Err` across:
1. `pkg/project/errors.go` (`IsNotFound`, `IsStale`)
2. `pkg/parser/errors.go` (`IsSyntaxError`, `IsNotFound`)
3. `pkg/diff/errors.go` (`IsNotFound`, `IsConflict`, `IsEmpty`)
4. `pkg/git/errors.go` (`IsNotRepository`, `IsNotFound`)
5. `pkg/cache/errors.go` (`IsNotFound`, `IsCorrupt`)
6. `pkg/lint/errors.go` (`IsLintFailure`, `IsNotFound`)
7. `pkg/spec/errors.go` (`IsNotFound`, `IsUnsupportedFormat`)
8. `pkg/pipeline/errors.go` (`IsPipelineAborted`, `IsNotFound`)

### Batch 3: Companion Test Suite Augmentation (8 Files)
Add typed nil pointer tests (`var typedNil *<Subsystem>Error = nil; require.False(...)`) and wrapped typed nil tests (`fmt.Errorf("%w", typedNil)`) to each test function in:
1. `pkg/project/errors_test.go`
2. `pkg/parser/errors_test.go`
3. `pkg/diff/errors_test.go`
4. `pkg/git/errors_test.go`
5. `pkg/cache/errors_test.go`
6. `pkg/lint/errors_test.go`
7. `pkg/spec/errors_test.go`
8. `pkg/pipeline/errors_test.go`

*Checkpoint*: Run `go test -v -count=1 ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline`.

### Batch 4: Godoc Architectural Alignment (8 Files)
Apply the exact drop-in replacements from Section 4 of `explorer_m3_fix_3/handoff.md` to align bracketed links and usage tier code examples with real exported package signatures:
1. `pkg/ingest/doc.go` (`HARToOpenAPI`, `HARToOpenAPIOpts`)
2. `pkg/cache/doc.go` (`LoadSecrets`, `TrafficIndex`, `StoreTraffic`, `GetTraffic`, `IsFresh`/`Put`)
3. `pkg/cfg/doc.go` (`New`, `WalkPaths`, `FindLoopBlocks`)
4. `pkg/diff/doc.go` (`DiffStack`, `LoadStack`, `Compare`)
5. `pkg/jsbundle/doc.go` (`ScanFiles`, `ScanBytes`)
6. `pkg/git/doc.go` (`ListProposalBranches`, `IsClean`)
7. `pkg/mirror/doc.go` (`CheckService`)
8. `pkg/parser/doc.go` (remove `[Lexer]`/`[Token]`, unbracket `generic.Optional[T]`)

### Verification & Acceptance
1. Verify with `go test -v -count=1 ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline`
2. Verify full workspace: `$env:GOWORK="off"; go test -count=1 ./...` (must pass 100% across all 41 packages)
3. Verify linter: `golangci-lint run --allow-parallel-runners ./...` (must report 0 issues)
4. Verify `go doc` renders cleanly without duplicate summary comments on `pkg/emitter`, `pkg/ingest`, `pkg/lint`, and `pkg/openapi`.
5. Write your complete handoff report to `d:/CodingProjects/vortex/.agents/worker_m3_remediation/handoff.md`.
6. Send a message to parent notifying that your handoff is ready.

## 2026-09-23T04:47:41Z
**Context**: Additional Reference Material for Batch 4
**Content**: `explorer_m3_fix_2` has also delivered its handoff report at `d:/CodingProjects/vortex/.agents/explorer_m3_fix_2/handoff.md` along with ready-to-use proposed `doc.go` files in `d:/CodingProjects/vortex/.agents/explorer_m3_fix_2/proposed_*_doc.go`. You may reference these alongside `explorer_m3_fix_3/handoff.md`.
**Action**: Continue execution of the 4 remediation batches.
