# DISPATCH — worker_m3_quickfix

## 2026-09-23T05:02:19Z

**Task**: Final 2 Line-Level Doc Edits for Milestone 3 Gate Passage
**Context**: Challenger challenger_m3_iter2_2 identified 2 minor doc comment defects in pkg/parser/doc.go and pkg/diff/doc.go. All other 4 gate agents (reviewer_m3_iter2_1, reviewer_m3_iter2_2, challenger_m3_iter2_1, auditor_m3_iter2_1) have APPROVED / verified CLEAN.

## Mandatory Reading
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md`
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md`
3. `d:/CodingProjects/vortex/.agents/challenger_m3_iter2_2/handoff.md`

## Instructions
1. **Fix `pkg/parser/doc.go` (around line 47)**:
   - Remove the unresolvable link line:
     `//   - [ParseDirectives]: Scans AST comment groups and extracts all declared directives.`
   - Keep `[ParseDirective]` and `[Directive]`.

2. **Fix `pkg/diff/doc.go` (around lines 64-67)**:
   - Change:
     ```go
     // Apply customized comparison options such as ignoring deprecated endpoints or custom tag filtering:
     //
     //	opts := diff.DiffOptions{IgnoreDeprecated: true}
     //	report := diff.CompareWithOptions(localRootIR, remoteDoc, "pkg/api", "openapi.yaml", opts)
     ```
     to:
     ```go
     // Apply customized comparison options such as additive mode to suppress ghost endpoint noise:
     //
     //	opts := diff.DiffOptions{Additive: true}
     //	report := diff.CompareWithOptions(localRootIR, remoteDoc, "pkg/api", "openapi.yaml", opts)
     ```

3. **Mandatory Integrity Warning**:
   DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_preview_auditor will independently verify your work.

4. **Verification**:
   - `$env:GOWORK="off"; go test -count=1 ./...`
   - `golangci-lint run --allow-parallel-runners ./...`
   - Verify `go doc ./pkg/parser ParseDirective` succeeds.

5. **Deliverables**:
   - Write handoff report to `d:/CodingProjects/vortex/.agents/worker_m3_quickfix/handoff.md`.
   - Send notification message to parent.
