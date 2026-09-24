# Dispatch: Survey 1 - DTO Codegen & Generic Monads Specification Mining

## Assignment
You are spec_miner_survey_1.
Working directory: d:/CodingProjects/vortex/.agents/spec_miner_survey_1/
Workspace root: d:/CodingProjects/vortex
Authoritative requirements: d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md

## Objective
Investigate and document the current implementation and exact specification requirements for:
1. R1: Zero-Allocation generic.Optional[T] & Empty Field DTO Codegen
   - Inspect pkg/emitter/dto.go, IR type resolution in pkg/ir/ (or relevant packages), foundation/generic/monads.go.
   - Check how generic.Optional[T] is currently defined in foundation/generic/monads.go, whether MarshalJSON / UnmarshalJSON exist or what signature/behavior is needed.
   - Check how DTO code is emitted in pkg/emitter/dto.go (look for fmt.Sprint calls, reflection, string concatenations).
   - Check how query, form-data, and url.Values serialization is handled.
   - Check handling of generic.Some("") vs generic.None().
2. R4: Performance Benchmarks & Zero-Alloc Guarantees
   - Identify existing benchmarks or test suites for DTO generation.
   - Determine how zero-alloc tests (b.ReportAllocs()) can be constructed and verified.

## Output Requirements
Produce a comprehensive handoff report at:
d:/CodingProjects/vortex/.agents/spec_miner_survey_1/handoff.md
Detailing existing code paths, exact functions/types to modify, interfaces, and recommendations for milestone execution.

## 2026-09-22T14:22:08Z
Task:
Mine and document the specifications and code paths for:
1. R1: Zero-Allocation generic.Optional[T] & Empty Field DTO Codegen (pkg/emitter/dto.go, IR type resolution, foundation/generic/monads.go).
2. R4: Performance Benchmarks & Zero-Alloc Guarantees (b.ReportAllocs(), benchmark tests).

Investigate:
- foundation/generic/monads.go definition of Optional[T], and requirements for MarshalJSON and UnmarshalJSON.
- pkg/emitter/dto.go: look for fmt.Sprint reflection anti-patterns, heap allocations, query/form/url.Values serialization.
- IR type resolution in the codebase.
- Serialization behavior of generic.Some("") (explicit empty string -> key=) vs generic.None() (omitted).
- Existing test suites and benchmarks for DTO codegen.

Write your comprehensive findings and handoff report to:
d:/CodingProjects/vortex/.agents/spec_miner_survey_1/handoff.md

When complete, use send_message to report your completion and the path to your handoff report.

