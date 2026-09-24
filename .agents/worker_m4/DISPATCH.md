# DISPATCH — worker_m4 (Milestone 4 Implementation)

**Task**: Implement Performance Benchmarks & Adversarial Test Coverage
**Working Directory**: `d:/CodingProjects/vortex/.agents/worker_m4/`
**Workspace Root**: `d:/CodingProjects/vortex`

## Mandatory Integrity Warning
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_preview_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

## Mandatory Reading
1. `d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md`
2. `d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md`
3. `d:/CodingProjects/vortex/.agents/explorer_m4_1/handoff.md`
4. `d:/CodingProjects/vortex/.agents/explorer_m4_2/handoff.md`
5. `d:/CodingProjects/vortex/.agents/explorer_m4_3/handoff.md`

## Implementation Steps

### Step 1: Foundation Monad Adversarial Tests
File: `d:/CodingProjects/foundation/generic/monads_adversarial_test.go`
Append the adversarial test suite covering:
1. `TestOptional_Adversarial_PrimitiveBoundariesRoundtrip`: extreme `int64`, extreme `uint64`, extreme `float64`, Unicode/Cyrillic/CJK/Emojis (`"🔥 VORTEX 🚀"`), booleans (`true`, `false`).
2. `TestOptional_Adversarial_CollectionsAndPointersRoundtrip`: `Optional[[]int]`, `Optional[map[string]int]`, `Optional[*int]`.
3. `TestOptional_Adversarial_DirectIsZeroMatrix`: matrix verifying `Some(zeroValue).IsZero() == false` and `None[T]().IsZero() == true` across all primitive and complex types.

### Step 2: DTO In-Process Fixture
File: `d:/CodingProjects/vortex/pkg/emitter/dto_fixture_test.go`
Create the fixture implementing `BenchmarkTestDTO` with all primitive `generic.Optional[T]` fields, emitted zero-alloc methods (`AppendFormData`, `AppendQuery`, `EncodeValues`), `appendQueryEscape`, fixture constructor `newFullBenchmarkDTO()`, and `TestDTO_FixtureMatchesEmitterCodegen` (as detailed in `explorer_m4_3/handoff.md` Section 4.3).

### Step 3: Top-Level & Sub-Process Regression Benchmarks
File: `d:/CodingProjects/vortex/pkg/emitter/dto_bench_test.go`
1. Add top-level benchmarks with `b.ReportAllocs()`:
   - `BenchmarkAppendQuery_Primitives`
   - `BenchmarkAppendQuery_Optionals_Some`
   - `BenchmarkAppendQuery_Optionals_None`
   - `BenchmarkAppendFormData_Primitives`
   - `BenchmarkAppendFormData_Optionals`
   - `BenchmarkEncodeValues_ZeroAlloc`
2. Expand sub-process benchmarks in `TestEmitter_DTO_AllPrimitives_Comprehensive`:
   - `Benchmark_AppendQuery_SomeEmptyString`
   - `Benchmark_AppendQuery_AllNone`
   - `Benchmark_AppendQuery_EscapedString`
   - `Benchmark_EncodeValues_AllNone`
   - `Benchmark_EncodeValues_AllPrimitives`
   - Update regex loop to enforce 0 allocs/op for byte-buffer benchmarks and bounded allocations for `EncodeValues_AllPrimitives`.

### Step 4: Top-Level & Sub-Process Adversarial Unit Tests
File: `d:/CodingProjects/vortex/pkg/emitter/dto_test.go`
1. Add top-level `TestZeroAlloc` suite testing fixture methods directly:
   - `TestZeroAlloc_AppendQuery_Primitives` (`testing.AllocsPerRun == 0`)
   - `TestZeroAlloc_AppendFormData_Primitives` (`testing.AllocsPerRun == 0`)
   - `TestZeroAlloc_SomeEmptyString` (`out == "empty_str="`, `allocs == 0`)
   - `TestZeroAlloc_AllNone` (`len(out) == 0`, `allocs == 0`)
   - `TestZeroAlloc_BoundaryValues` (`allocs == 0`)
   - `TestZeroAlloc_QueryEscaping` (`allocs == 0`)
   - `TestZeroAlloc_EncodeValues_None` (`allocs == 0`)
   - `TestZeroAlloc_NilReceiver` (`allocs == 0`)
2. Append `TestEmitter_DTO_Adversarial_FullSuite` for full compiler end-to-end integration (nil receivers, multi-field empty strings `q=&f=`, UTF-8/emoji query unescaping, buffer growth, slice collections).

## Verification Requirements
Run and document exact outputs of:
1. `go test -benchmem -run '^$' -bench 'BenchmarkAppend' ./pkg/emitter`
2. `go test -v ./pkg/emitter -run TestZeroAlloc`
3. `$env:GOWORK="off"; go test -count=1 ./...`
4. `golangci-lint run --allow-parallel-runners ./...`

Write your comprehensive handoff report to `d:/CodingProjects/vortex/.agents/worker_m4/handoff.md`.


## 2026-09-23T12:59:47Z
You are worker_m4 (teamwork_preview_worker).
Your working directory is d:/CodingProjects/vortex/.agents/worker_m4/.
Workspace root: d:/CodingProjects/vortex.

MANDATORY INTEGRITY WARNING:
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_preview_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

You MUST read:
1. d:/CodingProjects/vortex/.agents/ORIGINAL_REQUEST.md
2. d:/CodingProjects/vortex/.agents/orchestrator_3/PROJECT.md
3. d:/CodingProjects/vortex/.agents/explorer_m4_1/handoff.md
4. d:/CodingProjects/vortex/.agents/explorer_m4_2/handoff.md
5. d:/CodingProjects/vortex/.agents/explorer_m4_3/handoff.md
6. d:/CodingProjects/vortex/.agents/worker_m4/DISPATCH.md

## 2026-09-23T13:07:47Z
**Context**: Milestone 4 Step 4
**Content**: You stopped at Step 3 completion. Please resume from Step 4: Top-level & sub-process adversarial unit tests in `d:/CodingProjects/vortex/pkg/emitter/dto_test.go`. Then run the verification commands, write `handoff.md`, and notify parent.
**Action**: Resume Step 4 and complete handoff.
