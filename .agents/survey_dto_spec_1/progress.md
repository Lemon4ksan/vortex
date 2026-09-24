# Progress — survey_dto_spec_1

Last visited: 2026-09-22T17:34:00Z

- [x] Read ORIGINAL_REQUEST.md and DISPATCH.md
- [x] Initialized metadata (BRIEFING.md, DISPATCH.md, progress.md)
- [x] 1. Locate and inspect pkg/emitter/dto.go and IR type resolution logic across the codebase
- [x] 2. Locate where generic.Optional[T] is defined (check foundation/generic/monads.go, go.mod dependencies, or local packages)
- [x] 3. Investigate MarshalJSON and UnmarshalJSON implementations for generic.Optional[T]
- [x] 4. Analyze current DTO emission for query, form-data, and JSON: heap allocations (fmt.Sprint) vs zero-alloc encoders (strconv.Append*, direct slice/buffer writing)
- [x] 5. Investigate how empty string generic.Some("") is serialized (must be `key=`) vs unset generic.None() (omitted entirely)
- [x] 6. Investigate url.Values, form-data, and query parameter interoperability with generic.Optional[T]
- [x] 7. Synthesize findings, edge cases, feature inventory, and write handoff.md
- [x] 8. Send message to parent
