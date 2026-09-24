# Dispatch: explorer_m3_fix_2 (Godoc API Alignment & Sibling Comment Explorer)

- Target: Milestone 3 — Godoc API Alignment & Comment Decontamination
- Working directory: `d:/CodingProjects/vortex/.agents/explorer_m3_fix_2/`
- Workspace root: `d:/CodingProjects/vortex`

## Mandatory Inputs (Read Before Starting)
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md`
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md`
3. `d:/CodingProjects/vortex/.agents/orchestrator_3/GATE_STATUS.md`
4. `d:/CodingProjects/vortex/.agents/challenger_m3_2/handoff.md` (Full Godoc integrity report)
5. `d:/CodingProjects/vortex/.agents/reviewer_m3_1/handoff.md`

## Task
1. Inspect the 4 sibling files with duplicate package headers:
   - `pkg/emitter/emitter.go:5-6`
   - `pkg/ingest/namer.go:5-6`
   - `pkg/lint/rule.go:5-6`
   - `pkg/openapi/importer.go:5-7`
   Specify the exact lines to remove so that only `doc.go` provides package documentation.
2. Inspect the 8 `doc.go` files with hallucinated APIs / broken code examples:
   - `pkg/ingest/doc.go` (replace `ParseHAR` / `ConvertHARToOpenAPI` with `HARToOpenAPI` / `HARToOpenAPIOpts`)
   - `pkg/cache/doc.go` (replace `LoadSecretsVault` / `NewTrafficStore` with actual exported APIs in `pkg/cache`)
   - `pkg/cfg/doc.go` (replace `cfg.Build` with `cfg.New`)
   - `pkg/diff/doc.go` (replace `CheckpointStack` with `DiffStack`)
   - `pkg/jsbundle/doc.go` (replace `ScanDirectory` with `ScanFiles`/`ScanBytes`)
   - `pkg/git/doc.go` (replace `ListBranches` / `IsCleanWorkingTree` with `ListProposalBranches` / `IsClean`)
   - `pkg/mirror/doc.go` (replace `CheckAllServices` with `CheckService`)
   - `pkg/parser/doc.go` (fix brackets around `[T]` in diagram to avoid broken Godoc links, align types)
3. For each of the 8 packages, inspect the actual exported Go files in that package directory to retrieve the true type names and function signatures. Provide exact, verified drop-in replacement snippets.
4. Write your detailed remediation report in `d:/CodingProjects/vortex/.agents/explorer_m3_fix_2/handoff.md`.
5. Send a completion message to parent when done. Do NOT edit code directly (Explorer is read-only).

## 2026-09-23T04:41:30Z
You are explorer_m3_fix_2 (teamwork_preview_explorer).
Your working directory is d:/CodingProjects/vortex/.agents/explorer_m3_fix_2/.
Workspace root: d:/CodingProjects/vortex.

You MUST read:
1. d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md
2. d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md
3. d:/CodingProjects/vortex/.agents/orchestrator_3/GATE_STATUS.md
4. d:/CodingProjects/vortex/.agents/challenger_m3_2/handoff.md
5. d:/CodingProjects/vortex/.agents/reviewer_m3_1/handoff.md
6. d:/CodingProjects/vortex/.agents/explorer_m3_fix_2/DISPATCH.md

Task:
Inspect the 4 sibling files with duplicate package comments (emitter.go, namer.go, rule.go, importer.go) and the 8 doc.go files with hallucinated APIs (ingest, cache, cfg, diff, jsbundle, git, mirror, parser). Inspect actual exported Go declarations in each package and provide exact drop-in replacements for the doc.go usage examples and comments.
Write your report in d:/CodingProjects/vortex/.agents/explorer_m3_fix_2/handoff.md.
Do NOT edit production code directly (Explorer is read-only).
Send a message to parent when done.
