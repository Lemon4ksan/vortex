## 2026-09-22T14:43:12Z

You are worker_m1, implementing Milestone 1 (Zero-Allocation generic.Optional[T] & Empty Field DTO Codegen).

MANDATORY: You MUST read d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md and d:/CodingProjects/vortex/.agents/orchestrator_2/PROJECT.md before starting work. Do NOT proceed without reading them.

Also read the 3 explorer reports carefully:
- d:/CodingProjects/vortex/.agents/explorer_m1_1/handoff.md
- d:/CodingProjects/vortex/.agents/explorer_m1_2/handoff.md
- d:/CodingProjects/vortex/.agents/spec_miner_m1_3/handoff.md

Your working directory is: d:/CodingProjects/vortex/.agents/worker_m1/
Please create your metadata files (BRIEFING.md, progress.md) in your working directory.

MANDATORY INTEGRITY WARNING:
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_preview_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

Exclusive Write Ownership:
You own exclusively the following files:
- d:/CodingProjects/vortex/pkg/parser/binder.go
- d:/CodingProjects/vortex/pkg/parser/parser_test.go
- d:/CodingProjects/vortex/pkg/emitter/dto.go
- d:/CodingProjects/vortex/pkg/emitter/dto_test.go
- d:/CodingProjects/foundation/generic/monads.go
- d:/CodingProjects/foundation/generic/monads_test.go

Implementation Plan:
1. In `d:/CodingProjects/foundation/generic/monads.go`:
   - Implement `MarshalJSON() ([]byte, error)` on `Optional[T]`: returns `json.Marshal(o.val)` if valid, `[]byte("null")` if not valid.
   - Implement `UnmarshalJSON(data []byte) error` on `*Optional[T]`: if `bytes.Equal(data, []byte("null"))` or `len(data) == 0`, sets `*o = None[T]()`; otherwise unmarshals into a temporary `var val T` and on success sets `*o = Some(val)` (on error, leaves receiver untouched or resets and returns error).
   - Implement `IsZero() bool` on `Optional[T]`: returns `!o.valid` (supporting Go 1.24+ `omitzero`).
   - In `monads_test.go`: add thorough tests covering `Some("hello")`, `Some("")`, `Some(0)`, `Some(false)`, `None()`, JSON roundtrip, null, and `IsZero()`.
   - Run `go test ./generic/...` in `d:/CodingProjects/foundation`.

2. In `d:/CodingProjects/vortex/pkg/parser/binder.go`:
   - In `extractGoType(expr ast.Expr) ir.GoTypeIR`: handle `*ast.IndexExpr` and `*ast.IndexListExpr` to extract base name and element type(s). Set `goType.Name = fmt.Sprintf("%s[%s]", base.Name, elem.Name)`, `goType.ElemType = elem.Name`, `goType.IsCustomType = true`.
   - In `isDTOQueryStruct`: ensure `generic.Optional[` is not mistaken for a query DTO struct.
   - In `pkg/parser/parser_test.go`: add tests verifying generic Optional fields are correctly extracted into IR.

3. In `d:/CodingProjects/vortex/pkg/emitter/dto.go`:
   - Add `unwrapOptionalType(f *ir.FieldIR) (innerType string, isOptional bool)` that correctly recognizes both `generic.Optional[...]` and `Optional[...]` and uses `f.Type.ElemType`.
   - Update `emitFieldFormData`: for optional fields, directly unwrap inner primitives using zero-allocation encoders:
     - `string`: when `generic.Some("")`, emit `wire=` without `url.QueryEscape` (0 allocs). When `generic.None()`, omit completely. When non-empty, append escaped bytes.
     - `int*`: `strconv.AppendInt` (0 allocs).
     - `uint*`: `strconv.AppendUint` (0 allocs).
     - `float*`: `strconv.AppendFloat` (0 allocs).
     - `bool`: append `"true"`/`"false"` (or `"1"`/`"0"`) (0 allocs).
     - `time.Time`: `optVal.AppendFormat(dst, time.RFC3339)` (0 allocs, omit if `IsZero`).
     - Emit a local `appendQueryEscape` helper into emitted code or use byte-buffer appending to eliminate `url.QueryEscape` heap allocations.
   - Update `emitFieldEncodeValues`: specialize primitive optionals using `strconv.Format*` and direct string setting without `fmt.Sprint`.
   - Create `pkg/emitter/dto_test.go`: test compilation and execution of generated DTO serializers for all primitive optionals, validating that `Some("")` serializes as `wire=`, `None()` is omitted, values match expected, and zero allocations are achieved.

4. Verification:
   - Run `go test ./...` in `d:/CodingProjects/vortex` (must pass 100%).
   - Run `golangci-lint run ./...` in `d:/CodingProjects/vortex` (must pass with 0 issues).

5. Report:
   Write a comprehensive handoff report at:
   `d:/CodingProjects/vortex/.agents/worker_m1/handoff.md`
   Detail all changes, test execution commands, outputs, and verification. Send a message when done.
