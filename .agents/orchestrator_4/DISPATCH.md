# Dispatch: orchestrator_4

- Identity: orchestrator_4 (Project Orchestrator, successor gen4)
- Working directory: d:/CodingProjects/vortex/.agents/orchestrator_4/
- Parent: parent (`eb195458-ce5e-4c9b-a4b3-12d8748b49f4`)

## Mission
You are the successor Project Orchestrator (orchestrator_4) continuing the sovereign upgrade of Vortex.
Milestone 1 (Zero-Alloc DTO emission & generic.Optional[T]) is PASSED and verified.
Milestone 2 (Restrained High-Craft CLI Presentation via tuikit) is PASSED and verified.

## Next Steps
1. Initialize your BRIEFING.md, DISPATCH.md, and progress.md in `d:/CodingProjects/vortex/.agents/orchestrator_4/`.
2. Read:
   - `d:/CodingProjects/vortex/.agents/orchestrator_3/handoff.md`
   - `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md`
   - `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md`
3. Execute Milestone 3: Benchmark-Grade Code Documentation & Architecture
   - Comprehensive doc.go files across all packages in pkg/ following aoni/foundation standards.
   - Expand stubs in `pkg/emitter`, `pkg/ir`, `pkg/optimizer`, `pkg/parser`.
   - Add `doc.go` for internal packages (`internal/borrow`, `internal/inspector`).
   - Standardized sentinel errors and typed error predicates across key packages.
   - Run Explorer -> Worker -> Reviewers -> Challengers -> Auditor gate cycle.
4. Execute Milestone 4: Performance Benchmarks & Adversarial Test Coverage
   - Zero-allocation regression benchmarks (`b.ReportAllocs()`) for all generated DTO emitters (`AppendQuery`, `AppendFormData`, `EncodeValues`).
   - Comprehensive `generic.Optional[T]` tests.
5. Final Workspace Acceptance Gate: verify `go test ./...` and `golangci-lint run` pass 100% cleanly with 0 issues.
6. Report final completion with structured summary when all criteria are met.
