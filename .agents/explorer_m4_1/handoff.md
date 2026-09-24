# Handoff Report — explorer_m4_1: Performance Benchmarks & Adversarial Test Coverage

## 1. Observation

### 1.1 Source Files Inspected
1. `d:/CodingProjects/vortex/pkg/emitter/dto_bench_test.go` (lines 1–450)
2. `d:/CodingProjects/vortex/pkg/emitter/dto_test.go` (lines 1–265)
3. `d:/CodingProjects/vortex/pkg/emitter/dto.go` (lines 1–513)
4. `d:/CodingProjects/foundation/generic/monads.go` (lines 1–626)
5. `d:/CodingProjects/foundation/generic/monads_test.go` (lines 1–564)

### 1.2 Benchmark Coverage in `pkg/emitter/dto_bench_test.go`
In `pkg/emitter/dto_bench_test.go`:
- **Structure**:
  - `TestEmitter_DTO_AllPrimitives_Comprehensive(t *testing.T)` (lines 21–418):
    Generates a temporary Go module containing `AllPrimitivesDTO` with 21 fields (all primitive types, booleans with `@format bool_int` / `@format flag`, string, time.Time, wrapped in `generic.Optional[T]`). It emits `dto.gen.go` using `emitter.Emit`, writes an embedded test file `testAndBenchSrc`, and runs:
    ```go
    cmd := exec.Command("go", "test", "-v", "-bench=.", "-benchmem", ".")
    ```
    It parses output lines using regex:
    ```go
    benchRegex := regexp.MustCompile(`(Benchmark\w+)-\d+\s+\d+\s+([0-9.]+)\s+ns/op\s+(\d+)\s+B/op\s+(\d+)\s+allocs/op`)
    ```
    and enforces (lines 415–416):
    ```go
    require.Equalf(t, 0, bOp, "Benchmark %s must produce 0 B/op", benchName)
    require.Equalf(t, 0, allocsOp, "Benchmark %s must produce 0 allocs/op", benchName)
    ```
  - `BenchmarkEmitter_DTO_EmitCodegen(b *testing.B)` (lines 420–449):
    Benchmarks Vortex's compiler/emitter AST codegen itself (`emitter.Emit(root)`), NOT the emitted DTO methods.

- **Current Benchmarks in `testAndBenchSrc`** (lines 305–363):
  1. `Benchmark_AppendFormData_AllPrimitives(b *testing.B)` (lines 305–314)
  2. `Benchmark_AppendQuery_AllPrimitives(b *testing.B)` (lines 316–325)
  3. `Benchmark_AppendFormData_SomeEmptyString(b *testing.B)` (lines 327–336)
  4. `Benchmark_AppendFormData_AllNone(b *testing.B)` (lines 340–349)
  5. `Benchmark_AppendFormData_EscapedString(b *testing.B)` (lines 351–362)

- **Identified Benchmark Gaps**:
  - **GAP-B1 (Missing `AppendQuery` variants)**:
    `AppendQuery` is ONLY benchmarked for `AllPrimitives`. There are NO benchmarks for:
    - `Benchmark_AppendQuery_SomeEmptyString`
    - `Benchmark_AppendQuery_AllNone`
    - `Benchmark_AppendQuery_EscapedString`
  - **GAP-B2 (Missing Isolated Per-Primitive Benchmarks)**:
    Neither `AppendFormData` nor `AppendQuery` isolates individual primitive fields (`int64`, `uint64`, `float64`, `bool`, `time.Time`, `string`). The existing benchmarks only test all 21 fields together. If an allocation regression occurs in a single primitive handler, the aggregate benchmark obscures which type failed.
  - **GAP-B3 (Completely Missing `EncodeValues` Benchmarks)**:
    There are ZERO benchmarks for `EncodeValues`.
    - `Benchmark_EncodeValues_AllNone`: Must be added. When all fields are `generic.None`, `EncodeValues` performs 0 map operations and produces strictly 0 B/op and 0 allocs/op.
    - `Benchmark_EncodeValues_AllPrimitives`: Must be added. When fields are populated, stdlib `url.Values` allocates `[]string{value}` per `vals.Set`. Benchmarking this confirms elimination of `fmt.Sprint` reflection/boxing overhead, with allocations strictly bounded by map insertions.

### 1.3 Unit Test Coverage in `pkg/emitter/dto_test.go`
In `pkg/emitter/dto_test.go`:
- `TestEmitter_DTO_Emission(t *testing.T)` (lines 19–79):
  Static code string assertions on emitted Go source (checking for `AppendFormData`, `AppendQuery`, `EncodeValues`, `strconv.Append*`, and absence of `fmt.Sprint`). No runtime execution.
- `TestEmitter_DTO_ExecutionAndZeroAlloc(t *testing.T)` (lines 81–264):
  Compiles and runs `filter_test.go` against generated `SearchFilter` struct containing:
  - `Query generic.Optional[string]`
  - `Page generic.Optional[int]`
  - `Count generic.Optional[uint64]`
  - `Score generic.Optional[float64]`
  - `Active generic.Optional[bool]`
  - `Created generic.Optional[time.Time]`

- **Identified Unit Test Gaps**:
  - **GAP-T1 (Missing `AppendQuery` Zero-Alloc Assertions)**:
    `testing.AllocsPerRun(1000, func() { _ = filter.AppendFormData(buf[:0]) }) == 0` is checked in 4 tests, but `testing.AllocsPerRun` is NEVER called for `filter.AppendQuery(buf[:0])`.
  - **GAP-T2 (Missing `EncodeValues` Zero-Alloc Assertion on `None`)**:
    `filter.EncodeValues(vals)` is never verified with `testing.AllocsPerRun(1000, ...)` to ensure that an empty DTO (`AllNone`) causes 0 allocations.
  - **GAP-T3 (Missing Nil Receiver Safety and Zero-Alloc Assertions)**:
    `(*SearchFilter)(nil).AppendFormData(buf[:0])`, `(*SearchFilter)(nil).AppendQuery(buf[:0])`, and `(*SearchFilter)(nil).EncodeValues(vals)` are guarded by `if r == nil { return dst }` in `pkg/emitter/dto.go:23,36`, but no test asserts they execute safely with 0 allocations.
  - **GAP-T4 (Missing Adversarial Unicode Query Escaping Unit Test)**:
    In `dto_test.go`, query escaping is not tested on UTF-8 multi-byte characters (`日本語`, `привет`, `café`) with `testing.AllocsPerRun == 0`.
  - **GAP-T5 (Missing JSON Roundtrip & Go 1.24+ `omitzero` Integration)**:
    There is no test in `dto_test.go` validating that an emitted DTO struct containing `generic.Optional[T]` correctly roundtrips via `encoding/json` (`null` for `None`, inner value for `Some`, empty string `""` preserved, `0` preserved, `false` preserved) and integrates with `json:",omitzero"`.
  - **GAP-T6 (Missing Per-Primitive Isolated Zero-Alloc Assertions)**:
    In `dto_test.go`, fields are tested either all unset (`None`) or all set (`makeFullDTO`). Isolated per-field tests asserting `testing.AllocsPerRun == 0` for `Query` alone, `Page` alone, `Count` alone, `Score` alone, `Active` alone, and `Created` alone do not exist.

---

## 2. Logic Chain

1. **Premise 1 (Zero-Allocation Invariant)**: Vortex Milestone 4 requires verified 0 allocs/op for primitive DTO serialization methods (`AppendQuery`, `AppendFormData`) and unset `EncodeValues`.
2. **Premise 2 (Symmetric Method Delegation)**: In `pkg/emitter/dto.go:31–33`, `AppendQuery` delegates directly:
   ```go
   func (r *T) AppendQuery(dst []byte) []byte {
       return r.AppendFormData(dst)
   }
   ```
   However, `dto_test.go` only asserted `testing.AllocsPerRun == 0` on `AppendFormData`. An accidental refactoring or compiler inline failure could cause `AppendQuery` to escape `dst` to the heap without failing existing tests. Adding explicit `testing.AllocsPerRun == 0` assertions for `AppendQuery` ensures ironclad regression protection.
3. **Premise 3 (URL Values Allocation Semantics)**: `url.Values` is `map[string][]string`. Calling `vals.Set(k, v)` always allocates a 1-element slice `[]string{v}` in Go runtime. Therefore:
   - For `AllNone`, no `vals.Set` calls are executed; allocations must be strictly `0`.
   - For populated fields, allocations originate strictly from `vals.Set` slice creation (exactly equal to number of set fields) with zero `fmt.Sprint` overhead. The benchmark harness in `dto_bench_test.go` must assert `0 allocs/op` for `AllNone` and bounded allocations (<= non-empty field count) for `AllPrimitives`.
4. **Premise 4 (Acceptance Criteria Alignment)**: Acceptance criteria in `ORIGINAL_REQUEST.md` require:
   - `DTO emission for generic.Optional[T] produces exactly 0 allocs/op for primitives without calling fmt.Sprint.`
   - `Explicit empty string generic.Some("") serializes as key= in query and form payloads.`
   - `Unset optional (generic.None()) appends nothing to the payload.`
   - `generic.Optional[T] correctly marshals to and unmarshals from JSON matching the inner value or null.`
   Expanding `dto_bench_test.go` and `dto_test.go` with the designed suites directly satisfies all acceptance criteria.

---

## 3. Caveats

1. **No Production Compiler Modifications Needed**: `pkg/emitter/dto.go` and `foundation/generic/monads.go` are already fully implemented with zero-allocation serializers (`strconv.AppendInt`, `strconv.AppendUint`, `strconv.AppendFloat`, `appendQueryEscape`, `time.AppendFormat`). Only the test and benchmark suites need expansion.
2. **Subprocess Test Execution Overhead**: Both `dto_bench_test.go` and `dto_test.go` use `os/exec.Command("go", "test", ...)` within `t.TempDir()`. Running multiple benchmark iterations inside `exec.Command` adds ~15–25 seconds of test execution time. Benchmark iterations should be clean and scoped to avoid test timeouts.
3. **Stdlib `url.Values` Allocations**: As established, `vals.Set` allocates a slice. Any benchmark asserting `EncodeValues` on a populated DTO cannot demand `0 allocs/op` unless a custom value map is introduced. Asserting `0 allocs/op` for `Benchmark_EncodeValues_AllNone` and verifying `bOp <= 20 * sliceSize` for `Benchmark_EncodeValues_AllPrimitives` is the correct technical boundary.

---

## 4. Conclusion & Concrete Implementation Plan for `worker_m4`

`worker_m4` should implement two focused updates:

### 4.1 Concrete Additions for `pkg/emitter/dto_bench_test.go`

In `TestEmitter_DTO_AllPrimitives_Comprehensive`:
1. **Add `AppendQuery` symmetric benchmarks** to `testAndBenchSrc`:
   - `Benchmark_AppendQuery_SomeEmptyString(b *testing.B)`
   - `Benchmark_AppendQuery_AllNone(b *testing.B)`
   - `Benchmark_AppendQuery_EscapedString(b *testing.B)`
2. **Add Isolated Per-Primitive benchmarks** to `testAndBenchSrc`:
   - `Benchmark_AppendFormData_Int64(b *testing.B)`
   - `Benchmark_AppendQuery_Int64(b *testing.B)`
   - `Benchmark_AppendFormData_Uint64(b *testing.B)`
   - `Benchmark_AppendQuery_Uint64(b *testing.B)`
   - `Benchmark_AppendFormData_Float64(b *testing.B)`
   - `Benchmark_AppendQuery_Float64(b *testing.B)`
   - `Benchmark_AppendFormData_Bool(b *testing.B)`
   - `Benchmark_AppendQuery_Bool(b *testing.B)`
   - `Benchmark_AppendFormData_Time(b *testing.B)`
   - `Benchmark_AppendQuery_Time(b *testing.B)`
3. **Add `EncodeValues` benchmarks** to `testAndBenchSrc`:
   - `Benchmark_EncodeValues_AllNone(b *testing.B)`
   - `Benchmark_EncodeValues_AllPrimitives(b *testing.B)`
4. **Update Benchmark Result Parser** in `dto_bench_test.go:403–417`:
   Differentiate between zero-alloc byte-buffer benchmarks vs map-allocating `EncodeValues_AllPrimitives`:
   ```go
   for _, m := range matches {
       benchName := m[1]
       nsOp := m[2]
       bOp, _ := strconv.Atoi(m[3])
       allocsOp, _ := strconv.Atoi(m[4])

       t.Logf("Benchmark %-45s: %8s ns/op | %3d B/op | %3d allocs/op", benchName, nsOp, bOp, allocsOp)

       if benchName == "Benchmark_EncodeValues_AllPrimitives" {
           // EncodeValues on populated fields allocates slices via stdlib url.Values.Set,
           // but must NOT call fmt.Sprint (bounded to 1 alloc per populated field = <= 20 allocs/op).
           require.LessOrEqualf(t, allocsOp, 20, "Benchmark %s allocs must be bounded by url.Values map sets", benchName)
       } else {
           // All AppendFormData, AppendQuery, and EncodeValues_AllNone benchmarks must produce EXACTLY 0 B/op and 0 allocs/op.
           require.Equalf(t, 0, bOp, "Benchmark %s must produce 0 B/op", benchName)
           require.Equalf(t, 0, allocsOp, "Benchmark %s must produce 0 allocs/op", benchName)
       }
   }
   ```

### 4.2 Concrete Additions for `pkg/emitter/dto_test.go`

In `TestEmitter_DTO_ExecutionAndZeroAlloc`:
1. **Add `AppendQuery` zero-allocation assertions** to existing subtests:
   - In `TestSearchFilter_NoneOmissionAndZeroAlloc`: verify `filter.AppendQuery(buf[:0])` produces 0 allocs.
   - In `TestSearchFilter_SomeEmptyString`: verify `filter.AppendQuery(buf[:0]) == "q="` and produces 0 allocs.
   - In `TestSearchFilter_PrimitivesRoundtripAndZeroAlloc`: verify `filter.AppendQuery(buf[:0])` produces 0 allocs.
   - In `TestSearchFilter_ExplicitFalseAndZero`: verify `filter.AppendQuery(buf[:0])` produces 0 allocs.
2. **Add `EncodeValues` zero-allocation assertion on `None`**:
   - In `TestSearchFilter_NoneOmissionAndZeroAlloc`: verify `filter.EncodeValues(vals)` produces 0 allocs on empty DTO.
3. **Add Nil Receiver Safety test** (`TestSearchFilter_NilReceiver`):
   - Verifies `(*SearchFilter)(nil).AppendFormData(buf[:0])` returns `buf[:0]` with 0 allocs.
   - Verifies `(*SearchFilter)(nil).AppendQuery(buf[:0])` returns `buf[:0]` with 0 allocs.
   - Verifies `(*SearchFilter)(nil).EncodeValues(vals)` produces 0 allocs.
4. **Add Adversarial Unicode & Query Escaping test** (`TestSearchFilter_UnicodeAndAdversarialQueryEscaping`):
   - Tests `Query: generic.Some("search & replace = 100% / 日本語 + café ? #tag")`.
   - Asserts exact RFC 3986 escaped output for `AppendFormData` and `AppendQuery`.
   - Asserts `testing.AllocsPerRun(1000, ...) == 0` for both.
5. **Add JSON Roundtrip & Go 1.24+ `omitzero` test** (`TestSearchFilter_JSON_RoundtripAndOmitZero`):
   - Updates test struct tags to include `json:"q,omitzero"` etc.
   - Tests unset DTO marshals to `{}`.
   - Tests explicit empty string `Some("")`, zero int `Some(0)`, and false `Some(false)` are retained in JSON.
   - Tests unmarshaling populated JSON restores `generic.Some(...)`.
   - Tests unmarshaling `{"q":null}` resets field to `generic.None[string]()`.
6. **Add Isolated Per-Primitive Unit Tests** (`TestSearchFilter_IndividualPrimitivesZeroAlloc`):
   - Asserts `testing.AllocsPerRun(1000, ...) == 0` for each primitive field populated in isolation for both `AppendFormData` and `AppendQuery`.

---

## 5. Verification Method

1. **Run emitter package tests**:
   ```pwsh
   $env:GOWORK="off"
   go test -v -count=1 ./pkg/emitter
   ```
   *Expected output*: `PASS: TestEmitter_DTO_AllPrimitives_Comprehensive`, `PASS: TestEmitter_DTO_ExecutionAndZeroAlloc`, `PASS: TestEmitter_DTO_Emission`.

2. **Run full workspace test suite**:
   ```pwsh
   $env:GOWORK="off"
   go test -count=1 ./...
   ```
   *Expected output*: All packages pass cleanly across the entire workspace.

3. **Run linter**:
   ```pwsh
   golangci-lint run
   ```
   *Expected output*: `0 issues`.

4. **Invalidation Conditions**:
   - Any benchmark in `dto_bench_test.go` reporting `>0 B/op` or `>0 allocs/op` for `Append*` or `EncodeValues_AllNone`.
   - Any `testing.AllocsPerRun(1000, ...) != 0` in `dto_test.go`.
   - Any failure in JSON marshaling/unmarshaling roundtrip or `omitzero` handling.
   - Any lint failure in `golangci-lint run`.
