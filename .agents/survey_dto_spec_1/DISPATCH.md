## 2026-09-22T17:26:47Z

Scope of Investigation (R1: Zero-Allocation generic.Optional[T] & Empty Field DTO Codegen):
1. Locate and inspect pkg/emitter/dto.go and IR type resolution logic across the codebase.
2. Locate where generic.Optional[T] is defined (check foundation/generic/monads.go, go.mod dependencies, or local packages).
3. Investigate MarshalJSON and UnmarshalJSON implementations for generic.Optional[T].
4. Analyze current DTO emission for query, form-data, and JSON: where heap allocations (like fmt.Sprint) occur, how primitive types (string, int, uint, float, bool, time.Time) are currently encoded vs how they should be encoded with zero allocations (strconv.Append*, direct slice/buffer writing).
5. Investigate how empty string generic.Some("") is serialized (must be `key=`) vs unset generic.None() (must be omitted entirely).
6. Investigate url.Values, form-data, and query parameter interoperability with generic.Optional[T].

Deliverables:
Produce a comprehensive handoff report at:
d:/CodingProjects/vortex/.agents/survey_dto_spec_1/handoff.md
Include:
- Exact file paths and line numbers of relevant code
- Detailed technical specification for zero-alloc DTO codegen & generic.Optional[T]
- Edge cases (empty string, 0, false, null, None vs Some(""), time.Time zero value/formatting)
- Enumerated list of specific features required for the Feature Inventory
- Clear recommendations for implementation
