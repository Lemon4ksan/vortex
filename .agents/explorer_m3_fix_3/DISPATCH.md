# Dispatch: explorer_m3_fix_3 (Milestone 3 Remediation Synthesis & Plan)

- Target: Milestone 3 — Comprehensive Remediation Specification
- Working directory: `d:/CodingProjects/vortex/.agents/explorer_m3_fix_3/`
- Workspace root: `d:/CodingProjects/vortex`

## Mandatory Inputs (Read Before Starting)
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md`
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md`
3. `d:/CodingProjects/vortex/.agents/orchestrator_3/GATE_STATUS.md`
4. `d:/CodingProjects/vortex/.agents/challenger_m3_1/handoff.md` (Typed nil panic traces)
5. `d:/CodingProjects/vortex/.agents/challenger_m3_2/handoff.md` (Godoc rendering & API hallucinations)
6. `d:/CodingProjects/vortex/.agents/reviewer_m3_1/handoff.md`

## Task
Synthesize a complete, unified remediation blueprint for the implementation worker:
1. Synthesize the complete file modification inventory:
   - 8 `errors.go` files to patch with typed nil checks (`&& pErr != nil`)
   - 8 `errors_test.go` suites to augment with typed nil tests
   - 4 sibling files to clean of duplicate comments (`emitter.go`, `namer.go`, `rule.go`, `importer.go`)
   - 8 `doc.go` files to align with genuine package APIs
2. Provide exact verification commands:
   - `go test -v -count=1 ./pkg/project ./pkg/parser ./pkg/diff ./pkg/git ./pkg/cache ./pkg/lint ./pkg/spec ./pkg/pipeline`
   - `$env:GOWORK="off"; go test -count=1 ./...`
   - `golangci-lint run --allow-parallel-runners ./...`
   - `go doc` inspection on all 8 updated `doc.go` files and 4 cleaned packages.
3. Write your synthesized report in `d:/CodingProjects/vortex/.agents/explorer_m3_fix_3/handoff.md`.
4. Send a completion message to parent when done. Do NOT edit code directly (Explorer is read-only).

## 2026-09-23T04:41:30Z
User request received:
Synthesize a comprehensive, master remediation specification for the implementation worker:
- Inventory of all files to modify (8 errors.go, 8 errors_test.go, 4 sibling files, 8 doc.go).
- Exact changes, line targets, and verified Go snippets.
- Execution batches and verification commands.
Write your report in d:/CodingProjects/vortex/.agents/explorer_m3_fix_3/handoff.md.
Do NOT edit production code directly (Explorer is read-only).
Send a message to parent when done.
